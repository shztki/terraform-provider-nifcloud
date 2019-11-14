package nifcloud

import (
	"strings"
	
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/service/computing"

)

// Takes the result of schema.Set of strings and returns a []*string
func expandStringSet(configured *schema.Set) []*string {
	return expandStringList(configured.List())
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
func expandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, nifcloud.String(v.(string)))
		}
	}
	return vs
}

// Takes the result of flatmap.Expand for an array of listeners and
// returns ELB API compatible objects
func expandListeners(configured []interface{}) ([]*computing.RequestListenersStruct, error) {
	listeners := make([]*computing.RequestListenersStruct, 0, len(configured))

	// Loop over our configured listeners and create
	// an array of aws-sdk-go compatible objects
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		l := &computing.RequestListenersStruct{
			Protocol:         nifcloud.String(data["protocol"].(string)),
			BalancingType:    nifcloud.Int64(int64(data["balancing_type"].(int))),
		}

		if ip := data["instance_port"].(int); ip > 0 {
			l.SetInstancePort(int64(ip))
		}
		if lp := data["lb_port"].(int); lp > 0 {
			l.SetLoadBalancerPort(int64(lp))
		}

		listeners = append(listeners, l)
	}

	return listeners, nil
}

// Takes the result of flatmap.Expand for an array of listener and
// returns ELB API compatible objects
func expandListener(configured []interface{}) (*computing.RequestListenerStruct, error) {
	listeners := make([]*computing.RequestListenerStruct, 0, len(configured))

	// Loop over our configured listeners and create
	// an array of aws-sdk-go compatible objects
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		l := &computing.RequestListenerStruct{
			Protocol:         nifcloud.String(data["protocol"].(string)),
			BalancingType:    nifcloud.Int64(int64(data["balancing_type"].(int))),
		}

		if ip := data["instance_port"].(int); ip > 0 {
			l.SetInstancePort(int64(ip))
		}
		if lp := data["lb_port"].(int); lp > 0 {
			l.SetLoadBalancerPort(int64(lp))
		}

		listeners = append(listeners, l)
	}

	if len(listeners) != 1 {
		return nil, nil
	}

	return listeners[0], nil
}

// Takes the result of flatmap.Expand for an array of listeners and
// returns ELB API compatible objects
func expandRequestLoadBalancerNames(name string, configured []interface{}) ([]*computing.RequestLoadBalancerNamesStruct, error) {
	listeners := make([]*computing.RequestLoadBalancerNamesStruct, 0, len(configured))

	// Loop over our configured listeners and create
	// an array of aws-sdk-go compatible objects
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		l := &computing.RequestLoadBalancerNamesStruct{
			LoadBalancerName: nifcloud.String(name),
		}

		if ip := data["instance_port"].(int); ip > 0 {
			l.SetInstancePort(int64(ip))
		} else if p := data["protocol"].(string); p == "HTTP" {
			l.SetInstancePort(int64(80))
		} else if p := data["protocol"].(string); p == "HTTPS" {
			l.SetInstancePort(int64(443))
		} else if p := data["protocol"].(string); p == "FTP" {
			l.SetInstancePort(int64(21))
		}

		if lp := data["lb_port"].(int); lp > 0 {
			l.SetLoadBalancerPort(int64(lp))
		} else if p := data["protocol"].(string); p == "HTTP" {
			l.SetLoadBalancerPort(int64(80))
		} else if p := data["protocol"].(string); p == "HTTPS" {
			l.SetLoadBalancerPort(int64(443))
		} else if p := data["protocol"].(string); p == "FTP" {
			l.SetLoadBalancerPort(int64(21))
		}

		listeners = append(listeners, l)
	}

	return listeners, nil
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

// Expands an array of String Instance IDs into a []Instances
func expandInstanceString(list []interface{}) []*computing.RequestInstancesStruct {
	result := make([]*computing.RequestInstancesStruct, 0, len(list))
	for _, i := range list {
		result = append(result, &computing.RequestInstancesStruct{InstanceId: nifcloud.String(i.(string))})
	}
	return result
}


// Expands an array of String IpAddress IDs into a []IpAddresses
func expandAddFilter(list []interface{}) []*computing.RequestIPAddressesStruct {
	result := make([]*computing.RequestIPAddressesStruct, 0, len(list))

	// Loop over our configured listeners and create
	// an array of aws-sdk-go compatible objects
	for _, i := range list {
		l := &computing.RequestIPAddressesStruct{
			AddOnFilter: nifcloud.Bool(true),
			IPAddress:   nifcloud.String(i.(string)),
		}
		result = append(result, l)
	}
	return result
}

// Expands an array of String IpAddress IDs into a []IpAddresses
func expandDeleteFilter(list []interface{}) []*computing.RequestIPAddressesStruct {
	result := make([]*computing.RequestIPAddressesStruct, 0, len(list))

	// Loop over our configured listeners and create
	// an array of aws-sdk-go compatible objects
	for _, i := range list {
		l := &computing.RequestIPAddressesStruct{
			AddOnFilter: nifcloud.Bool(false),
			IPAddress:   nifcloud.String(i.(string)),
		}
		result = append(result, l)
	}
	return result
}


// Flattens an array of Instances into a []string
func flattenInstances(list []*computing.InstancesMemberItem) []string {
	result := make([]string, 0, len(list))
	for _, i := range list {
		result = append(result, *i.InstanceId)
	}
	return result
}

// Flattens an array of Instances into a []string
func flattenIpAddresses(list []*computing.IPAddressesMemberItem) []string {
	result := make([]string, 0, len(list))
	for _, i := range list {
		if *i.IPAddress == "*.*.*.*" {
			result = append(result, "")
		} else {
			result = append(result, *i.IPAddress)
		}
	}
	return result
}

// Flattens an array of Listeners into a []map[string]interface{}
func flattenListeners(list []*computing.ListenerDescriptionsMemberItem) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"instance_port":  *i.Listener.InstancePort,
			"protocol":       strings.ToUpper(*i.Listener.Protocol),
			"lb_port":        *i.Listener.LoadBalancerPort,
			"balancing_type": *i.Listener.BalancingType,
		}
		result = append(result, l)
	}
	return result
}

// Flattens an array of Listeners into a string
func flattenSSLCertificateID(list []*computing.ListenerDescriptionsMemberItem) []string {
	result := make([]string, 0, len(list))
	for _, i := range list {
		var l string
		if i.Listener.SSLCertificateId != nil {
			l = *i.Listener.SSLCertificateId;
		}
		result = append(result, l)
	}
	return result
}

// Flattens an array of Listeners into a string
func flattenSSLPolicyID(list []*computing.ListenerDescriptionsMemberItem) []string {
	result := make([]string, 0, len(list))
	for _, i := range list {
		var l string
		if i.Listener.SSLPolicy != nil {
			l = *i.Listener.SSLPolicy.SSLPolicyId;
		}
		result = append(result, l)
	}
	return result
}

// Flattens a health check into something that flatmap.Flatten()
// can handle
func flattenHealthCheck(check *computing.HealthCheck) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	chk := make(map[string]interface{})
	chk["unhealthy_threshold"] = *check.UnhealthyThreshold
//	chk["healthy_threshold"] = *check.HealthyThreshold
	chk["target"] = *check.Target
//	chk["timeout"] = *check.Timeout
	chk["interval"] = *check.Interval

	result = append(result, chk)

	return result
}