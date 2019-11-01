package nifcloud

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/service/rdb"
)

func resourceNifcloudDbParameterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudDbParameterGroupCreate,
		Read:   resourceNifcloudDbParameterGroupRead,
		Update: resourceNifcloudDbParameterGroupUpdate,
		Delete: resourceNifcloudDbParameterGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validateDbParamGroupName,
			},
			"family": {
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
			"parameter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"apply_method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "immediate",
						},
					},
				},
				Set: resourceNifcloudDbParameterHash,
			},
		},
	}
}

func resourceNifcloudDbParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	rdbconn := meta.(*NifcloudClient).rdbconn

	var groupName string
	if v, ok := d.GetOk("name"); ok {
		groupName = v.(string)
	}
	d.Set("name", groupName)

	createOpts := rdb.CreateDBParameterGroupInput{
		DBParameterGroupName:   nifcloud.String(groupName),
		DBParameterGroupFamily: nifcloud.String(d.Get("family").(string)),
		Description:            nifcloud.String(d.Get("description").(string)),
	}

	log.Printf("[DEBUG] Create DB Parameter Group: %#v", createOpts)
	resp, err := rdbconn.CreateDBParameterGroup(&createOpts)
	if err != nil {
		return fmt.Errorf("Error creating DB Parameter Group: %s", err)
	}

	d.Partial(true)
	d.SetPartial("name")
	d.SetPartial("family")
	d.SetPartial("description")
	d.Partial(false)

	d.SetId(nifcloud.StringValue(resp.DBParameterGroup.DBParameterGroupName))
	log.Printf("[INFO] DB Parameter Group ID: %s", d.Id())

	return resourceNifcloudDbParameterGroupUpdate(d, meta)
}

func resourceNifcloudDbParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	rdbconn := meta.(*NifcloudClient).rdbconn

	describeOpts := rdb.DescribeDBParameterGroupsInput{
		DBParameterGroupName: nifcloud.String(d.Id()),
	}

	describeResp, err := rdbconn.DescribeDBParameterGroups(&describeOpts)
	if err != nil {
		if isNifcloudErr(err, "Client.InvalidParameterNotFound.DBParameterGroup", "") {
			log.Printf("[WARN] DB Parameter Group (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if len(describeResp.DBParameterGroups) != 1 ||
		*describeResp.DBParameterGroups[0].DBParameterGroupName != d.Id() {
		return fmt.Errorf("Unable to find Parameter Group: %#v", describeResp.DBParameterGroups)
	}

	d.Set("name", describeResp.DBParameterGroups[0].DBParameterGroupName)
	d.Set("family", describeResp.DBParameterGroups[0].DBParameterGroupFamily)
	d.Set("description", describeResp.DBParameterGroups[0].Description)

	configParams := d.Get("parameter").(*schema.Set)
	describeParametersOpts := rdb.DescribeDBParametersInput{
		DBParameterGroupName: nifcloud.String(d.Id()),
	}
	if configParams.Len() < 1 {
		// if we don't have any params in the ResourceData already, two possibilities
		// first, we don't have a config available to us. Second, we do, but it has
		// no parameters. We're going to assume the first, to be safe. In this case,
		// we're only going to ask for the user-modified values, because any defaults
		// the user may have _also_ set are indistinguishable from the hundreds of
		// defaults AWS sets. If the user hasn't set any parameters, this will return
		// an empty list anyways, so we just make some unnecessary requests. But in
		// the more common case (I assume) of an import, this will make fewer requests
		// and "do the right thing".
		describeParametersOpts.Source = nifcloud.String("user")
	}

	var parameters []*rdb.Parameter
	err = rdbconn.DescribeDBParametersPages(&describeParametersOpts,
		func(describeParametersResp *rdb.DescribeDBParametersOutput, lastPage bool) bool {
			parameters = append(parameters, describeParametersResp.Parameters...)
			return !lastPage
		})
	if err != nil {
		return err
	}
//	log.Printf("[DEBUG] ******************* DescribeDBParametersPages: %v", parameters)
	
	var userParams []*rdb.Parameter
	if configParams.Len() < 1 {
		// if we have no config/no parameters in config, we've already asked for only
		// user-modified values, so we can just use the entire response.
		userParams = parameters
	} else {
		// if we have a config available to us, we have two possible classes of value
		// in the config. On the one hand, the user could have specified a parameter
		// that _actually_ changed things, in which case its Source would be set to
		// user. On the other, they may have specified a parameter that coincides with
		// the default value. In that case, the Source will be set to "system" or
		// "engine-default". We need to set the union of all "user" Source parameters
		// _and_ the "system"/"engine-default" Source parameters _that appear in the
		// config_ in the state, or the user gets a perpetual diff. See
		// terraform-providers/terraform-provider-aws#593 for more context and details.
		confParams, err := expandParameters(configParams.List())
		if err != nil {
			return err
		}
		for _, param := range parameters {
			if param.Source == nil || param.ParameterName == nil {
				continue
			}
			if *param.Source == "user" {
				userParams = append(userParams, param)
				continue
			}
			var paramFound bool
			for _, cp := range confParams {
				if cp.ParameterName == nil {
					continue
				}
				if *cp.ParameterName == *param.ParameterName {
					userParams = append(userParams, param)
					break
				}
			}
			if !paramFound {
				log.Printf("[DEBUG] Not persisting %s to state, as its source is %q and it isn't in the config", *param.ParameterName, *param.Source)
			}
		}
	}

	err = d.Set("parameter", flattenParameters(userParams))
	if err != nil {
		return fmt.Errorf("error setting 'parameter' in state: %#v", err)
	}

	return nil
}

func resourceNifcloudDbParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	rdbconn := meta.(*NifcloudClient).rdbconn

	d.Partial(true)

	if d.HasChange("parameter") {
		o, n := d.GetChange("parameter")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		// Expand the "parameter" set to aws-sdk-go compat []rdb.Parameter
		parameters, err := expandParameters(ns.Difference(os).List())
		if err != nil {
			return err
		}

		if len(parameters) > 0 {
			// We can only modify 20 parameters at a time, so walk them until
			// we've got them all.
			maxParams := 20
			for parameters != nil {
				var paramsToModify []*rdb.RequestParametersStruct
				if len(parameters) <= maxParams {
					paramsToModify, parameters = parameters[:], nil
				} else {
					paramsToModify, parameters = parameters[:maxParams], parameters[maxParams:]
				}
				modifyOpts := rdb.ModifyDBParameterGroupInput{
					DBParameterGroupName: nifcloud.String(d.Get("name").(string)),
					Parameters:           paramsToModify,
				}

				log.Printf("[DEBUG] Modify DB Parameter Group: %s", modifyOpts)
				_, err = rdbconn.ModifyDBParameterGroup(&modifyOpts)
				if err != nil {
					return fmt.Errorf("Error modifying DB Parameter Group: %s", err)
				}
			}
			d.SetPartial("parameter")
		}
	}

	d.Partial(false)

	return resourceNifcloudDbParameterGroupRead(d, meta)
}

func resourceNifcloudDbParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).rdbconn
	deleteOpts := rdb.DeleteDBParameterGroupInput{
		DBParameterGroupName: nifcloud.String(d.Id()),
	}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteDBParameterGroup(&deleteOpts)
		if err != nil {
//			if isNifcloudErr(err, "DBParameterGroupNotFoundFault", "") || isNifcloudErr(err, "InvalidDBParameterGroupState", "") {
//				return resource.RetryableError(err)
//			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = conn.DeleteDBParameterGroup(&deleteOpts)
	}
	if err != nil {
		if isNifcloudErr(err, "Client.InvalidParameterNotFound.DBParameterGroup", "") {
			return nil
		}
		return fmt.Errorf("Error deleting DB parameter group: %s", err)
	}
	return nil
}

func resourceNifcloudDbParameterHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	// Store the value as a lower case string, to match how we store them in flattenParameters
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["value"].(string))))

	return hashcode.String(buf.String())
}

func validateDbParamGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if (len(value) < 1) || (len(value) > 255) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 255 characters in length", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z]`).MatchString(value) {
		errors = append(errors, fmt.Errorf("first character of %q must be a letter", k))
	}

	if regexp.MustCompile(`--`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q cannot contain two consecutive hyphens", k))
	}

	if strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("%q cannot end with a - character", k))
	}

	if !regexp.MustCompile(`^[0-9a-zA-Z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q can only contain alphanumeric and %q characters", k, "-"))
	}

	return
}

// Takes the result of flatmap.Expand for an array of parameters and
// returns Parameter API compatible objects
func expandParameters(configured []interface{}) ([]*rdb.RequestParametersStruct, error) {
	var parameters []*rdb.RequestParametersStruct

	// Loop over our configured parameters and create
	// an array of aws-sdk-go compatible objects
	for _, pRaw := range configured {
		data := pRaw.(map[string]interface{})

		if data["name"].(string) == "" {
			continue
		}

		p := &rdb.RequestParametersStruct{
			ApplyMethod:    nifcloud.String(data["apply_method"].(string)),
			ParameterName:  nifcloud.String(data["name"].(string)),
			ParameterValue: nifcloud.String(data["value"].(string)),
		}

		parameters = append(parameters, p)
	}

	return parameters, nil
}

// Flattens an array of Parameters into a []map[string]interface{}
func flattenParameters(list []*rdb.Parameter) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		if i.ParameterName != nil {
			r := make(map[string]interface{})
			r["name"] = strings.ToLower(*i.ParameterName)
			// Default empty string, guard against nil parameter values
			r["value"] = ""
			if i.ParameterValue != nil {
				r["value"] = strings.ToLower(*i.ParameterValue)
			}
			if i.ApplyMethod != nil {
				r["apply_method"] = strings.ToLower(*i.ApplyMethod)
			}

			result = append(result, r)
		}
	}
	return result
}