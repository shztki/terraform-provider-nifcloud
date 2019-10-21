package nifcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
	"log"
)

func resourceNifcloudInstanceBackupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudInstanceBackupRuleCreate,
		Read:   resourceNifcloudInstanceBackupRuleRead,
		Update: resourceNifcloudInstanceBackupRuleUpdate,
		Delete: resourceNifcloudInstanceBackupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"backup_instance_max_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 15),
			},
			"time_slot_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1",
			},
			"instance_unique_id": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
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

func resourceNifcloudInstanceBackupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	var instanceUniqueIDS []*string
	if ius := d.Get("instance_unique_id").([]interface{}); ius != nil {
		for _, v := range ius {
			instanceUniqueIDS = append(instanceUniqueIDS, nifcloud.String(v.(string)))
		}
	}

	input := computing.CreateInstanceBackupRuleInput{
		BackupInstanceMaxCount: nifcloud.Int64(int64(d.Get("backup_instance_max_count").(int))),
		InstanceBackupRuleName: nifcloud.String(d.Get("name").(string)),
		TimeSlotId:             nifcloud.String(d.Get("time_slot_id").(string)),
		InstanceUniqueId:       instanceUniqueIDS,
		Description:            nifcloud.String(d.Get("description").(string)),
	}

	log.Printf("[INFO] Creating InstanceBackupRule: %v", input)
	out, err := conn.CreateInstanceBackupRule(&input)
	if err != nil {
		return fmt.Errorf("error creating InstanceBackupRule: %s", err)
	}

	rules := out.InstanceBackupRule
	log.Printf("[INFO] InstanceBackupRule ID: %s", *rules.InstanceBackupRuleId)
	d.SetId(*rules.InstanceBackupRuleId)

	return resourceNifcloudInstanceBackupRuleRead(d, meta)
}

func resourceNifcloudInstanceBackupRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf("[INFO] Reading InstanceBackupRule: %s", d.Id())
	input := computing.DescribeInstanceBackupRulesInput{
		InstanceBackupRuleId: []*string{nifcloud.String(d.Id())},
	}

	out, err := conn.DescribeInstanceBackupRules(&input)
	if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidParameterNotFound.InstanceBackupRuleId" {
		log.Printf("[WARN] InstanceBackupRule (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading InstanceBackupRule (%s): %s", d.Id(), err)
	}

	rules := out.InstanceBackupRulesSet[0]
	log.Printf("[INFO] Describe InstanceBackupRule: %v", rules)

	d.Set("name", rules.InstanceBackupRuleName)
	d.Set("backup_instance_max_count", rules.BackupInstanceMaxCount)
	d.Set("time_slot_id", rules.TimeSlotId)
	d.Set("instance_unique_id", rules.InstancesSet[0].InstanceUniqueId)
	d.Set("description", rules.Description)

	return nil
}

func resourceNifcloudInstanceBackupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf("[INFO] Updating InstanceBackupRule %s", d.Id())
	if d.HasChange("description") {
		input := computing.ModifyInstanceBackupRuleAttributeInput{
			InstanceBackupRuleId: nifcloud.String(d.Id()),
			Description:          nifcloud.String(d.Get("description").(string)),
		}
		_, err := conn.ModifyInstanceBackupRuleAttribute(&input)
		if err != nil {
			return fmt.Errorf("error description updating InstanceBackupRule (%s): %s", d.Id(), err)
		}	
	}
	if d.HasChange("time_slot_id") {
		input := computing.ModifyInstanceBackupRuleAttributeInput{
			InstanceBackupRuleId: nifcloud.String(d.Id()),
			TimeSlotId:           nifcloud.String(d.Get("time_slot_id").(string)),
		}
		_, err := conn.ModifyInstanceBackupRuleAttribute(&input)
		if err != nil {
			return fmt.Errorf("error time_slot_id updating InstanceBackupRule (%s): %s", d.Id(), err)
		}	
	}
	if d.HasChange("name") {
		input := computing.ModifyInstanceBackupRuleAttributeInput{
			InstanceBackupRuleId:   nifcloud.String(d.Id()),
			InstanceBackupRuleName: nifcloud.String(d.Get("name").(string)),
		}
		_, err := conn.ModifyInstanceBackupRuleAttribute(&input)
		if err != nil {
			return fmt.Errorf("error name updating InstanceBackupRule (%s): %s", d.Id(), err)
		}	
	}
	if d.HasChange("backup_instance_max_count") {
		input := computing.ModifyInstanceBackupRuleAttributeInput{
			InstanceBackupRuleId:   nifcloud.String(d.Id()),
			BackupInstanceMaxCount: nifcloud.Int64(int64(d.Get("backup_instance_max_count").(int))),
		}
		_, err := conn.ModifyInstanceBackupRuleAttribute(&input)
		if err != nil {
			return fmt.Errorf("error backup_instance_max_count updating InstanceBackupRule (%s): %s", d.Id(), err)
		}	
	}

	return resourceNifcloudInstanceBackupRuleRead(d, meta)
}

func resourceNifcloudInstanceBackupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf("[INFO] Deleting InstanceBackupRule: %s", d.Id())
	_, err := conn.DeleteInstanceBackupRule(&computing.DeleteInstanceBackupRuleInput{
		InstanceBackupRuleId: nifcloud.String(d.Id()),
	})
	if err != nil {
		return fmt.Errorf("error deleting InstanceBackupRule (%s): %s", d.Id(), err)
	}

	return nil
}
