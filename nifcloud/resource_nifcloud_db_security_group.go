package nifcloud

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/service/rdb"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNifcloudDbSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudDbSecurityGroupCreate,
		Read:   resourceNifcloudDbSecurityGroupRead,
		Update: resourceNifcloudDbSecurityGroupUpdate,
		Delete: resourceNifcloudDbSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Managed by Terraform",
			},

			"ingress": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"security_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceNifcloudDbSecurityGroupIngressHash,
			},
		},
	}
}

func resourceNifcloudDbSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).rdbconn

	var err error
	var errs []error

	opts := rdb.CreateDBSecurityGroupInput{
		DBSecurityGroupName:        nifcloud.String(d.Get("name").(string)),
		DBSecurityGroupDescription: nifcloud.String(d.Get("description").(string)),
		NiftyAvailabilityZone:      nifcloud.String(d.Get("availability_zone").(string)),
	}

	log.Printf("[DEBUG] DB Security Group create configuration: %#v", opts)
	_, err = conn.CreateDBSecurityGroup(&opts)
	if err != nil {
		return fmt.Errorf("Error creating DB Security Group: %s", err)
	}

	d.SetId(d.Get("name").(string))

	log.Printf("[INFO] DB Security Group ID: %s", d.Id())

	sg, err := resourceNifcloudDbSecurityGroupRetrieve(d, meta)
	if err != nil {
		return err
	}

	ingresses := d.Get("ingress").(*schema.Set)
	for _, ing := range ingresses.List() {
		err := resourceNifcloudDbSecurityGroupAuthorizeRule(ing, *sg.DBSecurityGroupName, conn)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return &multierror.Error{Errors: errs}
	}

	log.Println(
		"[INFO] Waiting for Ingress Authorizations to be authorized")

	stateConf := &resource.StateChangeConf{
		Pending: []string{"authorizing"},
		Target:  []string{"authorized"},
		Refresh: resourceNifcloudDbSecurityGroupStateRefreshFunc(d, meta),
		Timeout: 10 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceNifcloudDbSecurityGroupRead(d, meta)
}

func resourceNifcloudDbSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	sg, err := resourceNifcloudDbSecurityGroupRetrieve(d, meta)
	if err != nil {
		return err
	}

	d.Set("name", sg.DBSecurityGroupName)
	d.Set("description", sg.DBSecurityGroupDescription)

	// Create an empty schema.Set to hold all ingress rules
	rules := &schema.Set{
		F: resourceNifcloudDbSecurityGroupIngressHash,
	}

	for _, v := range sg.IPRanges {
		rule := map[string]interface{}{"cidr": *v.CIDRIP}
		rules.Add(rule)
	}

	for _, g := range sg.EC2SecurityGroups {
		rule := map[string]interface{}{}
		if g.EC2SecurityGroupName != nil {
			rule["security_group_name"] = *g.EC2SecurityGroupName
		}
		rules.Add(rule)
	}

	d.Set("ingress", rules)

	return nil
}

func resourceNifcloudDbSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).rdbconn

	d.Partial(true)

	if d.HasChange("ingress") {
		sg, err := resourceNifcloudDbSecurityGroupRetrieve(d, meta)
		if err != nil {
			return err
		}

		oi, ni := d.GetChange("ingress")
		if oi == nil {
			oi = new(schema.Set)
		}
		if ni == nil {
			ni = new(schema.Set)
		}

		ois := oi.(*schema.Set)
		nis := ni.(*schema.Set)
		removeIngress := ois.Difference(nis).List()
		newIngress := nis.Difference(ois).List()

		// DELETE old Ingress rules
		for _, ing := range removeIngress {
			err := resourceNifcloudDbSecurityGroupRevokeRule(ing, *sg.DBSecurityGroupName, conn)
			if err != nil {
				return err
			}
		}

		// ADD new/updated Ingress rules
		for _, ing := range newIngress {
			err := resourceNifcloudDbSecurityGroupAuthorizeRule(ing, *sg.DBSecurityGroupName, conn)
			if err != nil {
				return err
			}
		}
	}
	d.Partial(false)

	return resourceNifcloudDbSecurityGroupRead(d, meta)
}

func resourceNifcloudDbSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).rdbconn

	log.Printf("[DEBUG] DB Security Group destroy: %v", d.Id())

	opts := rdb.DeleteDBSecurityGroupInput{DBSecurityGroupName: nifcloud.String(d.Id())}

	log.Printf("[DEBUG] DB Security Group destroy configuration: %v", opts)
	_, err := conn.DeleteDBSecurityGroup(&opts)

	if err != nil {
		if isNifcloudErr(err, "Client.InvalidParameterNotFound.DBSecurityGroup", "") {
			return nil
		}
		return err
	}

	return nil
}

func resourceNifcloudDbSecurityGroupRetrieve(d *schema.ResourceData, meta interface{}) (*rdb.DBSecurityGroup, error) {
	conn := meta.(*NifcloudClient).rdbconn

	opts := rdb.DescribeDBSecurityGroupsInput{
		DBSecurityGroupName: nifcloud.String(d.Id()),
	}

	log.Printf("[DEBUG] DB Security Group describe configuration: %#v", opts)

	resp, err := conn.DescribeDBSecurityGroups(&opts)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving DB Security Groups: %s", err)
	}

	if len(resp.DBSecurityGroups) != 1 ||
		*resp.DBSecurityGroups[0].DBSecurityGroupName != d.Id() {
		return nil, fmt.Errorf("Unable to find DB Security Group: %#v", resp.DBSecurityGroups)
	}

	return resp.DBSecurityGroups[0], nil
}

// Authorizes the ingress rule on the db security group
func resourceNifcloudDbSecurityGroupAuthorizeRule(ingress interface{}, dbSecurityGroupName string, conn *rdb.Rdb) error {
	ing := ingress.(map[string]interface{})

	opts := rdb.AuthorizeDBSecurityGroupIngressInput{
		DBSecurityGroupName: nifcloud.String(dbSecurityGroupName),
	}

	if attr, ok := ing["cidr"]; ok && attr != "" {
		opts.CIDRIP = nifcloud.String(attr.(string))
	}

	if attr, ok := ing["security_group_name"]; ok && attr != "" {
		opts.EC2SecurityGroupName = nifcloud.String(attr.(string))
	}

	log.Printf("[DEBUG] Authorize ingress rule configuration: %#v", opts)

	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := conn.AuthorizeDBSecurityGroupIngress(&opts)
		
		if err != nil {
			if isNifcloudErr(err, "Client.ResourceIncorrectState.DBSecurityGroup.Unavailable", "") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error authorizing security group ingress: %s", err)
	}

	return nil
}

// Revokes the ingress rule on the db security group
func resourceNifcloudDbSecurityGroupRevokeRule(ingress interface{}, dbSecurityGroupName string, conn *rdb.Rdb) error {
	ing := ingress.(map[string]interface{})

	opts := rdb.RevokeDBSecurityGroupIngressInput{
		DBSecurityGroupName: nifcloud.String(dbSecurityGroupName),
	}

	if attr, ok := ing["cidr"]; ok && attr != "" {
		opts.CIDRIP = nifcloud.String(attr.(string))
	}

	if attr, ok := ing["security_group_name"]; ok && attr != "" {
		opts.EC2SecurityGroupName = nifcloud.String(attr.(string))
	}

	log.Printf("[DEBUG] Revoking ingress rule configuration: %#v", opts)

	_, err := conn.RevokeDBSecurityGroupIngress(&opts)

	if err != nil {
		return fmt.Errorf("Error revoking security group ingress: %s", err)
	}

	return nil
}

func resourceNifcloudDbSecurityGroupIngressHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["cidr"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["security_group_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func resourceNifcloudDbSecurityGroupStateRefreshFunc(
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := resourceNifcloudDbSecurityGroupRetrieve(d, meta)

		if err != nil {
			log.Printf("Error on retrieving DB Security Group when waiting: %s", err)
			return nil, "", err
		}

		statuses := make([]string, 0, len(v.EC2SecurityGroups)+len(v.IPRanges))
		for _, ec2g := range v.EC2SecurityGroups {
			statuses = append(statuses, *ec2g.Status)
		}
		for _, ips := range v.IPRanges {
			statuses = append(statuses, *ips.Status)
		}

		for _, stat := range statuses {
			// Not done
			if stat != "authorized" {
				return nil, "authorizing", nil
			}
		}

		return v, "authorized", nil
	}
}