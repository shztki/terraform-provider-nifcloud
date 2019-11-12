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

func resourceNifcloudRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudRouterCreate,
		Read:   resourceNifcloudRouterRead,
		Update: resourceNifcloudRouterUpdate,
		Delete: resourceNifcloudRouterDelete,
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
			"router_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "small",
			},
			"accounting_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2",
			},
			"network_interfaces": {
				Type:          schema.TypeSet,
				Optional:      true,
				MinItems:      1,
				MaxItems:      7,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipaddress": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dhcp": {
							Type:     schema.TypeBool,
							Optional: true,
							Default: false,
						},
						"dhcp_options_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dhcp_config_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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

func resourceNifcloudRouterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	var securityGroups []*string
	if sgs := d.Get("security_groups").([]interface{}); sgs != nil {
		for _, v := range sgs {
			securityGroups = append(securityGroups, nifcloud.String(v.(string)))
		}
	}

	var networkInterfaces []*computing.RequestNetworkInterfaceStruct
	if interfaces, ok := d.GetOk("network_interfaces"); ok {
		for _, ni := range interfaces.(*schema.Set).List() {
			networkInterface := &computing.RequestNetworkInterfaceStruct{}
			if v, ok := ni.(map[string]interface{}); ok {
				networkInterface.SetNetworkId(v["network_id"].(string))
				networkInterface.SetNetworkName(v["network_name"].(string))
				networkInterface.SetIpAddress(v["ipaddress"].(string))
				networkInterface.SetDhcp(v["dhcp"].(bool))
				networkInterface.SetDhcpOptionsId(v["dhcp_options_id"].(string))
				networkInterface.SetDhcpConfigId(v["dhcp_config_id"].(string))
			}
			networkInterfaces = append(networkInterfaces, networkInterface)
		}
	}

	createOpts := &computing.NiftyCreateRouterInput{
		RouterName:       nifcloud.String(d.Get("name").(string)),
		NetworkInterface: networkInterfaces,
		Description:      nifcloud.String(d.Get("description").(string)),
		Type:             nifcloud.String(d.Get("router_type").(string)),
		AccountingType:   nifcloud.String(d.Get("accounting_type").(string)),
		SecurityGroup:    securityGroups,
		AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
//		Placement:        &computing.RequestPlacementStruct{AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string))},
	}

	// Create the router.
	log.Printf("[DEBUG] Creating router")
	resp, err := conn.NiftyCreateRouter(createOpts)
//	var resp *computing.NiftyCreateRouterOutput
//	err := resource.Retry(20*time.Minute, func() *resource.RetryError {
//		var err error
//		resp, err = conn.NiftyCreateRouter(createOpts)
//
//		// Retry for ...
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
//		resp, err = conn.NiftyCreateRouter(createOpts)
//	}

	if err != nil {
		return fmt.Errorf("Error creating router: %s", err)
	}

	// Store the ID
	router := resp.Router
	d.SetId(*router.RouterId)
	log.Printf("[INFO] router ID: %s", *router.RouterId)

	// Wait for the router to be available.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "warning"},
		Target:     []string{"available"},
		Refresh:    routerRefreshFunc(conn, *router.RouterId),
		Timeout:    15 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for router (%s) to become ready: %s",
			*router.RouterId, stateErr)
	}

	return nil
}

func routerRefreshFunc(conn *computing.Computing, routerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		routerFilter := &computing.RequestFilterStruct{
			Name:         nifcloud.String("router-id"),
			RequestValue: []*string{nifcloud.String(routerID)},
		}

		resp, err := conn.NiftyDescribeRouters(&computing.NiftyDescribeRoutersInput{
			Filter: []*computing.RequestFilterStruct{routerFilter},
		})
		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.RouterId" {
				resp = nil
			} else {
				log.Printf("Error on routerRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil || len(resp.RouterSet) == 0 {
			// handle consistency issues
			return nil, "", nil
		}

		router := resp.RouterSet[0]
		return router, *router.State, nil
	}
}

func resourceNifcloudRouterExists(routerID string, conn *computing.Computing) (bool, error) {
	nameFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("router-id"),
		RequestValue: []*string{nifcloud.String(routerID)},
	}

	resp, err := conn.NiftyDescribeRouters(&computing.NiftyDescribeRoutersInput{
		Filter: []*computing.RequestFilterStruct{nameFilter},
	})
	if err != nil {
		return false, err
	}

	if len(resp.RouterSet) > 0 && *resp.RouterSet[0].State != "stopped" {
		return true, nil
	}

	return false, nil
}

func resourceNifcloudRouterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	routerFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("router-id"),
		RequestValue: []*string{nifcloud.String(d.Id())},
	}

	resp, err := conn.NiftyDescribeRouters(&computing.NiftyDescribeRoutersInput{
		Filter: []*computing.RequestFilterStruct{routerFilter},
	})
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.RouterId" {
			d.SetId("")
			return nil
		}
		log.Printf("[ERROR] Error finding router: %s", err)
		return err
	}

	if len(resp.RouterSet) != 1 {
		return fmt.Errorf("Error finding router: %s", d.Id())
	}

	log.Printf("[DEBUG] router describe %v", resp)
	if *resp.RouterSet[0].State == "deleted" {
		log.Printf("[INFO] router is in `deleted` state: %s", d.Id())
		d.SetId("")
		return nil
	}

	router := resp.RouterSet[0]
	d.Set("name", router.RouterName)
	d.Set("description", router.Description)
	d.Set("router_type", router.Type)
	d.Set("accounting_type", router.NextMonthAccountingType)
	d.Set("availability_zone", router.AvailabilityZone)

	sgs := make([]string, 0, len(router.GroupSet))
	for _, sg := range router.GroupSet {
		sgs = append(sgs, *sg.GroupId)
	}

	log.Printf("[DEBUG] Setting Security Group Ids: %#v", sgs)
	if err := d.Set("security_groups", sgs); err != nil {
		return err
	}

	return nil
}

func resourceNifcloudRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	d.Partial(true)

	log.Printf("[INFO] Updating router %s", d.Id())
	if d.HasChange("description") {
		input := computing.NiftyModifyRouterAttributeInput{
			RouterId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("description"),
			Value:        nifcloud.String(d.Get("description").(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyRouterAttribute(&input)
		if err != nil {
			return fmt.Errorf("error description updating router (%s): %s", d.Id(), err)
		}	
	}
	if d.HasChange("name") {
		input := computing.NiftyModifyRouterAttributeInput{
			RouterId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("routerName"),
			Value:        nifcloud.String(d.Get("name").(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyRouterAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating router (%s): %s", d.Id(), err)
		}

		// Wait for the router to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    routerRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for router (%s) to become ready: %s",
				d.Id(), stateErr)
		}		
	}
	if d.HasChange("accounting_type") {
		input := computing.NiftyModifyRouterAttributeInput{
			RouterId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("accountingType"),
			Value:        nifcloud.String(d.Get("accounting_type").(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyRouterAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating router (%s): %s", d.Id(), err)
		}

		// Wait for the router to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    routerRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for router (%s) to become ready: %s",
				d.Id(), stateErr)
		}
	}
	if d.HasChange("security_groups") {
		input := computing.NiftyModifyRouterAttributeInput{
			RouterId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("groupId"),
			Value:        nifcloud.String(d.Get("security_groups").([]interface{})[0].(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyRouterAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating router (%s): %s", d.Id(), err)
		}

		// Wait for the router to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    routerRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for router (%s) to become ready: %s",
				d.Id(), stateErr)
		}		
	}
	if d.HasChange("router_type") {
		input := computing.NiftyModifyRouterAttributeInput{
			RouterId: nifcloud.String(d.Id()),
			Attribute:    nifcloud.String("type"),
			Value:        nifcloud.String(d.Get("router_type").(string)),
			Agreement:    nifcloud.Bool(false),
		}
		_, err := conn.NiftyModifyRouterAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating router (%s): %s", d.Id(), err)
		}

		// Wait for the router to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    routerRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for router (%s) to become ready: %s",
				d.Id(), stateErr)
		}
	}
	if d.HasChange("network_interfaces") {
		var networkInterfaces []*computing.RequestNetworkInterfaceStruct
		if interfaces, ok := d.GetOk("network_interfaces"); ok {
			for _, ni := range interfaces.(*schema.Set).List() {
				networkInterface := &computing.RequestNetworkInterfaceStruct{}
				if v, ok := ni.(map[string]interface{}); ok {
					networkInterface.SetNetworkId(v["network_id"].(string))
					networkInterface.SetNetworkName(v["network_name"].(string))
					networkInterface.SetIpAddress(v["ipaddress"].(string))
					networkInterface.SetDhcp(v["dhcp"].(bool))
					networkInterface.SetDhcpOptionsId(v["dhcp_options_id"].(string))
					networkInterface.SetDhcpConfigId(v["dhcp_config_id"].(string))
				}
				networkInterfaces = append(networkInterfaces, networkInterface)
			}
		}

		input := computing.NiftyUpdateRouterNetworkInterfacesInput{
			RouterId:         nifcloud.String(d.Id()),
			NetworkInterface: networkInterfaces,
			NiftyReboot:      nifcloud.String("true"),
			Agreement:        nifcloud.Bool(false),
		}
		_, err := conn.NiftyUpdateRouterNetworkInterfaces(&input)
		if err != nil {
			return fmt.Errorf("error name updating router (%s): %s", d.Id(), err)
		}

		// Wait for the router to be available.
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "warning"},
			Target:     []string{"available"},
			Refresh:    routerRefreshFunc(conn, d.Id()),
			Timeout:    15 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmt.Errorf(
				"Error waiting for router (%s) to become ready: %s",
				d.Id(), stateErr)
		}		
	}

	d.Partial(false)

	return resourceNifcloudRouterRead(d, meta)
}

func resourceNifcloudRouterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	request := &computing.NiftyDeleteRouterInput{
		RouterId: nifcloud.String(d.Id()),
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.NiftyDeleteRouter(request)
		log.Printf("[DEBUG] deleting router %v", resp)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.RouterId" {
				return nil
			}
			return resource.RetryableError(err)
		}

		return nil
	})
/*
	_, err := conn.NiftyDeleteRouter(&computing.NiftyDeleteRouterInput{
		RouterId: nifcloud.String(d.Id()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.RouterId" {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting router: %s", err)
	}
*/
	routerFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("router-id"),
		RequestValue: []*string{nifcloud.String(d.Id())},
	}

	input := &computing.NiftyDescribeRoutersInput{
		Filter: []*computing.RequestFilterStruct{routerFilter},
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.NiftyDescribeRouters(input)
		log.Printf("[DEBUG] delete after describe 001 %v", resp)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.RouterId" {
				return nil
			}
			return resource.RetryableError(err)
		}

		err = checkRouterDeleteResponse(resp, d.Id())
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})

	if isResourceTimeoutError(err) {
		var resp *computing.NiftyDescribeRoutersOutput
		resp, err = conn.NiftyDescribeRouters(input)
		log.Printf("[DEBUG] delete after describe 002 %v", resp)

		if err != nil {
			return checkRouterDeleteResponse(resp, d.Id())
		}
	}

	if err != nil {
		return fmt.Errorf("Error deleting router: %s", err)
	}
	return nil

}

func checkRouterDeleteResponse(resp *computing.NiftyDescribeRoutersOutput, id string) error {
	if resp.RouterSet == nil {
		return nil
	}

	switch *resp.RouterSet[0].State {
	case "available":
		return fmt.Errorf("router (%s) in state (%s), retrying", id, *resp.RouterSet[0].State)
	case "pending":
		return nil
	default:
		return fmt.Errorf("Unrecognized state (%s) for router delete on (%s)", *resp.RouterSet[0].State, id)
	}
}

func expandNetworkInterface(configured []interface{}) ([]*computing.RequestNetworkInterfaceStruct, error) {
	var parameters []*computing.RequestNetworkInterfaceStruct

	// Loop over our configured parameters and create
	// an array of aws-sdk-go compatible objects
	for _, pRaw := range configured {
		data := pRaw.(map[string]interface{})

		if data["network_id"].(string) == "" {
			continue
		}

		p := &computing.RequestNetworkInterfaceStruct{
			NetworkId:     nifcloud.String(data["network_id"].(string)),
			NetworkName:   nifcloud.String(data["network_name"].(string)),
			IpAddress:     nifcloud.String(data["ipaddress"].(string)),
			Dhcp:          nifcloud.Bool(data["dhcp"].(bool)),
			DhcpOptionsId: nifcloud.String(data["dhcp_options_id"].(string)),
			DhcpConfigId:  nifcloud.String(data["dhcp_config_id"].(string)),
		}

		parameters = append(parameters, p)
	}

	return parameters, nil
}