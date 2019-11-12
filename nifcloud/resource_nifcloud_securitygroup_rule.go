package nifcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
	"log"
	"time"
	"bytes"
	"strings"
)

func resourceNifcloudSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudSecurityGroupRuleCreate,
		Read:   resourceNifcloudSecurityGroupRuleRead,
		Update: resourceNifcloudSecurityGroupRuleUpdate,
		Delete: resourceNifcloudSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Required:      true,
				ValidateFunc:  validation.StringLenBetween(1, 15),
			},

			"rules": {
				Type:       schema.TypeSet,
				Optional:   true,
				MinItems:   1,
				MaxItems:   1,
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

func resourceNifcloudSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn
	sgID := d.Get("name").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, err := findResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}
	if sg == nil {
		return err
	}

	perm := expandIPPerm(d.Get("rules"))
	if perm == nil {
		return nil
	}

//	log.Printf("[DEBUG] Authorize rule for Security Group for %s", sgID)
//	log.Printf("[INFO] **********************************\n ipPermissions : %v\n ***************************", perm)
	req := computing.AuthorizeSecurityGroupIngressInput{
		GroupName:     nifcloud.String(sgID),
		IpPermissions: perm,
	}

	if _, err = conn.AuthorizeSecurityGroupIngress(&req); err != nil {
		return fmt.Errorf(
			"Error authorizing rule for Security Group (%s): %s",
			sgID, err)
	}

	var rules []*computing.IpPermissionsSetItem
	id := ipPermissionIDHash(sgID, perm[0])
	log.Printf("[DEBUG] Computed group rule ID %s", id)

	d.SetId(id)

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		sg, err := findResourceSecurityGroup(conn, sgID)

		if err != nil {
			log.Printf("[DEBUG] Error finding Security Group (%s) for Rule (%s): %s", sgID, id, err)
			return resource.NonRetryableError(err)
		}

		rules = sg.IpPermissions
	
		rule := findRuleMatch(perm[0], rules)
		if rule == nil {
			log.Printf("[DEBUG] Unable to find matching Security Group Rule (%s) for Group %s",
				id, sgID)
			return resource.RetryableError(fmt.Errorf("No match found"))
		}

		log.Printf("[DEBUG] Found rule for Security Group Rule (%s): %s", id, rule)
		return nil
	})
	if isResourceTimeoutError(err) {
		sg, err := findResourceSecurityGroup(conn, sgID)
		if err != nil {
			return fmt.Errorf("Error finding security group: %s", err)
		}

		rules = sg.IpPermissions

		rule := findRuleMatch(perm[0], rules)
		if rule == nil {
			return fmt.Errorf("Error finding matching security group rule: %s", err)
		}
	}
	if err != nil {
		return fmt.Errorf("Error finding matching Security Group Rule (%s) for Group %s", id, sgID)
	}
	
	return nil
}

func resourceNifcloudSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn
	sgID := d.Get("name").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	sg, err := findResourceSecurityGroup(conn, sgID)
	if err != nil {
		return err
	}
	perm := expandIPPerm(d.Get("rules"))
	if perm == nil {
		return err
	}

	req := computing.RevokeSecurityGroupIngressInput{
		GroupName: nifcloud.String(sgID),
		IpPermissions: perm,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.RevokeSecurityGroupIngress(&req)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroupIngress" {
				return nil
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = conn.RevokeSecurityGroupIngress(&req)
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroupIngress" {
			return nil
		}
	}
	if err != nil {
		return fmt.Errorf(
			"Error deleting rule for Security Group (%s): %s",
			*sg.GroupName, err)
	}

	return nil
}

func resourceNifcloudSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn
	sgID := d.Get("name").(string)

	awsMutexKV.Lock(sgID)
	defer awsMutexKV.Unlock(sgID)

	if d.HasChange("rules") {
		before, after := d.GetChange("rules")
		ipPermissionsOld := expandIPPerm(before)
		log.Printf("[INFO] **********************************\n before ipPermissions : %v\n ***************************", ipPermissionsOld)
		ipPermissionsNew := expandIPPerm(after)
		log.Printf("[INFO] **********************************\n after ipPermissions : %v\n ***************************", ipPermissionsNew)
		if ipPermissionsOld != nil {
			req := computing.RevokeSecurityGroupIngressInput{
				GroupName: nifcloud.String(sgID),
				IpPermissions: ipPermissionsOld,
			}
	
			if _, err := conn.RevokeSecurityGroupIngress(&req); err != nil {
				return fmt.Errorf(
					"Error deleting rule for Security Group (%s): %s",
					d.Id(), err)
			}
		}
		if ipPermissionsNew != nil {
			req := computing.AuthorizeSecurityGroupIngressInput{
				GroupName: nifcloud.String(sgID),
				IpPermissions: ipPermissionsNew,
			}
	
			if _, err := conn.AuthorizeSecurityGroupIngress(&req); err != nil {
				return fmt.Errorf(
					"Error authorizing rule for Security Group (%s): %s",
					d.Id(), err)
			}
		}
	}

	return resourceNifcloudSecurityGroupRuleRead(d, meta)
}

func resourceNifcloudSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn
	sgID := d.Get("name").(string)
	sg, err := findResourceSecurityGroup(conn, sgID)
	if err != nil {
		return fmt.Errorf("Error finding security group (%s) for rule (%s): %s", sgID, d.Id(), err)
	}

	var rule *computing.IpPermissionsSetItem
	var rules []*computing.IpPermissionsSetItem
	rules = sg.IpPermissions
//	log.Printf("[DEBUG] Rules %v", rules)

	perm := expandIPPerm(d.Get("rules"))
	if perm == nil {
		return err
	}

	if len(rules) == 0 {
		log.Printf("[WARN] No rules were found for Security Group (%s) looking for Security Group Rule (%s)",
			sgID, d.Id())
		d.SetId("")
		return nil
	}

	rule = findRuleMatch(perm[0], rules)

	if rule == nil {
		log.Printf("[DEBUG] Unable to find matching Security Group Rule (%s) for Group %s",
			d.Id(), sgID)
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Found rule for Security Group Rule (%s): %s", d.Id(), rule)

	d.Set("name", sg.GroupName)
	d.Set("rules", d.Get("rules"))

	if strings.Contains(d.Id(), "_") {
		// import so fix the id
		id := ipPermissionIDHash(sgID, perm[0])
		d.SetId(id)
	}

	return nil
}

func findResourceSecurityGroup(conn *computing.Computing, id string) (*computing.SecurityGroupInfoSetItem, error) {
	input := computing.DescribeSecurityGroupsInput{
		GroupName: []*string{nifcloud.String(id)},
	}

	out, err := conn.DescribeSecurityGroups(&input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "Client.InvalidParameterNotFound.SecurityGroup" {
			return nil, err
		}
		return nil, fmt.Errorf("Couldn't find SecurityGroup resource: %s", err)
	}

	return out.SecurityGroupInfo[0], nil
}
/*
func findResourceSecurityGroup(conn *computing.Computing, id string) (*computing.SecurityGroupInfoSetItem, error) {
	req := computing.DescribeSecurityGroupsInput{
		GroupName: []*string{nifcloud.String(id)},
	}

	resp, err := conn.DescribeSecurityGroups(&req)
	if err, ok := err.(awserr.Error); ok && err.Code() == "Client.InvalidParameterNotFound.SecurityGroup" {
		return nil, securityGroupNotFound{id, nil}
	}
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, securityGroupNotFound{id, nil}
	}
	if len(resp.SecurityGroupInfo) != 1 || resp.SecurityGroupInfo[0] == nil {
		return nil, securityGroupNotFound{id, resp.SecurityGroupInfo}
	}

	return resp.SecurityGroupInfo[0], nil
}

type securityGroupNotFound struct {
	id             string
	securityGroups []*computing.SecurityGroupInfoSetItem
}

func (err securityGroupNotFound) Error() string {
	if err.securityGroups == nil {
		return fmt.Sprintf("No security group with ID %q", err.id)
	}
	return fmt.Sprintf("Expected to find one security group with ID %q, got: %#v",
		err.id, err.securityGroups)
}
*/

func findRuleMatch(p *computing.RequestIpPermissionsStruct, rules []*computing.IpPermissionsSetItem) *computing.IpPermissionsSetItem {
	var rule *computing.IpPermissionsSetItem
	for _, r := range rules {
		if p.ToPort != nil && r.ToPort != nil && *p.ToPort != *r.ToPort {
			continue
		}

		if p.FromPort != nil && r.FromPort != nil && *p.FromPort != *r.FromPort {
			continue
		}

		if p.IpProtocol != nil && r.IpProtocol != nil && *p.IpProtocol != *r.IpProtocol {
			continue
		}

		if p.Description != nil && r.Description != nil && *p.Description != *r.Description {
			continue
		}

		remaining := len(p.RequestIpRanges)
		for _, ip := range p.RequestIpRanges {
			for _, rip := range r.IpRanges {
				if ip.CidrIp == nil || rip.CidrIp == nil {
					continue
				}
				if *ip.CidrIp == *rip.CidrIp {
					remaining--
				}
			}
		}

		if remaining > 0 {
			continue
		}

		remaining = len(p.RequestGroups)
		for _, rg := range p.RequestGroups {
			for _, rrg := range r.Groups {
				if rg.GroupName == nil || rrg.GroupName == nil {
					continue
				}
				if *rg.GroupName == *rrg.GroupName {
					remaining--
				}
			}
		}

		if remaining > 0 {
			continue
		}

		rule = r
	}
	return rule
}

func ipPermissionIDHash(sgID string, ip *computing.RequestIpPermissionsStruct) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", sgID))
	if ip.FromPort != nil && *ip.FromPort > 0 {
		buf.WriteString(fmt.Sprintf("%d-", *ip.FromPort))
	}
	if ip.ToPort != nil && *ip.ToPort > 0 {
		buf.WriteString(fmt.Sprintf("%d-", *ip.ToPort))
	}
	buf.WriteString(fmt.Sprintf("%s-", *ip.IpProtocol))
	buf.WriteString(fmt.Sprintf("%s-", *ip.InOut))

	if ip.RequestIpRanges != nil && *ip.RequestIpRanges[0].CidrIp != "" {
		buf.WriteString(fmt.Sprintf("%s-", *ip.RequestIpRanges[0].CidrIp))
	}
	if ip.RequestGroups != nil && *ip.RequestGroups[0].GroupName != "" {
		buf.WriteString(fmt.Sprintf("%s-", *ip.RequestGroups[0].GroupName))
	}
	if ip.Description != nil && *ip.Description != "" {
		buf.WriteString(fmt.Sprintf("%s-", *ip.Description))
	}

	return fmt.Sprintf("sgrule-%d", hashcode.String(buf.String()))
}

func expandIPPerm(d interface{}) []*computing.RequestIpPermissionsStruct {
	var perms []*computing.RequestIpPermissionsStruct
	for _, ip := range d.(*schema.Set).List() {
		perm := &computing.RequestIpPermissionsStruct{}
		if v, ok := ip.(map[string]interface{}); ok {
			protocol := strings.ToUpper(v["protocol"].(string))
			if protocol == "ICMPV6-ALL" {
				protocol = "ICMPv6-all"
			}
			switch protocol {
			case "HTTP":
				perm.SetIpProtocol("TCP")
				perm.SetFromPort(int64(80))
				perm.SetToPort(int64(80))

			case "HTTPS":
				perm.SetIpProtocol("TCP")
				perm.SetFromPort(int64(443))
				perm.SetToPort(int64(443))

			case "SSH":
				perm.SetIpProtocol("TCP")
				perm.SetFromPort(int64(22))
				perm.SetToPort(int64(22))

			case "RDP":
				perm.SetIpProtocol("TCP")
				perm.SetFromPort(int64(3389))
				perm.SetToPort(int64(3389))

			case "L2TP":
				perm.SetIpProtocol("UDP")
				perm.SetFromPort(int64(1701))
				perm.SetToPort(int64(1701))

			case "TCP":
				perm.SetIpProtocol(protocol)
				perm.SetFromPort(int64(v["from_port"].(int)))
				perm.SetToPort(int64(v["to_port"].(int)))

			case "UDP":
				perm.SetIpProtocol(protocol)
				perm.SetFromPort(int64(v["from_port"].(int)))
				perm.SetToPort(int64(v["to_port"].(int)))

			default:
				perm.SetIpProtocol(protocol)

			}

			if v["inout"].(string) != "" {
				perm.SetInOut(v["inout"].(string))
			}

			if v["description"].(string) != "" {
				perm.SetDescription(v["description"].(string))
			}

			if v["cidr_blocks"].(string) != "" {
				tmp := []*computing.RequestIpRangesStruct {
					{
						CidrIp: nifcloud.String(v["cidr_blocks"].(string)),
					},
				}
				perm.SetRequestIpRanges(tmp)
			}
			
			if v["security_groups"].(string) != "" {
				tmp := []*computing.RequestGroupsStruct {
					{
						GroupName: nifcloud.String(v["security_groups"].(string)),
					},
				}
				perm.SetRequestGroups(tmp)
			}
		}
		perms = append(perms, perm)
	}

	return perms
}
