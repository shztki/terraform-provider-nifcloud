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

func resourceNifcloudImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudImageCreate,
		Read:   resourceNifcloudImageRead,
		Update: resourceNifcloudImageUpdate,
		Delete: resourceNifcloudImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"region_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringLenBetween(1, 40),
			},
			"left_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNifcloudImageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	req := computing.CreateImageInput{
		LeftInstance: nifcloud.Bool(d.Get("left_instance").(bool)),
		Name:         nifcloud.String(d.Get("name").(string)),
		Placement:    &computing.RequestPlacementStruct{RegionName: nifcloud.String(d.Get("region_name").(string)), AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string))},
		InstanceId:   nifcloud.String(d.Get("instance_id").(string)),
		Description:  nifcloud.String(d.Get("description").(string)),
	}

	log.Printf("[INFO] Creating Image: %v", req)
	res, err := conn.CreateImage(&req)
	if err != nil {
		return fmt.Errorf("error creating Image: %s", err)
	}

	log.Printf("[INFO] Image ID: %s", *res.ImageId)
	d.SetId(*res.ImageId)

	return resourceNifcloudImageRead(d, meta)
}

func resourceNifcloudImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf("[INFO] Reading Image: %s", d.Id())
	req := computing.DescribeImagesInput{
		ImageId: []*string{nifcloud.String(d.Id())},
	}

	var res *computing.DescribeImagesOutput
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		var err error
		res, err = conn.DescribeImages(&req)
		if err != nil {
			if err, ok := err.(awserr.Error); ok && err.Code() == "Client.InvalidParameterNotFound.ImageId" {
				if d.IsNewResource() {
					return resource.RetryableError(err)
				}

				log.Printf("[DEBUG] %s no longer exists, so we'll drop it from the state", d.Id())
				d.SetId("")
				return nil
			}

			return resource.NonRetryableError(err)
		}
		return nil
	})
	if isResourceTimeoutError(err) {
		res, err = conn.DescribeImages(&req)
	}
	if err != nil {
		return fmt.Errorf("Unable to find Image after retries: %s", err)
	}

	if len(res.ImagesSet) != 1 {
		d.SetId("")
		return nil
	}
	
	image := res.ImagesSet[0]
//	state := *image.ImageState
	log.Printf("[INFO] Describe Image: %v", image)

//	if state == "deregistered" {
//		d.SetId("")
//		return nil
//	}
//
//	if state != "available" {
//		return fmt.Errorf("AMI has become %s", state)
//	}
	
	d.Set("name", image.Name)
	d.Set("description", image.Description)

	return nil
}

func resourceNifcloudImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	d.Partial(true)

	log.Printf("[INFO] Updating Image %s", d.Id())
	if d.HasChange("description") {
		input := computing.ModifyImageAttributeInput{
			ImageId:   nifcloud.String(d.Id()),
			Attribute: nifcloud.String("description"),
			Value:     nifcloud.String(d.Get("description").(string)),
		}
		_, err := conn.ModifyImageAttribute(&input)
		if err != nil {
			return fmt.Errorf("error description updating Image (%s): %s", d.Id(), err)
		}	
	}
	if d.HasChange("name") {
		input := computing.ModifyImageAttributeInput{
			ImageId:   nifcloud.String(d.Id()),
			Attribute: nifcloud.String("imageName"),
			Value:     nifcloud.String(d.Get("name").(string)),
		}
		_, err := conn.ModifyImageAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating Image (%s): %s", d.Id(), err)
		}	
	}

	d.Partial(false)

	return resourceNifcloudImageRead(d, meta)
}

func resourceNifcloudImageDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf("[INFO] Deleting Image: %s", d.Id())
	_, err := conn.DeleteImage(&computing.DeleteImageInput{
		ImageId: nifcloud.String(d.Id()),
	})
	if err != nil {
		return fmt.Errorf("error deleting Image (%s): %s", d.Id(), err)
	}

	return nil
}
