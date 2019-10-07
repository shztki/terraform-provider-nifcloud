package nifcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
	"log"
	"time"
	"strconv"
//	"encoding/json"
)

func resourceNifcloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudSecurityGroupCreate,
		Read:   resourceNifcloudSecurityGroupRead,
		Update: resourceNifcloudSecurityGroupUpdate,
		Delete: resourceNifcloudSecurityGroupDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

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

			"ingress": {
				Type:       schema.TypeSet,
				Optional:   true,
//				Computed:   false,
//				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeString,
                            Optional: true,
						},

						"to_port": {
							Type:     schema.TypeString,
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
//							Set:      schema.HashString,
						},

						"description": {
							Type:         schema.TypeString,
							Optional:     true,
						},
					},
				},
//				Set: resourceAwsSecurityGroupRuleHash,
			},
//
//			"egress": {
//				Type:       schema.TypeSet,
//				Optional:   true,
//				Computed:   true,
////				ConfigMode: schema.SchemaConfigModeAttr,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"from_port": {
//							Type:     schema.TypeInt,
//							Optional: true,
//						},
//
//						"to_port": {
//							Type:     schema.TypeInt,
//							Optional: true,
//						},
//
//						"protocol": {
//							Type:      schema.TypeString,
//							Required:  true,
//							Default:   "ANY",
//						},
//
//						"cidr_blocks": {
//							Type:     schema.TypeList,
//							Optional: true,
//							Elem:     &schema.Schema{Type: schema.TypeString},
//						},
//
//						"security_groups": {
//							Type:     schema.TypeSet,
//							Optional: true,
//							Elem:     &schema.Schema{Type: schema.TypeString},
////							Set:      schema.HashString,
//						},
//
//						"description": {
//							Type:         schema.TypeString,
//							Optional:     true,
////							ValidateFunc: validateSecurityGroupRuleDescription,
//						},
//					},
//				},
////				Set: resourceAwsSecurityGroupRuleHash,
//			},

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
		log.Printf("[DEBUG] Authorize default ingress rule for Security Group for %s", d.Id())

		ipPermissions := setSecurityGroupIngress(d, meta)
		log.Printf("[INFO] **********************************\n main ipPermissions : %v\n ***************************", ipPermissions)
		if ipPermissions != nil {
			req := computing.AuthorizeSecurityGroupIngressInput{
				GroupName: nifcloud.String(d.Id()),
				IpPermissions: ipPermissions,
			}
	
			if _, err = conn.AuthorizeSecurityGroupIngress(&req); err != nil {
				return fmt.Errorf(
					"Error authorizing default ingress rule for Security Group (%s): %s",
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
			if isAWSErr(err, "InvalidGroup.NotFound", "") {
				return nil
			}
			if isAWSErr(err, "DependencyViolation", "") {
				// If it is a dependency violation, we want to retry
				return resource.RetryableError(err)
			}
			resource.NonRetryableError(err)
		}
		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = conn.DeleteSecurityGroup(&input)
		if isAWSErr(err, "InvalidGroup.NotFound", "") {
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

		if err != nil {
			return fmt.Errorf("Error UpdateSecurityGroup: %s", err)
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
		if ok && awsErr.Code() == "Client.InvalidParameterNotFound.GroupName" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find Instance resource: %s", err)
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
				if ec2err.Code() == "InvalidSecurityGroupID.NotFound" ||
					ec2err.Code() == "InvalidGroup.NotFound" {
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
	d.Set("ingress", securitygroup.IpPermissions)

	return nil
}

func setSecurityGroupIngress(d *schema.ResourceData, meta interface{}) []*computing.RequestIpPermissionsStruct {
	var ipPermissions []*computing.RequestIpPermissionsStruct
	if permissions, ok := d.GetOk("ingress"); ok {
		for _, ip := range permissions.(*schema.Set).List() {
			ipPermission := &computing.RequestIpPermissionsStruct{}
			if v, ok := ip.(map[string]interface{}); ok {
				if v["description"].(string) != "" {
					ipPermission.SetDescription(v["description"].(string))
				}
				
				if v["from_port"].(string) != "" {
					var from64,to64 int64
					from64, _ = strconv.ParseInt(v["from_port"].(string),10,64)
					to64, _ = strconv.ParseInt(v["to_port"].(string),10,64)
					ipPermission.SetFromPort(from64)
					ipPermission.SetToPort(to64)
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

				ipPermission.SetInOut("IN")
			}
			ipPermissions = append(ipPermissions, ipPermission)
		}
	}

	return ipPermissions
}
