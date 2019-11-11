package nifcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNifcloudVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudVpnGatewayCreate,
		Read:   resourceNifcloudVpnGatewayRead,
		Update: resourceNifcloudVpnGatewayUpdate,
		Delete: resourceNifcloudVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 15),
			},
			"vpn_gateway_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "small",
			},
			"accounting_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2",
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNifcloudVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	var securityGroups []*string
	if sgs := d.Get("security_groups").([]interface{}); sgs != nil {
		for _, v := range sgs {
			securityGroups = append(securityGroups, nifcloud.String(v.(string)))
		}
	}

	niftyNetwork := &computing.RequestNiftyNetworkStruct{
		NetworkId: nifcloud.String(d.Get("network_id").(string)),
		IpAddress: nifcloud.String(d.Get("private_ip_address").(string)),
	}

	createOpts := &computing.CreateVpnGatewayInput{
		NiftyVpnGatewayName:        nifcloud.String(d.Get("name").(string)),
		NiftyNetwork:               niftyNetwork,
		NiftyVpnGatewayDescription: nifcloud.String(d.Get("description").(string)),
		NiftyVpnGatewayType:        nifcloud.String(d.Get("vpn_gateway_type").(string)),
		AccountingType:             nifcloud.String(d.Get("accounting_type").(string)),
		SecurityGroup:              securityGroups,
		Placement:                  &computing.RequestPlacementStruct{AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string))},
	}

	// Create the Vpn Gateway.
	log.Printf("[DEBUG] Creating vpn gateway")
	resp, err := conn.CreateVpnGateway(createOpts)
//	var resp *computing.CreateVpnGatewayOutput
//	err := resource.Retry(20*time.Minute, func() *resource.RetryError {
//		var err error
//		resp, err = conn.CreateVpnGateway(createOpts)
//
//		// Retry for ...
//		//log.Printf("[DEBUG] Creating vpn gateway err: %s", err)
//		if isNifcloudErr(err, "Server.ResourceIncorrectState.Network.Processing", "") {
//			return resource.RetryableError(err)
//		}
//
//		if err != nil {
//			return resource.NonRetryableError(err)
//		}
//
//		return nil
//	})
//
//	if isResourceTimeoutError(err) {
//		resp, err = conn.CreateVpnGateway(createOpts)
//	}

	if err != nil {
		return fmt.Errorf("Error creating vpn gateway: %s", err)
	}

	// Store the ID
	vpnGateway := resp.VpnGateway
	d.SetId(*vpnGateway.VpnGatewayId)
	log.Printf("[INFO] Vpn gateway ID: %s", *vpnGateway.VpnGatewayId)

	// Wait for the VpnGateway to be available.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "warning"},
		Target:     []string{"available"},
		Refresh:    vpnGatewayRefreshFunc(conn, *vpnGateway.VpnGatewayId),
		Timeout:    15 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for vpn gateway (%s) to become ready: %s",
			*vpnGateway.VpnGatewayId, stateErr)
	}

	return nil
}

func vpnGatewayRefreshFunc(conn *computing.Computing, gatewayID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		gatewayFilter := &computing.RequestFilterStruct{
			Name:         nifcloud.String("vpn-gateway-id"),
			RequestValue: []*string{nifcloud.String(gatewayID)},
		}

		resp, err := conn.DescribeVpnGateways(&computing.DescribeVpnGatewaysInput{
			Filter: []*computing.RequestFilterStruct{gatewayFilter},
		})
		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.VpnGatewayId" {
				resp = nil
			} else {
				log.Printf("Error on VpnGatewayRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil || len(resp.VpnGatewaySet) == 0 {
			// handle consistency issues
			return nil, "", nil
		}

		gateway := resp.VpnGatewaySet[0]
		return gateway, *gateway.State, nil
	}
}

func resourceNifcloudVpnGatewayExists(name string, conn *computing.Computing) (bool, error) {
	nameFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("nifty-vpn-gateway-name"),
		RequestValue: []*string{nifcloud.String(name)},
	}

	resp, err := conn.DescribeVpnGateways(&computing.DescribeVpnGatewaysInput{
		Filter: []*computing.RequestFilterStruct{nameFilter},
	})
	if err != nil {
		return false, err
	}

	if len(resp.VpnGatewaySet) > 0 && *resp.VpnGatewaySet[0].State != "stopped" {
		return true, nil
	}

	return false, nil
}

func resourceNifcloudVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	gatewayFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("vpn-gateway-id"),
		RequestValue: []*string{nifcloud.String(d.Id())},
	}

	resp, err := conn.DescribeVpnGateways(&computing.DescribeVpnGatewaysInput{
		Filter: []*computing.RequestFilterStruct{gatewayFilter},
	})
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.VpnGatewayId" {
			d.SetId("")
			return nil
		}
		log.Printf("[ERROR] Error finding VpnGateway: %s", err)
		return err
	}

	if len(resp.VpnGatewaySet) != 1 {
		return fmt.Errorf("Error finding VpnGateway: %s", d.Id())
	}

	log.Printf("[DEBUG] vpn gateway describe %v", resp)
	if *resp.VpnGatewaySet[0].State == "deleted" {
		log.Printf("[INFO] Vpn Gateway is in `deleted` state: %s", d.Id())
		d.SetId("")
		return nil
	}

	vpnGateway := resp.VpnGatewaySet[0]
	d.Set("name", vpnGateway.NiftyVpnGatewayName)
	d.Set("ip_address", vpnGateway.IpAddress)
	d.Set("description", vpnGateway.NiftyVpnGatewayDescription)
	d.Set("vpn_gateway_type", vpnGateway.NiftyVpnGatewayType)
	d.Set("accounting_type", vpnGateway.NextMonthAccountingType)
	d.Set("availability_zone", vpnGateway.AvailabilityZone)

	sgs := make([]string, 0, len(vpnGateway.GroupSet))
	for _, sg := range vpnGateway.GroupSet {
		sgs = append(sgs, *sg.GroupId)
	}

	log.Printf("[DEBUG] Setting Security Group Ids: %#v", sgs)
	if err := d.Set("security_groups", sgs); err != nil {
		return err
	}

	return nil
}

func resourceNifcloudVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	d.Partial(true)

	log.Printf("[INFO] Updating VpnGateway %s", d.Id())
	if d.HasChange("description") {
		input := computing.NiftyModifyVpnGatewayAttributeInput{
			VpnGatewayId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("niftyVpnGatewayDescription"),
			Value:        nifcloud.String(d.Get("description").(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyVpnGatewayAttribute(&input)
		if err != nil {
			return fmt.Errorf("error description updating VpnGateway (%s): %s", d.Id(), err)
		}	
	}
	if d.HasChange("name") {
		input := computing.NiftyModifyVpnGatewayAttributeInput{
			VpnGatewayId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("niftyVpnGatewayName"),
			Value:        nifcloud.String(d.Get("name").(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyVpnGatewayAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating VpnGateway (%s): %s", d.Id(), err)
		}

		// Wait for the VpnGateway to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    vpnGatewayRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for vpn gateway (%s) to become ready: %s",
				d.Id(), stateErr)
		}		
	}
	if d.HasChange("accounting_type") {
		input := computing.NiftyModifyVpnGatewayAttributeInput{
			VpnGatewayId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("niftyVpnGatewayAccountingType"),
			Value:        nifcloud.String(d.Get("accounting_type").(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyVpnGatewayAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating VpnGateway (%s): %s", d.Id(), err)
		}

		// Wait for the VpnGateway to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    vpnGatewayRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for vpn gateway (%s) to become ready: %s",
				d.Id(), stateErr)
		}
	}
	if d.HasChange("security_groups") {
		input := computing.NiftyModifyVpnGatewayAttributeInput{
			VpnGatewayId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("groupId"),
			Value:        nifcloud.String(d.Get("security_groups").([]interface{})[0].(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyVpnGatewayAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating VpnGateway (%s): %s", d.Id(), err)
		}

		// Wait for the VpnGateway to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    vpnGatewayRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for vpn gateway (%s) to become ready: %s",
				d.Id(), stateErr)
		}		
	}
	if d.HasChange("private_ip_address") {
		niftyNetwork := &computing.RequestNetworkInterfaceStruct{
			NetworkId: nifcloud.String(d.Get("network_id").(string)),
			IpAddress: nifcloud.String(d.Get("private_ip_address").(string)),
		}
	
		input := computing.NiftyUpdateVpnGatewayNetworkInterfacesInput{
			VpnGatewayId:     nifcloud.String(d.Id()),
			NetworkInterface: niftyNetwork,
			NiftyReboot:      nifcloud.String("true"),
			Agreement:        nifcloud.Bool(false),
		}
		_, err := conn.NiftyUpdateVpnGatewayNetworkInterfaces(&input)
		if err != nil {
			return fmt.Errorf("error name updating VpnGateway (%s): %s", d.Id(), err)
		}

		// Wait for the VpnGateway to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    vpnGatewayRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for vpn gateway (%s) to become ready: %s",
				d.Id(), stateErr)
		}		
	}
	if d.HasChange("vpn_gateway_type") {
		input := computing.NiftyModifyVpnGatewayAttributeInput{
			VpnGatewayId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("niftyVpnGatewayType"),
			Value:        nifcloud.String(d.Get("vpn_gateway_type").(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyVpnGatewayAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating VpnGateway (%s): %s", d.Id(), err)
		}

		// Wait for the VpnGateway to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    vpnGatewayRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for vpn gateway (%s) to become ready: %s",
				d.Id(), stateErr)
		}
	}

	d.Partial(false)

	return resourceNifcloudVpnGatewayRead(d, meta)
}

func resourceNifcloudVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	request := &computing.DeleteVpnGatewayInput{
		VpnGatewayId: nifcloud.String(d.Id()),
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DeleteVpnGateway(request)
		log.Printf("[DEBUG] deleting vpn gateway %v", resp)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.VpnGatewayId" {
				return nil
			}
			return resource.RetryableError(err)
		}

		return nil
	})
/*
	_, err := conn.DeleteVpnGateway(&computing.DeleteVpnGatewayInput{
		VpnGatewayId: nifcloud.String(d.Id()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.VpnGatewayId" {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting VpnGateway: %s", err)
	}
*/
	gatewayFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("vpn-gateway-id"),
		RequestValue: []*string{nifcloud.String(d.Id())},
	}

	input := &computing.DescribeVpnGatewaysInput{
		Filter: []*computing.RequestFilterStruct{gatewayFilter},
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DescribeVpnGateways(input)
		log.Printf("[DEBUG] delete after describe 001 %v", resp)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.VpnGatewayId" {
				return nil
			}
			return resource.RetryableError(err)
		}

		err = checkVpnGatewayDeleteResponse(resp, d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if isResourceTimeoutError(err) {
		var resp *computing.DescribeVpnGatewaysOutput
		resp, err = conn.DescribeVpnGateways(input)
		log.Printf("[DEBUG] delete after describe 002 %v", resp)

		if err != nil {
			return checkVpnGatewayDeleteResponse(resp, d.Id())
		}
	}

	if err != nil {
		return fmt.Errorf("Error deleting vpn gateway: %s", err)
	}
	return nil

}

func checkVpnGatewayDeleteResponse(resp *computing.DescribeVpnGatewaysOutput, id string) error {
	if resp.VpnGatewaySet == nil {
		return nil
	}

	switch *resp.VpnGatewaySet[0].State {
	case "available":
		return fmt.Errorf("Gateway (%s) in state (%s), retrying", id, *resp.VpnGatewaySet[0].State)
	case "pending":
		return nil
	default:
		return fmt.Errorf("Unrecognized state (%s) for Vpn Gateway delete on (%s)", *resp.VpnGatewaySet[0].State, id)
	}
}
