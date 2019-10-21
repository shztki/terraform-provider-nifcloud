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
	"time"
)

func resourceNifcloudVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudVolumeCreate,
		Read:   resourceNifcloudVolumeRead,
		Update: resourceNifcloudVolumeUpdate,
		Delete: resourceNifcloudVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 32),
			},
			"disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceNifcloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	request := &computing.CreateVolumeInput{
		Size:           nifcloud.Int64(int64(d.Get("size").(int))),
		VolumeId:       nifcloud.String(d.Get("name").(string)),
		DiskType:       nifcloud.String(d.Get("disk_type").(string)),
		InstanceId:     nifcloud.String(d.Get("instance_id").(string)),
		AccountingType: nifcloud.String(d.Get("accounting_type").(string)),
		Description:    nifcloud.String(d.Get("description").(string)),
	}

	log.Printf(
		"[DEBUG] EBS Volume create opts: %s", request)
	result, err := conn.CreateVolume(request)
	if err != nil {
		return fmt.Errorf("Error creating EC2 volume: %s", err)
	}

	log.Println("[DEBUG] Waiting for Volume to become available")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available", "in-use"},
		Refresh:    volumeStateRefreshFunc(conn, *result.VolumeId),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Volume (%s) to become available: %s",
			*result.VolumeId, err)
	}

	d.SetId(*result.VolumeId)

	return resourceNifcloudVolumeRead(d, meta)
}

func resourceNifcloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	if d.HasChange("description") {
		_, err := conn.ModifyVolumeAttribute(&computing.ModifyVolumeAttributeInput{
			VolumeId:  nifcloud.String(d.Id()),
			Attribute: nifcloud.String("description"),
			Value:     nifcloud.String(d.Get("description").(string)),
		})
		if err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"creating", "configuring"},
			Target:     []string{"available", "in-use"},
			Refresh:    volumeStateRefreshFunc(conn, d.Id()),
			Timeout:    5 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for Volume (%s) to become available: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("name") {
		_, err := conn.ModifyVolumeAttribute(&computing.ModifyVolumeAttributeInput{
			VolumeId:  nifcloud.String(d.Id()),
			Attribute: nifcloud.String("volumeName"),
			Value:     nifcloud.String(d.Get("name").(string)),
		})
		if err != nil {
			return err
		}

		d.SetId(d.Get("name").(string))

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"creating", "configuring"},
			Target:     []string{"available", "in-use"},
			Refresh:    volumeStateRefreshFunc(conn, d.Id()),
			Timeout:    5 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for Volume (%s) to become available: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("accounting_type") {
		_, err := conn.ModifyVolumeAttribute(&computing.ModifyVolumeAttributeInput{
			VolumeId:  nifcloud.String(d.Id()),
			Attribute: nifcloud.String("accountingType"),
			Value:     nifcloud.String(d.Get("accounting_type").(string)),
		})
		if err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"creating", "configuring"},
			Target:     []string{"available", "in-use"},
			Refresh:    volumeStateRefreshFunc(conn, d.Id()),
			Timeout:    5 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for Volume (%s) to become available: %s",
				d.Id(), err)
		}
	}

	return resourceNifcloudVolumeRead(d, meta)
}

// volumeStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a the state of a Volume. Returns successfully when volume is available
func volumeStateRefreshFunc(conn *computing.Computing, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeVolumes(&computing.DescribeVolumesInput{
			VolumeId: []*string{&volumeID},
		})

		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok {
				// Set this to nil as if we didn't find anything.
				log.Printf("Error on Volume State Refresh: message: \"%s\", code:\"%s\"", ec2err.Message(), ec2err.Code())
				resp = nil
				return nil, "", err
			} 

			log.Printf("Error on Volume State Refresh: %s", err)
			return nil, "", err
		}

		v := resp.VolumeSet[0]
		return v, *v.Status, nil
	}
}

func resourceNifcloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	request := &computing.DescribeVolumesInput{
		VolumeId: []*string{nifcloud.String(d.Id())},
	}

	response, err := conn.DescribeVolumes(request)
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidParameterNotFound.Volume" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading EC2 volume %s: %s", d.Id(), err)
	}

	if response == nil || len(response.VolumeSet) == 0 || response.VolumeSet[0] == nil {
		return fmt.Errorf("error reading EC2 Volume (%s): empty response", d.Id())
	}

	volume := response.VolumeSet[0]

	d.Set("name", volume.VolumeId)
	d.Set("size", volume.Size)
	d.Set("disk_type", voDiskTypes()[nifcloud.StringValue(volume.DiskType)])
	//d.Set("accounting_type", volume.AccountingType)
	d.Set("accounting_type", volume.NextMonthAccountingType)

	return nil
}

func resourceNifcloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	detach := computing.DetachVolumeInput{
		Agreement:  nifcloud.Bool(true),
		VolumeId:   nifcloud.String(d.Id()),
		InstanceId: nifcloud.String(d.Get("instance_id").(string)),
	}

	if _, err := conn.DetachVolume(&detach); err != nil {
		return fmt.Errorf("Error DetachVolumeInput: %s", err)
	}

	input := &computing.DeleteVolumeInput{
		VolumeId: nifcloud.String(d.Id()),
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteVolume(input)

		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "InvalidParameterNotFound.Volume" {
			return nil
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if isResourceTimeoutError(err) {
		_, err = conn.DeleteVolume(input)
	}

	if err != nil {
		return fmt.Errorf("error deleting EBS Volume (%s): %s", d.Id(), err)
	}

	describeInput := &computing.DescribeVolumesInput{
		VolumeId: []*string{nifcloud.String(d.Id())},
	}

	var output *computing.DescribeVolumesOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		output, err = conn.DescribeVolumes(describeInput)

		if err != nil {
			return resource.NonRetryableError(err)
		}

		for _, volume := range output.VolumeSet {
			if nifcloud.StringValue(volume.VolumeId) == d.Id() {
				state := nifcloud.StringValue(volume.Status)

				if state != "" {
					return resource.RetryableError(fmt.Errorf("EBS Volume (%s) still deleting", d.Id()))
				}

				return resource.NonRetryableError(fmt.Errorf("EBS Volume (%s) in unexpected state after deletion: %s", d.Id(), state))
			}
		}

		return nil
	})

	if isResourceTimeoutError(err) {
		output, err = conn.DescribeVolumes(describeInput)
	}

	if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "InvalidParameterNotFound.Volume" {
		return nil
	}

	for _, volume := range output.VolumeSet {
		if nifcloud.StringValue(volume.VolumeId) == d.Id() {
			return fmt.Errorf("EBS Volume (%s) in unexpected state after deletion: %s", d.Id(), nifcloud.StringValue(volume.Status))
		}
	}

	return nil
}

func voDiskTypes() map[string]string {
	return map[string]string{
		"Standard Storage":      "2",
		"High-Speed Storage A":  "3",
		"High-Speed Storage B":  "4",
		"Flash Storage":         "5",
	}
}
