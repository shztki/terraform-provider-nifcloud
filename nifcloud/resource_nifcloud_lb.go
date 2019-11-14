package nifcloud

import (
	"fmt"
	"log"
	"time"
	"strconv"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNifcloudLb() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudLbCreate,
		Read:   resourceNifcloudLbRead,
		Update: resourceNifcloudLbUpdate,
		Delete: resourceNifcloudLbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateLbName,
			},

//			"availability_zone": {
//				Type:     schema.TypeString,
//				Optional: true,
//				Computed: true,
//				ForceNew: true,
//			},

			"network_volume": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntInSlice([]int{10,20,30,40,100,200,300,400,500,600,700,800,900,1000,1100,1200,1300,1400,1500,1600,1700,1800,1900,2000}),
			},

			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "v4",
				ValidateFunc: validation.StringInSlice([]string{"v4","v6",""}, true),
			},

			"accounting_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "2",
				ValidateFunc: validation.StringInSlice([]string{"1","2",""}, true),
			},

			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard","ats",""}, true),
			},

			"instances": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
//				Computed: true,
				Set:      schema.HashString,
			},

			"filter_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
//				Default:      "1",
				ValidateFunc: validation.StringInSlice([]string{"1","2"}, true),
			},

			"filter_ip_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 10,
				Optional: true,
//				Computed: true,
				Set:      schema.HashString,
			},

			"session_stickiness_policy_enable": {
				Type:         schema.TypeBool,
				Optional:     true,
//				Computed:     true,
			},

			"session_stickiness_policy_expiration_period": {
				Type:     schema.TypeInt,
				Optional: true,
//				Computed: true,
				ValidateFunc: validation.IntBetween(3, 60),
			},

			"sorry_page_enable": {
				Type:         schema.TypeBool,
				Optional:     true,
//				Computed:     true,
			},

			"sorry_page_status_code": {
				Type:     schema.TypeInt,
				Optional: true,
//				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{200,503}),
			},

			"ssl_certificate_id": {
				Type:         schema.TypeString,
				Optional:     true,
			},

			"ssl_policy_id": {
				Type:         schema.TypeString,
				Optional:     true,
			},

			"listener": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateListenerProtocol(),
						},

						"lb_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"instance_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"balancing_type": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntInSlice([]int{1,2}),
						},
					},
				},	
				//Set: resourceNifcloudLbListenerHash,
			},

			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
//						"healthy_threshold": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(2, 10),
//						},

						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},

						"target": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateHeathCheckTarget,
						},

						"interval": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(5, 300),
						},

//						"timeout": {
//							Type:         schema.TypeInt,
//							Required:     true,
//							ValidateFunc: validation.IntBetween(2, 60),
//						},
					},
				},
			},
		},
	}
}

func resourceNifcloudLbCreate(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*NifcloudClient).computingconn

	// Expand the "RequestListenerStruct" set to nifcloud-sdk-go compat []*computing.RequestListenerStruct
	listeners, err := expandListeners(d.Get("listener").(*schema.Set).List())
	if err != nil {
		return err
	}

	var elbName string
	if v, ok := d.GetOk("name"); ok {
		elbName = v.(string)
	}

	// Provision the elb
	elbOpts := &computing.CreateLoadBalancerInput{
		LoadBalancerName: nifcloud.String(elbName),
		Listeners:        listeners,
		AccountingType:   nifcloud.String(d.Get("accounting_type").(string)),
		IpVersion:        nifcloud.String(d.Get("ip_version").(string)),
		NetworkVolume:    nifcloud.Int64(int64(d.Get("network_volume").(int))),
		PolicyType:       nifcloud.String(d.Get("policy_type").(string)),
	}

//	if v, ok := d.GetOk("availability_zone"); ok {
//		elbOpts.AvailabilityZones = expandStringList(v.(*schema.Set).List())
//	}

	log.Printf("[DEBUG] LB create configuration: %v", elbOpts)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := elbconn.CreateLoadBalancer(elbOpts)

		if err != nil {
//			if awsErr, ok := err.(awserr.Error); ok {
//				// Check for IAM SSL Cert error, eventual consistancy issue
//				if awsErr.Code() == "CertificateNotFound" {
//					return resource.RetryableError(
//						fmt.Errorf("Error creating LB Listener with SSL Cert, retrying: %s", err))
//				}
//			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = elbconn.CreateLoadBalancer(elbOpts)
	}
	if err != nil {
		return fmt.Errorf("Error creating LB: %s", err)
	}

	// Assign the elb's unique identifier for use later
	d.SetId(elbName)
	log.Printf("[INFO] LB ID: %s", d.Id())

	// Enable partial mode and record what we set
	d.Partial(true)
	d.SetPartial("name")
	d.SetPartial("accounting_type")
//	d.SetPartial("availability_zones")
	d.SetPartial("listener")
	d.SetPartial("ip_version")
	d.SetPartial("network_volume")
	d.SetPartial("policy_type")

	return resourceNifcloudLbUpdate(d, meta)
}

func resourceNifcloudLbRead(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*NifcloudClient).computingconn
	elbName := d.Id()

	// Expand the "RequestLoadBalancerNamesStruct" set to nifcloud-sdk-go compat []*computing.RequestLoadBalancerNamesStruct
	loadBalancerNames, err := expandRequestLoadBalancerNames(elbName, d.Get("listener").(*schema.Set).List())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] LB describe loadBalancerNames: %v", loadBalancerNames)
	
	// Retrieve the LB properties for updating the state
	describeElbOpts := &computing.DescribeLoadBalancersInput{
		LoadBalancerNames: loadBalancerNames,
	}

	describeResp, err := elbconn.DescribeLoadBalancers(describeElbOpts)
	if err != nil {
		if isLoadBalancerNotFound(err) {
			// The LB is gone now, so just remove it from the state
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving LB: %s", err)
	}
	if len(describeResp.DescribeLoadBalancersResult.LoadBalancerDescriptions) != 1 {
		return fmt.Errorf("Unable to find LB: %v", describeResp.DescribeLoadBalancersResult.LoadBalancerDescriptions)
	}
	log.Printf("[DEBUG] LB describe : %v", describeResp.DescribeLoadBalancersResult.LoadBalancerDescriptions)
//	var describeResp *computing.DescribeLoadBalancersOutput
//	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
//		var err error
//		describeResp, err = elbconn.DescribeLoadBalancers(describeElbOpts)
//		if err != nil {
//			if isLoadBalancerNotFound(err) {
//				// The LB is gone now, so just remove it from the state
//				d.SetId("")
//				return nil
//			}
//			return resource.NonRetryableError(err)
//		}
//		if len(describeResp.LoadBalancerDescriptions) != 1 {
//			return resource.RetryableError(err)
//		}
//		return nil
//	})
//	if isResourceTimeoutError(err) {
//		describeResp, err = elbconn.DescribeLoadBalancers(describeElbOpts)
//	}
//	if err != nil {
//		return fmt.Errorf("Error describe LB: %s", err)
//	}
//	if len(describeResp.LoadBalancerDescriptions) != 1 {
//		return fmt.Errorf("Unable to find LB: %v", describeResp.LoadBalancerDescriptions)
//	}

	return flatflattenNifcloudLbResource(d, describeResp.DescribeLoadBalancersResult.LoadBalancerDescriptions[0])
}

// flatflattenNifcloudLbResource takes a *elbv2.LoadBalancer and populates all respective resource fields.
func flatflattenNifcloudLbResource(d *schema.ResourceData, lb *computing.LoadBalancerDescriptionsMemberItem) error {

	d.Set("name", lb.LoadBalancerName)
//	d.Set("availability_zone", lb.AvailabilityZones[0])
	d.Set("listener", flattenListeners(lb.ListenerDescriptions))
	d.Set("network_volume", lb.NetworkVolume)
	d.Set("accounting_type", lb.NextMonthAccountingType)
	d.Set("policy_type", lb.PolicyType)
	d.Set("dns_name", lb.DNSName)
	d.Set("instances", flattenInstances(lb.Instances))
	d.Set("filter_type", lb.Filter.FilterType)
	addresses := flattenIpAddresses(lb.Filter.IPAddresses)
	if len(addresses) != 1 && addresses[0] == "" {
		d.Set("filter_ip_addresses", addresses)
	}
	if lb.Option != nil && lb.Option.SessionStickinessPolicy != nil {
		d.Set("session_stickiness_policy_enable", *lb.Option.SessionStickinessPolicy.Enabled)
		d.Set("session_stickiness_policy_expiration_period", *lb.Option.SessionStickinessPolicy.ExpirationPeriod)
	}
	if lb.Option != nil && lb.Option.SorryPage != nil {
		d.Set("sorry_page_enable", *lb.Option.SorryPage.Enabled)
		d.Set("sorry_page_status_code", *lb.Option.SorryPage.StatusCode)
	}
	if v := flattenSSLCertificateID(lb.ListenerDescriptions); v != nil {
		d.Set("ssl_certificate_id", v[0])
	}
	if v := flattenSSLPolicyID(lb.ListenerDescriptions); v != nil {
		d.Set("ssl_policy_id", v[0])
	}

//	// There's only one health check, so save that to state as we
//	// currently can
	if *lb.HealthCheck.Target != "" {
		d.Set("health_check", flattenHealthCheck(lb.HealthCheck))
	}

	return nil
}

func resourceNifcloudLbUpdate(d *schema.ResourceData, meta interface{}) error {

	elbconn := meta.(*NifcloudClient).computingconn

	d.Partial(true)

	req := &computing.UpdateLoadBalancerInput{
		LoadBalancerName: nifcloud.String(d.Id()),
	}
	requestUpdate := false
	var oldListener *computing.RequestListenerStruct
	var err error

	if d.HasChange("accounting_type") {
		d.SetPartial("accounting_type")
		var accountingType int
		accountingType, _ = strconv.Atoi(d.Get("accounting_type").(string))
		req.AccountingTypeUpdate = nifcloud.Int64(int64(accountingType))
		requestUpdate = true
	}
	if d.HasChange("network_volume") {
		d.SetPartial("network_volume")
		req.NetworkVolumeUpdate = nifcloud.Int64(int64(d.Get("network_volume").(int)))
		requestUpdate = true
	}
	if d.HasChange("name") {
		d.SetPartial("network_volume")
		req.LoadBalancerNameUpdate = nifcloud.String(d.Get("name").(string))
		requestUpdate = true
	}
	if d.HasChange("listener") {
		o, n := d.GetChange("listener")
		oldListener, err = expandListener(o.(*schema.Set).List())
		if err != nil {
			return err
		}
		newListener, err := expandListener(n.(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener == nil {
			oldListener = newListener
		}

		updateListner := &computing.RequestListenerUpdateStruct{
			InstancePort:          oldListener.InstancePort,
			LoadBalancerPort:      oldListener.LoadBalancerPort,
			RequestListenerStruct: newListener,
		}
		req.ListenerUpdate = updateListner
		requestUpdate = true
		d.SetPartial("listener")
	}

	if d.HasChange("health_check") {
		log.Printf("[INFO] Updating HealthCheck %s ", d.Id())
		hc := d.Get("health_check").([]interface{})
		if len(hc) > 0 {
			check := hc[0].(map[string]interface{})
			listener, err := expandListener(d.Get("listener").(*schema.Set).List())
			if err != nil {
				return err
			}
			if oldListener != nil {
				listener = oldListener
			}

			configureHealthCheckOpts := computing.ConfigureHealthCheckInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
				HealthCheck: &computing.RequestHealthCheckStruct{
//					HealthyThreshold:   nifcloud.Int64(int64(check["healthy_threshold"].(int))),
					UnhealthyThreshold: nifcloud.Int64(int64(check["unhealthy_threshold"].(int))),
					Interval:           nifcloud.Int64(int64(check["interval"].(int))),
					Target:             nifcloud.String(check["target"].(string)),
//					Timeout:            nifcloud.Int64(int64(check["timeout"].(int))),
				},
			}
			_, err = elbconn.ConfigureHealthCheck(&configureHealthCheckOpts)
			if err != nil {
				return fmt.Errorf("Failure configuring health check for LB: %s", err)
			}
			d.SetPartial("health_check")
		}
	}

	// If we currently have instances, or did have instances,
	// we want to figure out what to add and remove from the load
	// balancer
	if d.HasChange("instances") {
		log.Printf("[INFO] Updating Instances %s ", d.Id())
		o, n := d.GetChange("instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := expandInstanceString(os.Difference(ns).List())
		add := expandInstanceString(ns.Difference(os).List())

		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}
		
		if len(add) > 0 {
			registerInstancesOpts := computing.RegisterInstancesWithLoadBalancerInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
				Instances:        add,
			}

			_, err := elbconn.RegisterInstancesWithLoadBalancer(&registerInstancesOpts)
			if err != nil {
				return fmt.Errorf("Failure registering instances with LB: %s", err)
			}
		}
		if len(remove) > 0 {
			deRegisterInstancesOpts := computing.DeregisterInstancesFromLoadBalancerInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
				Instances:        remove,
			}

			_, err := elbconn.DeregisterInstancesFromLoadBalancer(&deRegisterInstancesOpts)
			if err != nil {
				return fmt.Errorf("Failure deregistering instances from LB: %s", err)
			}
		}

		d.SetPartial("instances")
	}

	if d.HasChange("filter_ip_addresses") {
		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}

		if d.HasChange("filter_type") {
			log.Printf("[INFO] First Updating FilterType")
			setFilterForLoadBalancerOpts := computing.SetFilterForLoadBalancerInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
				FilterType:       nifcloud.String(d.Get("filter_type").(string)),
			}

			_, err := elbconn.SetFilterForLoadBalancer(&setFilterForLoadBalancerOpts)
			if err != nil {
				return fmt.Errorf("Failure setting filter for LB: %s", err)
			}

			d.SetPartial("filter_type")
		}
		o, n := d.GetChange("filter_ip_addresses")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		// DELETE old filter
		if len(remove) > 0 {
			log.Printf("[INFO] Updating Filter Delete")
			addresses := expandDeleteFilter(remove)
			setFilterForLoadBalancerOpts := computing.SetFilterForLoadBalancerInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
				FilterType:       nifcloud.String(d.Get("filter_type").(string)),
				IPAddresses:      addresses,
			}

			_, err := elbconn.SetFilterForLoadBalancer(&setFilterForLoadBalancerOpts)
			if err != nil {
				return fmt.Errorf("Failure setting filter for LB: %s", err)
			}
		}

		// ADD new/updated filter
		if len(add) > 0 {
			log.Printf("[INFO] Updating Filter Add")
			addresses := expandAddFilter(add)
			setFilterForLoadBalancerOpts := computing.SetFilterForLoadBalancerInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
				FilterType:       nifcloud.String(d.Get("filter_type").(string)),
				IPAddresses:      addresses,
			}

			_, err := elbconn.SetFilterForLoadBalancer(&setFilterForLoadBalancerOpts)
			if err != nil {
				return fmt.Errorf("Failure setting filter for LB: %s", err)
			}
		}

		d.SetPartial("filter_ip_addresses")

	} else if d.HasChange("filter_type") {
		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}

		log.Printf("[INFO] Updating FilterType only")
		setFilterForLoadBalancerOpts := computing.SetFilterForLoadBalancerInput{
			LoadBalancerName: nifcloud.String(d.Id()),
			InstancePort:     listener.InstancePort,
			LoadBalancerPort: listener.LoadBalancerPort,
			FilterType:       nifcloud.String(d.Get("filter_type").(string)),
		}

		_, err = elbconn.SetFilterForLoadBalancer(&setFilterForLoadBalancerOpts)
		if err != nil {
			return fmt.Errorf("Failure setting filter for LB: %s", err)
		}

		d.SetPartial("filter_type")
	}

	if d.HasChange("session_stickiness_policy_enable") {
		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}

		log.Printf("[INFO] Updating LoadBalancerOption")
		enableOpts := &computing.RequestSessionStickinessPolicyUpdateStruct{Enable: nifcloud.Bool(false)}
		if v, ok := d.GetOk("session_stickiness_policy_enable"); ok {
			enableOpts.Enable = nifcloud.Bool(v.(bool))
		}
		if v, ok := d.GetOk("session_stickiness_policy_expiration_period"); ok {
			enableOpts.ExpirationPeriod = nifcloud.Int64(int64(v.(int)))
		}
		updateLoadBalancerOptionOpts := computing.UpdateLoadBalancerOptionInput{
			LoadBalancerName:              nifcloud.String(d.Id()),
			InstancePort:                  listener.InstancePort,
			LoadBalancerPort:              listener.LoadBalancerPort,
			SessionStickinessPolicyUpdate: enableOpts,
		}

		_, err = elbconn.UpdateLoadBalancerOption(&updateLoadBalancerOptionOpts)
		if err != nil {
			return fmt.Errorf("Failure updating LoadBalancerOption for LB: %s", err)
		}

		d.SetPartial("session_stickiness_policy_enable")
	} else if d.HasChange("session_stickiness_policy_expiration_period") {
		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}

		log.Printf("[INFO] Updating LoadBalancerOption")
		enableOpts := &computing.RequestSessionStickinessPolicyUpdateStruct{Enable: nifcloud.Bool(false)}
		if v, ok := d.GetOk("session_stickiness_policy_enable"); ok {
			enableOpts.Enable = nifcloud.Bool(v.(bool))
		}
		if v, ok := d.GetOk("session_stickiness_policy_expiration_period"); ok {
			enableOpts.ExpirationPeriod = nifcloud.Int64(int64(v.(int)))
		}
		updateLoadBalancerOptionOpts := computing.UpdateLoadBalancerOptionInput{
			LoadBalancerName:              nifcloud.String(d.Id()),
			InstancePort:                  listener.InstancePort,
			LoadBalancerPort:              listener.LoadBalancerPort,
			SessionStickinessPolicyUpdate: enableOpts,
		}

		_, err = elbconn.UpdateLoadBalancerOption(&updateLoadBalancerOptionOpts)
		if err != nil {
			return fmt.Errorf("Failure updating LoadBalancerOption for LB: %s", err)
		}

		d.SetPartial("session_stickiness_policy_enable")
	}

	if d.HasChange("sorry_page_enable") {
		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}

		log.Printf("[INFO] Updating LoadBalancerOption")
		enableOpts := &computing.RequestSorryPageUpdateStruct{Enable: nifcloud.Bool(false)}
		if v, ok := d.GetOk("sorry_page_enable"); ok {
			enableOpts.Enable = nifcloud.Bool(v.(bool))
		}
		if v, ok := d.GetOk("sorry_page_status_code"); ok {
			enableOpts.StatusCode = nifcloud.Int64(int64(v.(int)))
		}
		updateLoadBalancerOptionOpts := computing.UpdateLoadBalancerOptionInput{
			LoadBalancerName: nifcloud.String(d.Id()),
			InstancePort:     listener.InstancePort,
			LoadBalancerPort: listener.LoadBalancerPort,
			SorryPageUpdate:  enableOpts,
		}

		_, err = elbconn.UpdateLoadBalancerOption(&updateLoadBalancerOptionOpts)
		if err != nil {
			return fmt.Errorf("Failure updating LoadBalancerOption for LB: %s", err)
		}

		d.SetPartial("sorry_page_enable")
	} else if d.HasChange("sorry_page_status_code") {
		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}

		log.Printf("[INFO] Updating LoadBalancerOption")
		enableOpts := &computing.RequestSorryPageUpdateStruct{Enable: nifcloud.Bool(false)}
		if v, ok := d.GetOk("sorry_page_enable"); ok {
			enableOpts.Enable = nifcloud.Bool(v.(bool))
		}
		if v, ok := d.GetOk("sorry_page_status_code"); ok {
			enableOpts.StatusCode = nifcloud.Int64(int64(v.(int)))
		}
		updateLoadBalancerOptionOpts := computing.UpdateLoadBalancerOptionInput{
			LoadBalancerName: nifcloud.String(d.Id()),
			InstancePort:     listener.InstancePort,
			LoadBalancerPort: listener.LoadBalancerPort,
			SorryPageUpdate:  enableOpts,
		}

		_, err = elbconn.UpdateLoadBalancerOption(&updateLoadBalancerOptionOpts)
		if err != nil {
			return fmt.Errorf("Failure updating LoadBalancerOption for LB: %s", err)
		}

		d.SetPartial("sorry_page_status_code")
	}

	if d.HasChange("ssl_certificate_id") {
		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}

		log.Printf("[INFO] Updating SetLoadBalancerListenerSSLCertificate")
		if v, ok := d.GetOk("ssl_certificate_id"); ok {
			setLoadBalancerListenerSSLCertificateOpts := computing.SetLoadBalancerListenerSSLCertificateInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
				SSLCertificateId: nifcloud.String(v.(string)),
			}
	
			_, err = elbconn.SetLoadBalancerListenerSSLCertificate(&setLoadBalancerListenerSSLCertificateOpts)
			if err != nil {
				return fmt.Errorf("Failure SetLoadBalancerListenerSSLCertificate for LB: %s", err)
			}	
		} else {
			unsetLoadBalancerListenerSSLCertificateOpts := computing.UnsetLoadBalancerListenerSSLCertificateInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
			}
	
			_, err = elbconn.UnsetLoadBalancerListenerSSLCertificate(&unsetLoadBalancerListenerSSLCertificateOpts)
			if err != nil {
				return fmt.Errorf("Failure UnsetLoadBalancerListenerSSLCertificate for LB: %s", err)
			}
		}

		d.SetPartial("ssl_certificate_id")
	}

	if d.HasChange("ssl_policy_id") {
		listener, err := expandListener(d.Get("listener").(*schema.Set).List())
		if err != nil {
			return err
		}
		if oldListener != nil {
			listener = oldListener
		}

		log.Printf("[INFO] Updating NiftySetLoadBalancerSSLPoliciesOfListener")
		if v, ok := d.GetOk("ssl_policy_id"); ok {
			niftySetLoadBalancerSSLPoliciesOfListenerOpts := computing.NiftySetLoadBalancerSSLPoliciesOfListenerInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
				SSLPolicyId:      nifcloud.String(v.(string)),
			}
	
			_, err = elbconn.NiftySetLoadBalancerSSLPoliciesOfListener(&niftySetLoadBalancerSSLPoliciesOfListenerOpts)
			if err != nil {
				return fmt.Errorf("Failure NiftySetLoadBalancerSSLPoliciesOfListener for LB: %s", err)
			}	
		} else {
			niftyUnsetLoadBalancerSSLPoliciesOfListenerOpts := computing.NiftyUnsetLoadBalancerSSLPoliciesOfListenerInput{
				LoadBalancerName: nifcloud.String(d.Id()),
				InstancePort:     listener.InstancePort,
				LoadBalancerPort: listener.LoadBalancerPort,
			}
	
			_, err = elbconn.NiftyUnsetLoadBalancerSSLPoliciesOfListener(&niftyUnsetLoadBalancerSSLPoliciesOfListenerOpts)
			if err != nil {
				return fmt.Errorf("Failure NiftyUnsetLoadBalancerSSLPoliciesOfListener for LB: %s", err)
			}
		}

		d.SetPartial("ssl_policy_id")
	}

	if requestUpdate {
		log.Printf("[INFO] Updating LoadBalancer %s : update param %v", d.Id(), req)

		err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			_, err := elbconn.UpdateLoadBalancer(req)

			// Retry for ...
			if isNifcloudErr(err, "Client.InvalidParameterNotFound.LoadBalancerPort", "") {
				return resource.RetryableError(err)
			}

			if err != nil {
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if isResourceTimeoutError(err) {
			_, err = elbconn.UpdateLoadBalancer(req)
		}

		if err != nil {
			return fmt.Errorf("Error updating LoadBalancer %s: %s", d.Id(), err)
		}
		
		d.SetId(d.Get("name").(string))

//		log.Printf("[DEBUG] Waiting for DB Instance (%s) to be available", d.Id())
//		err = waitUntilNifcloudDbInstanceIsAvailableAfterUpdate(d.Id(), conn, d.Timeout(schema.TimeoutUpdate))
//		if err != nil {
//			return fmt.Errorf("error waiting for DB Instance (%s) to be available: %s", d.Id(), err)
//		}
	}

	d.Partial(false)

	return resourceNifcloudLbRead(d, meta)
}

func resourceNifcloudLbDelete(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*NifcloudClient).computingconn

	listeners, err := expandListeners(d.Get("listener").(*schema.Set).List())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting LB: %s, Listener: %v", d.Id(), listeners)

	// Destroy the load balancer
	deleteElbOpts := computing.DeleteLoadBalancerInput{
		LoadBalancerName: nifcloud.String(d.Id()),
		InstancePort:     listeners[0].InstancePort,
		LoadBalancerPort: listeners[0].LoadBalancerPort,
	}
	if _, err := elbconn.DeleteLoadBalancer(&deleteElbOpts); err != nil {
		return fmt.Errorf("Error deleting LB: %s", err)
	}

	return nil
}

/*
func resourceNifcloudLbListenerHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["instance_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-",
		strings.ToLower(m["instance_protocol"].(string))))
	buf.WriteString(fmt.Sprintf("%d-", m["lb_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-",
		strings.ToLower(m["lb_protocol"].(string))))

	if v, ok := m["ssl_certificate_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}
*/

func isLoadBalancerNotFound(err error) bool {
	elberr, ok := err.(awserr.Error)
	return ok && elberr.Code() == "Client.InvalidParameterNotFound.LoadBalancer"
}

/*
func validateHeathCheckTarget(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	// Parse the Health Check target value.
	matches := regexp.MustCompile(`\A(\w+):(\d+)(.+)?\z`).FindStringSubmatch(value)

	// Check if the value contains a valid target.
	if matches == nil || len(matches) < 1 {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check: %s",
			k, value))

		// Invalid target? Return immediately,
		// there is no need to collect other
		// errors.
		return
	}

	// Check if the value contains a valid protocol.
	if !isValidProtocol(matches[1]) {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check protocol %q. "+
				"Valid protocols are either %q, %q, %q, or %q.",
			k, matches[1], "TCP", "SSL", "HTTP", "HTTPS"))
	}

	// Check if the value contains a valid port range.
	port, _ := strconv.Atoi(matches[2])
	if port < 1 || port > 65535 {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check target port \"%d\". "+
				"Valid port is in the range from 1 to 65535 inclusive.",
			k, port))
	}

	switch strings.ToLower(matches[1]) {
	case "tcp", "ssl":
		// Check if value is in the form <PROTOCOL>:<PORT> for TCP and/or SSL.
		if matches[3] != "" {
			errors = append(errors, fmt.Errorf(
				"%q cannot contain a path in the Health Check target: %s",
				k, value))
		}

	case "http", "https":
		// Check if value is in the form <PROTOCOL>:<PORT>/<PATH> for HTTP and/or HTTPS.
		if matches[3] == "" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a path in the Health Check target: %s",
				k, value))
		}

		// Cannot be longer than 1024 multibyte characters.
		if len([]rune(matches[3])) > 1024 {
			errors = append(errors, fmt.Errorf("%q cannot contain a path longer "+
				"than 1024 characters in the Health Check target: %s",
				k, value))
		}

	}

	return
}

func isValidProtocol(s string) bool {
	if s == "" {
		return false
	}
	s = strings.ToLower(s)

	validProtocols := map[string]bool{
		"http":  true,
		"https": true,
		"ssl":   true,
		"tcp":   true,
	}

	if _, ok := validProtocols[s]; !ok {
		return false
	}

	return true
}

// LB automatically creates ENI(s) on creation
// but the cleanup is asynchronous and may take time
// which then blocks IGW, SG or VPC on deletion
// So we make the cleanup "synchronous" here
func cleanupELBNetworkInterfaces(conn *ec2.EC2, name string) error {
	out, err := conn.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   nifcloud.String("attachment.instance-owner-id"),
				Values: []*string{nifcloud.String("amazon-elb")},
			},
			{
				Name:   nifcloud.String("description"),
				Values: []*string{nifcloud.String("LB " + name)},
			},
		},
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Found %d ENIs to cleanup for LB %q",
		len(out.NetworkInterfaces), name)

	if len(out.NetworkInterfaces) == 0 {
		// Nothing to cleanup
		return nil
	}

	err = detachNetworkInterfaces(conn, out.NetworkInterfaces)
	if err != nil {
		return err
	}

	err = deleteNetworkInterfaces(conn, out.NetworkInterfaces)
	return err
}

func detachNetworkInterfaces(conn *ec2.EC2, nis []*ec2.NetworkInterface) error {
	log.Printf("[DEBUG] Trying to detach %d leftover ENIs", len(nis))
	for _, ni := range nis {
		if ni.Attachment == nil {
			log.Printf("[DEBUG] ENI %s is already detached", *ni.NetworkInterfaceId)
			continue
		}
		_, err := conn.DetachNetworkInterface(&ec2.DetachNetworkInterfaceInput{
			AttachmentId: ni.Attachment.AttachmentId,
			Force:        aws.Bool(true),
		})
		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "InvalidAttachmentID.NotFound" {
				log.Printf("[DEBUG] ENI %s is already detached", *ni.NetworkInterfaceId)
				continue
			}
			return err
		}

		log.Printf("[DEBUG] Waiting for ENI (%s) to become detached", *ni.NetworkInterfaceId)
		stateConf := &resource.StateChangeConf{
			Pending: []string{"true"},
			Target:  []string{"false"},
			Refresh: networkInterfaceAttachmentRefreshFunc(conn, *ni.NetworkInterfaceId),
			Timeout: 10 * time.Minute,
		}

		if _, err := stateConf.WaitForState(); err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "InvalidNetworkInterfaceID.NotFound" {
				continue
			}
			return fmt.Errorf(
				"Error waiting for ENI (%s) to become detached: %s", *ni.NetworkInterfaceId, err)
		}
	}
	return nil
}

func deleteNetworkInterfaces(conn *ec2.EC2, nis []*ec2.NetworkInterface) error {
	log.Printf("[DEBUG] Trying to delete %d leftover ENIs", len(nis))
	for _, ni := range nis {
		_, err := conn.DeleteNetworkInterface(&ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: ni.NetworkInterfaceId,
		})
		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "InvalidNetworkInterfaceID.NotFound" {
				log.Printf("[DEBUG] ENI %s is already deleted", *ni.NetworkInterfaceId)
				continue
			}
			return err
		}
	}
	return nil
}
*/