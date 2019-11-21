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
)

func resourceNifcloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudEipCreate,
		Read:   resourceNifcloudEipRead,
		Update: resourceNifcloudEipUpdate,
		Delete: resourceNifcloudEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:     schema.TypeString,
				Optional: true,
//				Computed: true,
			},

//			"network_interface": {
//				Type:     schema.TypeString,
//				Optional: true,
//				Computed: true,
//			},

//			"allocation_id": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},

//			"association_id": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},

			"nifty_private_ip": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

		},
	}
}

func resourceNifcloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	allocOpts := &computing.AllocateAddressInput{
		NiftyPrivateIp: nifcloud.Bool(d.Get("nifty_private_ip").(bool)),
		Placement:      &computing.RequestPlacementStruct{AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string))},
	}

	log.Printf("[DEBUG] EIP create configuration: %#v", allocOpts)
	allocResp, err := conn.AllocateAddress(allocOpts)
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}

	// Assign the eips (unique) allocation id for use later
	// the EIP api has a conditional unique ID (really), so
	// if we're in a VPC we need to save the ID as such, otherwise
	// it defaults to using the public IP
	log.Printf("[DEBUG] EIP Allocate: %#v", allocResp)
	if d.Get("nifty_private_ip").(bool) {
		d.SetId(*allocResp.PrivateIpAddress)
	} else {
		d.SetId(*allocResp.PublicIp)
	}

	return resourceNifcloudEipUpdate(d, meta)
}

func resourceNifcloudEipRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	id := d.Id()

	req := &computing.DescribeAddressesInput{}

	if d.Get("nifty_private_ip").(bool) {
		req.PrivateIpAddress = []*string{nifcloud.String(id)}
	} else {
		req.PublicIp = []*string{nifcloud.String(id)}
	}

	log.Printf(
		"[DEBUG] EIP describe configuration: %s", req)

	var err error
	var describeAddresses *computing.DescribeAddressesOutput

	if d.IsNewResource() {
		err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			describeAddresses, err = conn.DescribeAddresses(req)
			if err != nil {
				awsErr, ok := err.(awserr.Error)
				if ok && (awsErr.Code() == "Client.InvalidAssociationId.NotFound" ||
					awsErr.Code() == "Client.InvalidParameterNotFound.IpAddress") {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}
			return nil
		})
		if isResourceTimeoutError(err) {
			describeAddresses, err = conn.DescribeAddresses(req)
		}
		if err != nil {
			return fmt.Errorf("Error retrieving EIP: %s", err)
		}
	} else {
		describeAddresses, err = conn.DescribeAddresses(req)
		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok && (awsErr.Code() == "Client.InvalidAssociationId.NotFound" ||
				awsErr.Code() == "Client.InvalidParameterNotFound.IpAddress") {
				log.Printf("[WARN] EIP not found, removing from state: %s", req)
				d.SetId("")
				return nil
			}
			return err
		}
	}

	var address *computing.AddressesSetItem

	// In the case that AWS returns more EIPs than we intend it to, we loop
	// over the returned addresses to see if it's in the list of results
	for _, addr := range describeAddresses.AddressesSet {
		if nifcloud.StringValue(addr.PrivateIpAddress) == id || nifcloud.StringValue(addr.PublicIp) == id {
			address = addr
			break
		}
	}

	if address == nil {
		log.Printf("[WARN] EIP %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	//if address.InstanceId != nil {
	log.Printf("[DEBUG] EIP InstanceID: %#v", address.InstanceId)
	if v, ok := d.GetOk("instance"); ok && v != "" {
		d.Set("instance", address.InstanceId)
	} else {
		d.Set("instance", "")
	}

	d.Set("private_ip", address.PrivateIpAddress)
	d.Set("public_ip", address.PublicIp)

	return nil
}

func resourceNifcloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	if d.HasChange("description") {
		input := computing.NiftyModifyAddressAttributeInput{
			Attribute: nifcloud.String("description"),
			Value:     nifcloud.String(d.Get("description").(string)),
		}

		if d.Get("nifty_private_ip").(bool) {
			input.PrivateIpAddress = nifcloud.String(d.Id())
		} else {
			input.PublicIp = nifcloud.String(d.Id())
		}

		_, err := conn.NiftyModifyAddressAttribute(&input)
		if err != nil {
			return fmt.Errorf("error description updating EIP (%s): %s", d.Id(), err)
		}
	}

	// If we are updating an EIP that is not newly created, and we are attached to
	// an instance or interface, detach first.
	disassociate := false
	if !d.IsNewResource() {
		o, _ := d.GetChange("instance")
		if d.HasChange("instance") && o.(string) != "" {
			disassociate = true
		}
	}
	if disassociate {
		if err := disassociateEip(d, meta); err != nil {
			return err
		}
		if err := waitForInstanceID(d, meta, conn, ""); err != nil {
			return err
		}
	}

	// Associate to instance or interface if specified
	associate := false
	v_instance, ok_instance := d.GetOk("instance")

	if d.HasChange("instance") && ok_instance {
		associate = true
	}
	if associate {
		instanceId := v_instance.(string)

		assocOpts := &computing.AssociateAddressInput{
			NiftyReboot: nifcloud.String("true"),
			InstanceId:  nifcloud.String(instanceId),
		}

		if d.Get("nifty_private_ip").(bool) {
			assocOpts.PrivateIpAddress = nifcloud.String(d.Id())
		} else {
			assocOpts.PublicIp = nifcloud.String(d.Id())
		}

		log.Printf("[DEBUG] EIP associate configuration: %s", assocOpts)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := conn.AssociateAddress(assocOpts)
			if err != nil {
				if isNifcloudErr(err, "Server.ResourceIncorrectState.IpAddress.Processing", "") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if isResourceTimeoutError(err) {
			_, err = conn.AssociateAddress(assocOpts)
		}
		if err != nil {
			// Prevent saving instance if association failed
			// e.g. missing internet gateway in VPC
			d.Set("instance", "")
			return fmt.Errorf("Failure associating EIP: %s", err)
		}

		if err := waitForInstanceID(d, meta, conn, instanceId); err != nil {
			return err
		}
	}

	return resourceNifcloudEipRead(d, meta)
}

func resourceNifcloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	if err := resourceNifcloudEipRead(d, meta); err != nil {
		return err
	}
	if d.Id() == "" {
		// This might happen from the read
		return nil
	}

	// If we are attached to an instance or interface, detach first.
	if d.Get("instance").(string) != "" {
		if err := disassociateEip(d, meta); err != nil {
			return err
		}
	}

	var input *computing.ReleaseAddressInput
	if d.Get("nifty_private_ip").(bool) {
		input = &computing.ReleaseAddressInput{
			PrivateIpAddress: nifcloud.String(d.Id()),
		}
	} else {
		input = &computing.ReleaseAddressInput{
			PublicIp: nifcloud.String(d.Id()),
		}
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error
		_, err = conn.ReleaseAddress(input)

		if err == nil {
			return nil
		}
		if _, ok := err.(awserr.Error); !ok {
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(err)
	})
	if isResourceTimeoutError(err) {
		_, err = conn.ReleaseAddress(input)
	}
	if err != nil {
		return fmt.Errorf("Error releasing EIP address: %s", err)
	}
	return nil
}

func disassociateEip(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn
	log.Printf("[DEBUG] Disassociating EIP: %s", d.Id())
	var err error
	disAssociateOpts := &computing.DisassociateAddressInput{
		NiftyReboot: nifcloud.String("true"),
	}

	if d.Get("nifty_private_ip").(bool) {
		disAssociateOpts.PrivateIpAddress = nifcloud.String(d.Id())
	} else {
		disAssociateOpts.PublicIp = nifcloud.String(d.Id())
	}
	_, err = conn.DisassociateAddress(disAssociateOpts)
	
	// First check if the association ID is not found. If this
	// is the case, then it was already disassociated somehow,
	// and that is okay. The most commmon reason for this is that
	// the instance or ENI it was attached it was destroyed.
	if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidAssociationID.NotFound" {
		err = nil
	}
	return err
}

func waitForInstanceID(d *schema.ResourceData, meta interface{}, conn *computing.Computing, instanceID string) error {
	id := d.Id()

	req := &computing.DescribeAddressesInput{}

	if d.Get("nifty_private_ip").(bool) {
		req.PrivateIpAddress = []*string{nifcloud.String(id)}
	} else {
		req.PublicIp = []*string{nifcloud.String(id)}
	}

	log.Printf("[DEBUG] EIP describe waitForInstanceID: %s", req)

	var err error
	var describeAddresses *computing.DescribeAddressesOutput

	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		describeAddresses, err = conn.DescribeAddresses(req)
		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok && (awsErr.Code() == "Client.InvalidAssociationId.NotFound" ||
				awsErr.Code() == "Client.InvalidParameterNotFound.IpAddress") {
				log.Printf("[DEBUG] EIP describe NonRetryableError: %s", req)
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}

		var address *computing.AddressesSetItem
		for _, addr := range describeAddresses.AddressesSet {
			if nifcloud.StringValue(addr.PrivateIpAddress) == id || nifcloud.StringValue(addr.PublicIp) == id {
				address = addr
				break
			}
		}
		if address == nil {
			return resource.RetryableError(fmt.Errorf("Error DescribeAddresses, retrying"))
		}
		if instanceID == "" {
			if address.InstanceId != nil {
				return resource.RetryableError(fmt.Errorf("Error DescribeAddresses InstanceID Exists, retrying"))
			}
		} else {
			if address.InstanceId == nil {
				return resource.RetryableError(fmt.Errorf("Error DescribeAddresses InstanceID NotExists, retrying"))
			}
		}
		err = nil
		return nil
	})
	if isResourceTimeoutError(err) {
		describeAddresses, err = conn.DescribeAddresses(req)
	}
	if err != nil {
		return fmt.Errorf("Error retrieving EIP: %s", err)
	}

	log.Printf("[DEBUG] EIP describe waitForInstance after: %v", describeAddresses)
	return err
}