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

func resourceNifcloudNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudNetworkCreate,
		Read:   resourceNifcloudNetworkRead,
		Update: resourceNifcloudNetworkUpdate,
		Delete: resourceNifcloudNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 15),
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"accounting_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNifcloudNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	createOpts := &computing.NiftyCreatePrivateLanInput{
		PrivateLanName:   nifcloud.String(d.Get("name").(string)),
		CidrBlock:        nifcloud.String(d.Get("cidr_block").(string)),
		AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		AccountingType:   nifcloud.String(d.Get("accounting_type").(string)),
		Description:      nifcloud.String(d.Get("description").(string)),
	}

	var err error
	resp, err := conn.NiftyCreatePrivateLan(createOpts)

	if err != nil {
		return fmt.Errorf("Error creating subnet: %s", err)
	}

	// Get the ID and store it
	subnet := resp.PrivateLan
	d.SetId(*subnet.NetworkId)
	log.Printf("[INFO] Subnet ID: %s", d.Id())

	// Wait for the Subnet to become available
	log.Printf("[DEBUG] Waiting for subnet (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Refresh: SubnetStateRefreshFunc(conn, d.Id()),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for subnet (%s) to become ready: %s",
			d.Id(), err)
	}

	return resourceNifcloudNetworkRead(d, meta)
}

func resourceNifcloudNetworkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	resp, err := conn.NiftyDescribePrivateLans(&computing.NiftyDescribePrivateLansInput{
		NetworkId: []*string{nifcloud.String(d.Id())},
	})

	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidParameterNotFound.NetworkId" {
			// Update state to indicate the subnet no longer exists.
			d.SetId("")
			return nil
		}
		return err
	}
	if resp == nil {
		return nil
	}

	subnet := resp.PrivateLanSet[0]

	d.Set("name", subnet.PrivateLanName)
	d.Set("cidr_block", subnet.CidrBlock)
	d.Set("availability_zone", subnet.AvailabilityZone)
	//d.Set("accounting_type", subnet.AccountingType)
	d.Set("accounting_type", subnet.NextMonthAccountingType)
	d.Set("description", subnet.Description)

	return nil
}

func resourceNifcloudNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	d.Partial(true)

	updateStateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    SubnetStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if d.HasChange("description") {
		_, err := conn.NiftyModifyPrivateLanAttribute(&computing.NiftyModifyPrivateLanAttributeInput{
			NetworkId: nifcloud.String(d.Id()),
			Attribute: nifcloud.String("description"),
			Value:     nifcloud.String(d.Get("description").(string)),
		})
		if err != nil {
			return fmt.Errorf("Error NiftyModifyPrivateLanAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("name") {
		_, err := conn.NiftyModifyPrivateLanAttribute(&computing.NiftyModifyPrivateLanAttributeInput{
			NetworkId: nifcloud.String(d.Id()),
			Attribute: nifcloud.String("privateLanName"),
			Value:     nifcloud.String(d.Get("name").(string)),
		})
		if err != nil {
			return fmt.Errorf("Error NiftyModifyPrivateLanAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("accounting_type") {
		_, err := conn.NiftyModifyPrivateLanAttribute(&computing.NiftyModifyPrivateLanAttributeInput{
			NetworkId: nifcloud.String(d.Id()),
			Attribute: nifcloud.String("accountingType"),
			Value:     nifcloud.String(d.Get("accounting_type").(string)),
		})
		if err != nil {
			return fmt.Errorf("Error NiftyModifyPrivateLanAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("cidr_block") {
		_, err := conn.NiftyModifyPrivateLanAttribute(&computing.NiftyModifyPrivateLanAttributeInput{
			NetworkId: nifcloud.String(d.Id()),
			Attribute: nifcloud.String("cidrBlock"),
			Value:     nifcloud.String(d.Get("cidr_block").(string)),
		})
		if err != nil {
			return fmt.Errorf("Error NiftyModifyPrivateLanAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceNifcloudNetworkRead(d, meta)
}

func resourceNifcloudNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf("[INFO] Deleting subnet: %s", d.Id())

	req := &computing.NiftyDeletePrivateLanInput{
		NetworkId: nifcloud.String(d.Id()),
	}

	wait := resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"destroyed"},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			_, err := conn.NiftyDeletePrivateLan(req)
			if err != nil {
				if apiErr, ok := err.(awserr.Error); ok {
					if apiErr.Code() == "DependencyViolation" {
						// There is some pending operation, so just retry
						// in a bit.
						return 42, "pending", nil
					}

					if apiErr.Code() == "InvalidParameterNotFound.NetworkId" {
						return 42, "destroyed", nil
					}
				}

				return 42, "failure", err
			}

			return 42, "destroyed", nil
		},
	}

	if _, err := wait.WaitForState(); err != nil {
		return fmt.Errorf("Error deleting subnet: %s", err)
	}

	return nil
}

// SubnetStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch a Subnet.
func SubnetStateRefreshFunc(conn *computing.Computing, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.NiftyDescribePrivateLans(&computing.NiftyDescribePrivateLansInput{
			NetworkId: []*string{nifcloud.String(id)},
		})
		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidParameterNotFound.NetworkId" {
				resp = nil
			} else {
				log.Printf("Error on SubnetStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		subnet := resp.PrivateLanSet[0]
		return subnet, *subnet.State, nil
	}
}
