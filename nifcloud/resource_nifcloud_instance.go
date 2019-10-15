package nifcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
	"log"
	"strconv"
	"time"
)

func resourceNifcloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudInstanceCreate,
		Read:   resourceNifcloudInstanceRead,
		Update: resourceNifcloudInstanceUpdate,
		Delete: resourceNifcloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				conn := meta.(*NifcloudClient).computingconn
				out, err := conn.DescribeInstances(&computing.DescribeInstancesInput{})

				if err != nil {
					return nil, fmt.Errorf("Error Import resource: %s", err)
				}

				for _, r := range out.ReservationSet {
					i := r.InstancesSet[0]
					if *i.InstanceId == d.Id() {
						d.Set("name", i.InstanceId)

						return []*schema.ResourceData{d}, nil
					}
				}

				return nil, fmt.Errorf("Error Import resource: %s", d.Id())
			},
		},

		SchemaVersion: 1,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 15),
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"admin"},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"disable_api_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"accounting_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2",
			},
			"admin": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_name"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_name"},
			},
//			"ip_type": {
//				Type:          schema.TypeString,
//				Optional:      true,
//				ConflictsWith: []string{"network_interfaces"},
//				//Default:       "static",
//			},
//			"public_ip": {
//				Type:          schema.TypeString,
//				Optional:      true,
//				ConflictsWith: []string{"network_interfaces"},
//			},
			"agreement": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_interfaces": {
				Type:          schema.TypeSet,
				Optional:      true,
				MinItems:      1,
				MaxItems:      2,
//				ConflictsWith: []string{"ip_type","public_ip"},
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
					},
				},
			},
			"license": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"license_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"license_num": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"user_data": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNifcloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
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
			}

			networkInterfaces = append(networkInterfaces, networkInterface)
		}
	}

	var licenses []*computing.RequestLicenseStruct
	if licensesSet, ok := d.GetOk("license"); ok {
		for _, l := range licensesSet.(*schema.Set).List() {
			license := &computing.RequestLicenseStruct{}
            if v, ok := l.(map[string]interface{}); ok {
		  		license.SetLicenseName(v["license_name"].(string))
				license.SetLicenseNum(v["license_num"].(string))
			}

			licenses = append(licenses, license)
		}
	}

	input := computing.RunInstancesInput{
		InstanceId:            nifcloud.String(d.Get("name").(string)),
		ImageId:               nifcloud.String(d.Get("image_id").(string)),
		KeyName:               nifcloud.String(d.Get("key_name").(string)),
		SecurityGroup:         securityGroups,
		UserData:              nifcloud.String(d.Get("user_data").(string)),
		InstanceType:          nifcloud.String(d.Get("instance_type").(string)),
		Placement:             &computing.RequestPlacementStruct{AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string))},
		DisableApiTermination: nifcloud.Bool(d.Get("disable_api_termination").(bool)),
		AccountingType:        nifcloud.String(d.Get("accounting_type").(string)),
		Admin:                 nifcloud.String(d.Get("admin").(string)),
		Password:              nifcloud.String(d.Get("password").(string)),
//		IpType:                nifcloud.String(d.Get("ip_type").(string)),
//		PublicIp:              nifcloud.String(d.Get("public_ip").(string)),
		Agreement:             nifcloud.Bool(d.Get("agreement").(bool)),
		Description:           nifcloud.String(d.Get("description").(string)),
		NetworkInterface:      networkInterfaces,
		License:               licenses,
	}

	out, err := conn.RunInstances(&input)
	if err != nil {
		return fmt.Errorf("Error RunInstancesInput: %s", err)
	}

	instance := out.InstancesSet[0]

	log.Printf("[INFO] Instance Id: %s", *instance.InstanceId)

	d.SetId(*instance.InstanceId)
	d.Set("unique_id", instance.InstanceUniqueId)

	log.Printf("[DEBUG] Waiting for instance (%s) to become running", *instance.InstanceId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"running","warning"},
		Refresh:    InstanceStateRefreshFunc(meta, d.Id(), []string{"terminated"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			*instance.InstanceId, err)
	}

	return resourceNifcloudInstanceRead(d, meta)
}

func resourceNifcloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	stopInstancesInput := computing.StopInstancesInput{
		InstanceId: []*string{nifcloud.String(d.Id())},
	}
	if _, err := conn.StopInstances(&stopInstancesInput); err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "Server.ProcessingFailure.Instance.Stop" {
			// 何もしないで継続
		} else {
			return fmt.Errorf("Error StopInstances: %s", err)
		}
	}

	log.Printf("[DEBUG] Waiting for instance (%s) to become stopped", d.Id())

	stopStateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running"},
		Target:     []string{"stopped"},
		Refresh:    InstanceStateRefreshFunc(meta, d.Id(), []string{"warning"}),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stopStateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to stopped: %s", d.Id(), err)
	}

	terminateInstancesInput := computing.TerminateInstancesInput{
		InstanceId: []*string{nifcloud.String(d.Id())},
	}
	if _, err := conn.TerminateInstances(&terminateInstancesInput); err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidInstanceID.NotFound" {
			return nil
		}
		return fmt.Errorf("Error terminating instance: %s", err)
	}

	log.Printf("[DEBUG] Waiting for instance (%s) to become terminate", d.Id())

	terminateStateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "stopped"},
		Target:     []string{"terminated"},
		Refresh:    InstanceStateRefreshFunc(meta, d.Id(), []string{"warning"}),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := terminateStateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to terminate: %s", d.Id(), err)
	}

	return nil
}

func resourceNifcloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	updateStateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"running", "stopped", "warning"},
		Refresh:    InstanceStateRefreshFunc(meta, d.Id(), []string{"waiting", "terminated"}),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	d.Partial(true)

	if d.HasChange("description") {
		_, err := conn.ModifyInstanceAttribute(&computing.ModifyInstanceAttributeInput{
			InstanceId: nifcloud.String(d.Id()),
			Attribute:  nifcloud.String("description"),
			Value:      nifcloud.String(d.Get("description").(string)),
		})
		if err != nil {
			return fmt.Errorf("Error ModifyInstanceAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("instance_type") {
		_, err := conn.ModifyInstanceAttribute(&computing.ModifyInstanceAttributeInput{
			InstanceId: nifcloud.String(d.Id()),
			Attribute:  nifcloud.String("instanceType"),
			Value:      nifcloud.String(d.Get("instance_type").(string)),
		})
		if err != nil {
			return fmt.Errorf("Error ModifyInstanceAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("disable_api_termination") {
		_, err := conn.ModifyInstanceAttribute(&computing.ModifyInstanceAttributeInput{
			InstanceId: nifcloud.String(d.Id()),
			Attribute:  nifcloud.String("disableApiTermination"),
			Value:      nifcloud.String(strconv.FormatBool(d.Get("disable_api_termination").(bool))),
		})
		if err != nil {
			return fmt.Errorf("Error ModifyInstanceAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("name") {
		_, err := conn.ModifyInstanceAttribute(&computing.ModifyInstanceAttributeInput{
			InstanceId:  nifcloud.String(d.Id()),
			Attribute:   nifcloud.String("instanceName"),
			NiftyReboot: nifcloud.String(d.Get("name").(string)),
		})

		if err != nil {
			return fmt.Errorf("Error ModifyInstanceAttribute: %s", err)
		}

		d.SetId(d.Get("name").(string))

		updateStateConf := &resource.StateChangeConf{
			Pending:    []string{"pending", "terminated"},
			Target:     []string{"running", "stopped", "warning"},
			Refresh:    InstanceStateRefreshFunc(meta, d.Id(), []string{"waiting"}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 5 * time.Second,
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("accounting_type") {
		_, err := conn.ModifyInstanceAttribute(&computing.ModifyInstanceAttributeInput{
			InstanceId: nifcloud.String(d.Id()),
			Attribute:  nifcloud.String("accountingType"),
			Value:      nifcloud.String(d.Get("accounting_type").(string)),
		})
		if err != nil {
			return fmt.Errorf("Error ModifyInstanceAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("security_groups") {
		_, err := conn.ModifyInstanceAttribute(&computing.ModifyInstanceAttributeInput{
			InstanceId: nifcloud.String(d.Id()),
			Attribute:  nifcloud.String("groupId"),
			Value:      nifcloud.String(d.Get("security_groups").([]interface{})[0].(string)),
		})
		if err != nil {
			return fmt.Errorf("Error ModifyInstanceAttribute: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				d.Id(), err)
		}
	}

//	if d.HasChange("ip_type") {
//		_, err := conn.ModifyInstanceAttribute(&computing.ModifyInstanceAttributeInput{
//			InstanceId: nifcloud.String(d.Id()),
//			Attribute:  nifcloud.String("ipType"),
//			Value:      nifcloud.String(d.Get("ip_type").(string)),
//		})
//		if err != nil {
//			return fmt.Errorf("Error ModifyInstanceAttribute: %s", err)
//		}
//
//		if _, err := updateStateConf.WaitForState(); err != nil {
//			return fmt.Errorf(
//				"Error waiting for instance (%s) to become ready: %s",
//				d.Id(), err)
//		}
//	}
	if d.HasChange("network_interfaces") {
		var networkInterfaces []*computing.RequestNetworkInterfaceStruct
		if interfaces, ok := d.GetOk("network_interfaces"); ok {
			for _, ni := range interfaces.(*schema.Set).List() {
				networkInterface := &computing.RequestNetworkInterfaceStruct{}
				if v, ok := ni.(map[string]interface{}); ok {
					networkInterface.SetNetworkId(v["network_id"].(string))
					networkInterface.SetNetworkName(v["network_name"].(string))
					networkInterface.SetIpAddress(v["ipaddress"].(string))
				}

				networkInterfaces = append(networkInterfaces, networkInterface)
			}
		}

		_, err := conn.NiftyUpdateInstanceNetworkInterfaces(&computing.NiftyUpdateInstanceNetworkInterfacesInput{
			InstanceId:       nifcloud.String(d.Id()),
			NetworkInterface: networkInterfaces,
			//NiftyReboot:      nifcloud.String("true"),
		})
		if err != nil {
			return fmt.Errorf("Error NiftyUpdateInstanceNetworkInterfaces: %s", err)
		}

		if _, err := updateStateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceNifcloudInstanceRead(d, meta)
}

func resourceNifcloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	input := computing.DescribeInstancesInput{}

	out, err := conn.DescribeInstances(&input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "Client.InvalidParameterNotFound.Instance" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find Instance resource: %s", err)
	}

	reservations := out.ReservationSet
	var reservation *computing.ReservationSetItem = nil
	for _, v := range reservations {
		if *v.InstancesSet[0].InstanceId == d.Id() {
			reservation = v
		}
	}

	if reservation == nil {
		return fmt.Errorf("Couldn't find Instance resource: %s", err)
	}

	return setInstanceResourceData(d, meta, reservation)
}

// InstanceStateRefreshFunc is function
func InstanceStateRefreshFunc(meta interface{}, instanceID string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn := meta.(*NifcloudClient).computingconn

		input := computing.DescribeInstancesInput{
			InstanceId: []*string{nifcloud.String(instanceID)},
		}

		out, err := conn.DescribeInstances(&input)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.Instance" {
				return "", "terminated", nil
			}

			log.Printf("Error on InstanceStateRefresh: %s", err)
			return nil, "", err
		}

		instance := out.ReservationSet[0].InstancesSet[0]
		state := *instance.InstanceState.Name

		for _, failState := range failStates {
			if state == failState {
				return instance, state, fmt.Errorf("Failed to reach target state. Reason: %s", state)
			}
		}

		return instance, state, nil
	}
}

func setInstanceResourceData(d *schema.ResourceData, meta interface{}, reservation *computing.ReservationSetItem) error {
	conn := meta.(*NifcloudClient).computingconn

	instance := reservation.InstancesSet[0]
	d.Set("name", instance.InstanceId)
	d.SetId(*instance.InstanceId)
	d.Set("unique_id", instance.InstanceUniqueId)

	outDisableAPITermination, err := conn.DescribeInstanceAttribute(&computing.DescribeInstanceAttributeInput{
		InstanceId: nifcloud.String(d.Id()),
		Attribute:  nifcloud.String("disableApiTermination"),
	})

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "Client.InvalidParameterNotFound.Instance" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Instance: %s", err)
	}

	outUserData, err := conn.DescribeInstanceAttribute(&computing.DescribeInstanceAttributeInput{
		InstanceId: nifcloud.String(d.Id()),
		Attribute:  nifcloud.String("userData"),
	})

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "Client.InvalidParameterNotFound.Instance" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Instance: %s", err)
	}

	log.Printf("[INFO] **********************************\n instance : %v\n ***************************", instance)
	log.Printf("[INFO] **********************************\n NetworkInterfaceSet : %v\n ***************************", instance.NetworkInterfaceSet)
	d.Set("image_id", instance.ImageId)
	d.Set("instance_type", instance.InstanceType)
	//d.Set("accounting_type", instance.AccountingType)
	d.Set("accounting_type", instance.NextMonthAccountingType)
	d.Set("description", instance.Description)
	d.Set("availability_zone", instance.Placement.AvailabilityZone)
	d.Set("user_data", outUserData.UserData)
//	d.Set("ip_type", instance.IpType)
	d.Set("ip_address", instance.IpAddress)

	d.Set("disable_api_termination", outDisableAPITermination.DisableApiTermination.Value)

	// only windows
	//d.Set("admin", instance.Admin)
	d.Set("admin", nifcloud.String(d.Get("admin").(string)))

	// only linux
	d.Set("key_name", instance.KeyName)

	sgs := make([]string, 0, len(reservation.GroupSet))
	for _, sg := range reservation.GroupSet {
		sgs = append(sgs, *sg.GroupId)
	}

	log.Printf("[DEBUG] Setting Security Group Ids: %#v", sgs)
	if err := d.Set("security_groups", sgs); err != nil {
		return err
	}

	return nil
}
