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

func resourceNifcloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudSecurityGroupCreate,
		Read:   resourceNifcloudSecurityGroupRead,
		Update: resourceNifcloudSecurityGroupUpdate,
		Delete: resourceNifcloudSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				conn := meta.(*NifcloudClient).computingconn
				out, err := conn.DescribeSecurityGroups(&computing.DescribeSecurityGroupsInput{})

				if err != nil {
					return nil, fmt.Errorf("Error Import resource: %s", err)
				}

				for _, r := range out.SecurityGroupInfo {
					if *r.GroupName == d.Id() {
						d.Set("name", r.GroupName)

						return []*schema.ResourceData{d}, nil
					}
				}

				return nil, fmt.Errorf("Error Import resource: %s", d.Id())
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Required:      true,
				ValidateFunc:  validation.StringLenBetween(1, 15),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
//				ValidateFunc: validation.StringLenBetween(0, 40),
			},

			"rules": {
				Type:       schema.TypeSet,
				Optional:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"to_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"protocol": {
							Type:      schema.TypeString,
							Required:  true,
//							Default:  "TCP",
						},

						"cidr_blocks": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"security_groups": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"description": {
							Type:         schema.TypeString,
							Optional:     true,
						},

						"inout": {
							Type:         schema.TypeString,
							Required:     true,
						},
					},
				},
			},

		},
	}
}

func resourceNifcloudSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	input := computing.CreateSecurityGroupInput{
		GroupName:        nifcloud.String(d.Get("name").(string)),
		GroupDescription: nifcloud.String(d.Get("description").(string)),
	}

	securitygroup, err := conn.CreateSecurityGroup(&input)
	if err != nil {
		return fmt.Errorf("Error CreateSecurityGroupInput: %s", err)
	}

	log.Printf("[INFO] SecurityGroup RequestId: %s", *securitygroup.RequestId)

	d.SetId(d.Get("name").(string))

	log.Printf("[DEBUG] Waiting for (%s) to become running", *securitygroup.RequestId)

	resp, err := waitForSgToExist(conn, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
		"Error waiting for Security Group (%s) to become available: %s",
		d.Id(), err)
	}
	group := resp.(*computing.SecurityGroupInfoSetItem)
	log.Printf("[INFO] SecurityGroup info: %s", *group.GroupName)

	if group.GroupName != nil && *group.GroupName != "" {
		log.Printf("[DEBUG] Authorize default rule for Security Group for %s", d.Id())

		ipPermissions := setSecurityGroupRule(d.Get("rules"))
//		log.Printf("[INFO] **********************************\n ipPermissions : %v\n ***************************", ipPermissions)
		if ipPermissions != nil && len(ipPermissions) != 0 {
			req := computing.AuthorizeSecurityGroupIngressInput{
				GroupName: nifcloud.String(d.Id()),
				IpPermissions: ipPermissions,
			}
	
			if _, err = conn.AuthorizeSecurityGroupIngress(&req); err != nil {
				return fmt.Errorf(
					"Error authorizing default rule for Security Group (%s): %s",
					d.Id(), err)
			}
		}
	}

	return resourceNifcloudSecurityGroupRead(d, meta)
}

func resourceNifcloudSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	input := computing.DeleteSecurityGroupInput{
		GroupName: nifcloud.String(d.Id()),
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.DeleteSecurityGroup(&input)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroup" {
				return nil
			}
			resource.RetryableError(err)
		}
		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = conn.DeleteSecurityGroup(&input)
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroup" {
			return nil
		}
	}
	if err != nil {
		return fmt.Errorf("Error deleting security group: %s", err)
	}

	return nil
}

func resourceNifcloudSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	if d.HasChange("description") {
		_, err := conn.UpdateSecurityGroup(&computing.UpdateSecurityGroupInput{
			GroupName:              nifcloud.String(d.Get("name").(string)),
			GroupDescriptionUpdate: nifcloud.String(d.Get("description").(string)),
		})
		if err != nil {
			return fmt.Errorf("Error UpdateSecurityGroup: %s", err)
		}
	}

	if d.HasChange("name") {
        before, after := d.GetChange("name")
		_, err := conn.UpdateSecurityGroup(&computing.UpdateSecurityGroupInput{
			GroupName:       nifcloud.String(before.(string)),
			GroupNameUpdate: nifcloud.String(after.(string)),
		})

		d.SetId(d.Get("name").(string))

		if err != nil {
			return fmt.Errorf("Error UpdateSecurityGroup: %s", err)
		}
	}

	if d.HasChange("rules") {
		before, after := d.GetChange("rules")
		ipPermissionsOld := setSecurityGroupRule(before)
//		log.Printf("[INFO] **********************************\n before ipPermissions : %v\n ***************************", ipPermissionsOld)
		ipPermissionsNew := setSecurityGroupRule(after)
//		log.Printf("[INFO] **********************************\n after ipPermissions : %v\n ***************************", ipPermissionsNew)
		if ipPermissionsOld != nil {
			req := computing.RevokeSecurityGroupIngressInput{
				GroupName: nifcloud.String(d.Id()),
				IpPermissions: ipPermissionsOld,
			}
	
			if _, err := conn.RevokeSecurityGroupIngress(&req); err != nil {
				return fmt.Errorf(
					"Error deleting default ingress rule for Security Group (%s): %s",
					d.Id(), err)
			}
		}
		if ipPermissionsNew != nil {
			req := computing.AuthorizeSecurityGroupIngressInput{
				GroupName: nifcloud.String(d.Id()),
				IpPermissions: ipPermissionsNew,
			}
	
			if _, err := conn.AuthorizeSecurityGroupIngress(&req); err != nil {
				return fmt.Errorf(
					"Error authorizing default ingress rule for Security Group (%s): %s",
					d.Id(), err)
			}
		}
	}

	return resourceNifcloudSecurityGroupRead(d, meta)
}

func resourceNifcloudSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	input := computing.DescribeSecurityGroupsInput{
		GroupName: []*string{nifcloud.String(d.Get("name").(string))},
	}

	out, err := conn.DescribeSecurityGroups(&input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroup" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SecurityGroup resource: %s", err)
	}

	return setSecurityGroupResourceData(d, meta, out)
}

func waitForSgToExist(conn *computing.Computing, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("[DEBUG] Waiting for Security Group (%s) to exist", id)
	stateConf := &resource.StateChangeConf{
		Pending: []string{""},
		Target:  []string{"exists"},
		Refresh: SGStateRefreshFunc(conn, id),
		Timeout: timeout,
	}

	return stateConf.WaitForState()
}

// SGStateRefreshFunc is function
func SGStateRefreshFunc(conn *computing.Computing, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := computing.DescribeSecurityGroupsInput{
			GroupName: []*string{nifcloud.String(id)},
		}
		resp, err := conn.DescribeSecurityGroups(&req)
		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok {
				if ec2err.Code() == "Client.InvalidParameterNotFound.SecurityGroup" {
					resp = nil
					err = nil
				}
			}

			if err != nil {
				log.Printf("Error on SGStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			return nil, "", nil
		}

		group := resp.SecurityGroupInfo[0]
		return group, "exists", nil
	}
}

func setSecurityGroupResourceData(d *schema.ResourceData, meta interface{}, out *computing.DescribeSecurityGroupsOutput) error {
	securitygroup := out.SecurityGroupInfo[0]

	d.Set("name", securitygroup.GroupName)
	d.Set("description", securitygroup.GroupDescription)
	d.Set("rules", securitygroup.IpPermissions)

	return nil
}

func setSecurityGroupRule(permissions interface{}) []*computing.RequestIpPermissionsStruct {
	var ipPermissions []*computing.RequestIpPermissionsStruct
	for _, ip := range permissions.(*schema.Set).List() {
		ipPermission := &computing.RequestIpPermissionsStruct{}
		if v, ok := ip.(map[string]interface{}); ok {
			if v["description"].(string) != "" {
				ipPermission.SetDescription(v["description"].(string))
			}
			
			if v["from_port"].(int) > 0 {
				ipPermission.SetFromPort(int64(v["from_port"].(int)))
				ipPermission.SetToPort(int64(v["to_port"].(int)))
			}
			
			if v["cidr_blocks"].(string) != "" {
				tmp := []*computing.RequestIpRangesStruct {
					{
						CidrIp: nifcloud.String(v["cidr_blocks"].(string)),
					},
				}
				ipPermission.SetRequestIpRanges(tmp)
			}
			
			if v["security_groups"].(string) != "" {
				tmp := []*computing.RequestGroupsStruct {
					{
						GroupName: nifcloud.String(v["security_groups"].(string)),
					},
				}
				ipPermission.SetRequestGroups(tmp)
			}

			if v["protocol"].(string) != "" {
				ipPermission.SetIpProtocol(v["protocol"].(string))
			}

			if v["inout"].(string) != "" {
				ipPermission.SetInOut(v["inout"].(string))
			}
		}
		ipPermissions = append(ipPermissions, ipPermission)
	}

	return ipPermissions
}
