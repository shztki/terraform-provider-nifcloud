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

func resourceNifcloudCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudCustomerGatewayCreate,
		Read:   resourceNifcloudCustomerGatewayRead,
		Update: resourceNifcloudCustomerGatewayUpdate,
		Delete: resourceNifcloudCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 15),
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lan_side_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lan_side_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNifcloudCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	name := d.Get("name").(string)

	alreadyExists, err := resourceNifcloudCustomerGatewayExists(name, conn)
	if err != nil {
		return err
	}

	if alreadyExists {
		return fmt.Errorf("An existing customer gateway for Name: %s has been found", name)
	}

	createOpts := &computing.CreateCustomerGatewayInput{
		NiftyCustomerGatewayName:        nifcloud.String(name),
		IpAddress:                       nifcloud.String(d.Get("ip_address").(string)),
		NiftyCustomerGatewayDescription: nifcloud.String(d.Get("description").(string)),
		NiftyLanSideCidrBlock:           nifcloud.String(d.Get("lan_side_cidr_block").(string)),
		NiftyLanSideIpAddress:           nifcloud.String(d.Get("lan_side_ip_address").(string)),
	}

	// Create the Customer Gateway.
	log.Printf("[DEBUG] Creating customer gateway")
	resp, err := conn.CreateCustomerGateway(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating customer gateway: %s", err)
	}

	// Store the ID
	customerGateway := resp.CustomerGateway
	d.SetId(*customerGateway.CustomerGatewayId)
	log.Printf("[INFO] Customer gateway ID: %s", *customerGateway.CustomerGatewayId)

	// Wait for the CustomerGateway to be available.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    customerGatewayRefreshFunc(conn, *customerGateway.CustomerGatewayId),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for customer gateway (%s) to become ready: %s",
			*customerGateway.CustomerGatewayId, err)
	}

	return nil
}

func customerGatewayRefreshFunc(conn *computing.Computing, gatewayID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		gatewayFilter := &computing.RequestFilterStruct{
			Name:         nifcloud.String("customer-gateway-id"),
			RequestValue: []*string{nifcloud.String(gatewayID)},
		}

		resp, err := conn.DescribeCustomerGateways(&computing.DescribeCustomerGatewaysInput{
			Filter: []*computing.RequestFilterStruct{gatewayFilter},
		})
		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidCustomerGatewayID.NotFound" {
				resp = nil
			} else {
				log.Printf("Error on CustomerGatewayRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil || len(resp.CustomerGatewaySet) == 0 {
			// handle consistency issues
			return nil, "", nil
		}

		gateway := resp.CustomerGatewaySet[0]
		return gateway, *gateway.State, nil
	}
}

func resourceNifcloudCustomerGatewayExists(name string, conn *computing.Computing) (bool, error) {
	nameFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("nifty-customer-gateway-name"),
		RequestValue: []*string{nifcloud.String(name)},
	}

	resp, err := conn.DescribeCustomerGateways(&computing.DescribeCustomerGatewaysInput{
		Filter: []*computing.RequestFilterStruct{nameFilter},
	})
	if err != nil {
		return false, err
	}

	if len(resp.CustomerGatewaySet) > 0 && *resp.CustomerGatewaySet[0].State != "stopped" {
		return true, nil
	}

	return false, nil
}

func resourceNifcloudCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	gatewayFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("customer-gateway-id"),
		RequestValue: []*string{nifcloud.String(d.Id())},
	}

	resp, err := conn.DescribeCustomerGateways(&computing.DescribeCustomerGatewaysInput{
		Filter: []*computing.RequestFilterStruct{gatewayFilter},
	})
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidCustomerGatewayID.NotFound" {
			d.SetId("")
			return nil
		}
		log.Printf("[ERROR] Error finding CustomerGateway: %s", err)
		return err
	}

	if len(resp.CustomerGatewaySet) != 1 {
		return fmt.Errorf("Error finding CustomerGateway: %s", d.Id())
	}

	if *resp.CustomerGatewaySet[0].State == "deleted" {
		log.Printf("[INFO] Customer Gateway is in `deleted` state: %s", d.Id())
		d.SetId("")
		return nil
	}

	customerGateway := resp.CustomerGatewaySet[0]
	d.Set("name", customerGateway.NiftyCustomerGatewayName)
	d.Set("ip_address", customerGateway.IpAddress)
	d.Set("description", customerGateway.NiftyCustomerGatewayDescription)
	d.Set("lan_side_cidr_block", customerGateway.NiftyLanSideCidrBlock)
	d.Set("lan_side_ip_address", customerGateway.NiftyLanSideIpAddress)

	return nil
}

func resourceNifcloudCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	d.Partial(true)

	log.Printf("[INFO] Updating CustomerGateway %s", d.Id())
	if d.HasChange("description") {
		input := computing.NiftyModifyCustomerGatewayAttributeInput{
			CustomerGatewayId: nifcloud.String(d.Id()),
			Attribute:         nifcloud.String("niftyCustomerGatewayDescription"),
			Value:             nifcloud.String(d.Get("description").(string)),
		}
		_, err := conn.NiftyModifyCustomerGatewayAttribute(&input)
		if err != nil {
			return fmt.Errorf("error description updating CustomerGateway (%s): %s", d.Id(), err)
		}	
	}
	if d.HasChange("name") {
		input := computing.NiftyModifyCustomerGatewayAttributeInput{
			CustomerGatewayId: nifcloud.String(d.Id()),
			Attribute:         nifcloud.String("niftyCustomerGatewayName"),
			Value:             nifcloud.String(d.Get("name").(string)),
		}
		_, err := conn.NiftyModifyCustomerGatewayAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating CustomerGateway (%s): %s", d.Id(), err)
		}	
	}

	d.Partial(false)

	return resourceNifcloudCustomerGatewayRead(d, meta)
}

func resourceNifcloudCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	_, err := conn.DeleteCustomerGateway(&computing.DeleteCustomerGatewayInput{
		CustomerGatewayId: nifcloud.String(d.Id()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.CustomerGatewayId" {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting CustomerGateway: %s", err)
	}

	gatewayFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("customer-gateway-id"),
		RequestValue: []*string{nifcloud.String(d.Id())},
	}

	input := &computing.DescribeCustomerGatewaysInput{
		Filter: []*computing.RequestFilterStruct{gatewayFilter},
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DescribeCustomerGateways(input)
		log.Printf("[DEBUG] delete after describe 001 %v", resp)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.CustomerGatewayId" {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		err = checkGatewayDeleteResponse(resp, d.Id())
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})

	if isResourceTimeoutError(err) {
		var resp *computing.DescribeCustomerGatewaysOutput
		resp, err = conn.DescribeCustomerGateways(input)
		log.Printf("[DEBUG] delete after describe 002 %v", resp)

		if err != nil {
			return checkGatewayDeleteResponse(resp, d.Id())
		}
	}

	if err != nil {
		return fmt.Errorf("Error deleting customer gateway: %s", err)
	}
	return nil

}

func checkGatewayDeleteResponse(resp *computing.DescribeCustomerGatewaysOutput, id string) error {
	if resp.CustomerGatewaySet == nil {
		return nil
	}

	switch *resp.CustomerGatewaySet[0].State {
	case "pending", "available":
		return fmt.Errorf("Gateway (%s) in state (%s), retrying", id, *resp.CustomerGatewaySet[0].State)
	case "stopped":
		return nil
	default:
		return fmt.Errorf("Unrecognized state (%s) for Customer Gateway delete on (%s)", *resp.CustomerGatewaySet[0].State, id)
	}
}
