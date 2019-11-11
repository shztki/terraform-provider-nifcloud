package nifcloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNifcloudVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudVpnConnectionCreate,
		Read:   resourceNifcloudVpnConnectionRead,
//		Update: resourceNifcloudVpnConnectionUpdate,
		Delete: resourceNifcloudVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// IPSec
			"ipsec": {
				Type:     schema.TypeSet,
				Required: true,
//				Optional: true,
				ForceNew: true,
				MinItems: 0,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dh": {
							Type:     schema.TypeInt,
							Default:  2,
//							Required: true,
							Optional: true,
//							ForceNew: true,
						},
						"esp_life_time": {
							Type:     schema.TypeInt,
							Default:  3600,
//							Required: true,
							Optional: true,
//							ForceNew: true,
						},
						"ike_life_time": {
							Type:     schema.TypeInt,
							Default:  28800,
//							Required: true,
							Optional: true,
//							ForceNew: true,
						},
						"encryption_algorithm": {
							Type:     schema.TypeString,
							Default:  "AES128",
//							Required: true,
							Optional: true,
//							ForceNew: true,
						},
						"hash_algorithm": {
							Type:     schema.TypeString,
							Default:  "SHA1",
//							Required: true,
							Optional: true,
//							ForceNew: true,
						},
						"ike_version": {
							Type:     schema.TypeString,
							Default:  "IKEv1",
//							Required: true,
							Optional: true,
//							ForceNew: true,
						},
						"pre_shared_key": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateVpnConnectionTunnelPreSharedKey,
//							ForceNew:     true,
						},
					},
				},
			},

			// "L2TPv3 / IPSec" Only
			"tunnel": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MinItems: 0,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Default:  "L2TPv3",
//							Required: true,
							Optional: true,
//							ForceNew: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Required: true,
//							ForceNew: true,
						},
						"encapsulation": {
							Type:     schema.TypeString,
							Required: true,
//							ForceNew: true,
						},
						"mtu": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "1500",
//							ForceNew: true,
						},
						"peer_session_id": {
							Type:     schema.TypeString,
							Optional: true,
//							ForceNew: true,
						},
						"peer_tunnel_id": {
							Type:     schema.TypeString,
							Optional: true,
//							ForceNew: true,
						},
						"session_id": {
							Type:     schema.TypeString,
							Optional: true,
//							ForceNew: true,
						},
						"tunnel_id": {
							Type:     schema.TypeString,
							Optional: true,
//							ForceNew: true,
						},
						"destination_port": {
							Type:     schema.TypeString,
							Optional: true,
//							ForceNew: true,
						},
						"source_port": {
							Type:     schema.TypeString,
							Optional: true,
//							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceNifcloudVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

//	awsMutexKV.Lock(d.Get("vpn_gateway_id").(string))
//	defer awsMutexKV.Unlock(d.Get("vpn_gateway_id").(string))
//	awsMutexKV.Lock(d.Get("customer_gateway_id").(string))
//	defer awsMutexKV.Unlock(d.Get("customer_gateway_id").(string))

	createOpts := &computing.CreateVpnConnectionInput{
		Agreement:                     nifcloud.Bool(true),
		CustomerGatewayId:             nifcloud.String(d.Get("customer_gateway_id").(string)),
		NiftyVpnConnectionDescription: nifcloud.String(d.Get("description").(string)),
		VpnGatewayId:                  nifcloud.String(d.Get("vpn_gateway_id").(string)),
		Type:                          nifcloud.String(d.Get("type").(string)),
//		NiftyVpnConnectionMtu:         nifcloud.String(d.Get("mtu").(string)),
	}

	if tunnels, ok := d.GetOk("tunnel"); ok {
		tunnel := &computing.RequestNiftyTunnelStruct{}
		for _, tmp := range tunnels.(*schema.Set).List() {
			if v, ok := tmp.(map[string]interface{}); ok {
				createOpts.SetNiftyVpnConnectionMtu(v["mtu"].(string))
				tunnel.SetEncapsulation(v["encapsulation"].(string))
				tunnel.SetMode(v["mode"].(string))
				tunnel.SetType(v["type"].(string))
				if v["peer_session_id"].(string) != "" {
					tunnel.SetPeerSessionId(v["peer_session_id"].(string))
				}
				if v["peer_tunnel_id"].(string) != "" {
					tunnel.SetPeerTunnelId(v["peer_tunnel_id"].(string))
				}
				if v["session_id"].(string) != "" {
					tunnel.SetSessionId(v["session_id"].(string))
				}
				if v["destination_port"].(string) != "" {
					tunnel.SetDestinationPort(v["destination_port"].(string))
				}
				if v["source_port"].(string) != "" {
					tunnel.SetSourcePort(v["source_port"].(string))
				}
				if v["tunnel_id"].(string) != "" {
					tunnel.SetTunnelId(v["tunnel_id"].(string))
				}
			}
		}
		createOpts.SetNiftyTunnel(tunnel)

	} else if ipsecs, ok := d.GetOk("ipsec"); ok {
		ipsec := &computing.RequestNiftyIpsecConfigurationStruct{}
		for _, tmp := range ipsecs.(*schema.Set).List() {
			if v, ok := tmp.(map[string]interface{}); ok {
				ipsec.SetDiffieHellmanGroup(int64(v["dh"].(int)))
				ipsec.SetEncapsulatingSecurityPayloadLifetime(int64(v["esp_life_time"].(int)))
				ipsec.SetEncryptionAlgorithm(v["encryption_algorithm"].(string))
				ipsec.SetHashAlgorithm(v["hash_algorithm"].(string))
				ipsec.SetInternetKeyExchange(v["ike_version"].(string))
				ipsec.SetInternetKeyExchangeLifetime(int64(v["ike_life_time"].(int)))
				if v["pre_shared_key"].(string) != "" {
					ipsec.SetPreSharedKey(v["pre_shared_key"].(string))
				}
			}
		}
		createOpts.SetNiftyIpsecConfiguration(ipsec)

	} else {
		log.Printf("error VPN connection needs tunnel or ipsec parameter.")
		return nil
	}

	// Create the VPN Connection
	log.Printf("[DEBUG] Creating vpn connection")
	resp, err := conn.CreateVpnConnection(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating vpn connection: %s", err)
	}

	d.SetId(nifcloud.StringValue(resp.VpnConnection.VpnConnectionId))
	log.Printf("[DEBUG] vpn connection create %v", resp)

	if err := waitForEc2VpnConnectionAvailable(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for VPN connection (%s) to become available: %s", d.Id(), err)
	}

	// Read off the API to populate our RO fields.
	return resourceNifcloudVpnConnectionRead(d, meta)
}

func vpnConnectionRefreshFunc(conn *computing.Computing, connectionID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeVpnConnections(&computing.DescribeVpnConnectionsInput{
			VpnConnectionId: []*string{nifcloud.String(connectionID)},
		})

		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.VpnConnectionId" {
				resp = nil
			} else {
				log.Printf("Error on VPNConnectionRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil || len(resp.VpnConnectionSet) == 0 {
			return nil, "", nil
		}

		connection := resp.VpnConnectionSet[0]
		return connection, *connection.State, nil
	}
}

func resourceNifcloudVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	resp, err := conn.DescribeVpnConnections(&computing.DescribeVpnConnectionsInput{
		VpnConnectionId: []*string{nifcloud.String(d.Id())},
	})

	if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.VpnConnectionId" {
		log.Printf("[WARN] EC2 VPN Connection (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EC2 VPN Connection (%s): %s", d.Id(), err)
	}

	if resp == nil || len(resp.VpnConnectionSet) == 0 || resp.VpnConnectionSet[0] == nil {
		return fmt.Errorf("error reading EC2 VPN Connection (%s): empty response", d.Id())
	}

	if len(resp.VpnConnectionSet) > 1 {
		return fmt.Errorf("error reading EC2 VPN Connection (%s): multiple responses", d.Id())
	}

	log.Printf("[DEBUG] vpn connection describe %v", resp)
	vpnConnection := resp.VpnConnectionSet[0]

	// Set attributes under the user's control.
	d.Set("vpn_gateway_id", vpnConnection.VpnGatewayId)
	d.Set("customer_gateway_id", vpnConnection.CustomerGatewayId)
	d.Set("type", vpnConnection.Type)
	d.Set("description", vpnConnection.NiftyVpnConnectionDescription)

	if vpnConnection.NiftyIpsecConfiguration != nil {
		ipsec := vpnConnection.NiftyIpsecConfiguration
		d.Set("mtu", ipsec.Mtu)
		if err := d.Set("ipsec", ipsecToMapList(ipsec)); err != nil {
			return err
		}
	}

	if vpnConnection.NiftyTunnel != nil {
		tunnel := vpnConnection.NiftyTunnel
		if err := d.Set("tunnel", tunnelToMapList(tunnel)); err != nil {
			return err
		}
	}

	return nil
}

/*
func resourceNifcloudVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	// Update tags if required.
	if err := setTags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	return resourceNifcloudVpnConnectionRead(d, meta)
}
*/

func resourceNifcloudVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	awsMutexKV.Lock(d.Get("vpn_gateway_id").(string))
	defer awsMutexKV.Unlock(d.Get("vpn_gateway_id").(string))
	awsMutexKV.Lock(d.Get("customer_gateway_id").(string))
	defer awsMutexKV.Unlock(d.Get("customer_gateway_id").(string))

	request := &computing.DeleteVpnConnectionInput{
		VpnConnectionId: nifcloud.String(d.Id()),
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DeleteVpnConnection(request)
		log.Printf("[DEBUG] deleting vpn connection %v", resp)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "Client.InvalidParameterNotFound.VpnConnectionId" {
				return nil
			}
			return resource.RetryableError(err)
		}

		return nil
	})
/*
	_, err := conn.DeleteVpnConnection(&computing.DeleteVpnConnectionInput{
		VpnConnectionId: nifcloud.String(d.Id()),
	})

	if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.VpnConnectionId" {
		return nil
	}
*/
	if err != nil {
		return fmt.Errorf("error deleting VPN Connection (%s): %s", d.Id(), err)
	}

	return nil
}

func waitForEc2VpnConnectionAvailable(conn *computing.Computing, id string) error {
	// Wait for the connection to become available. This has an obscenely
	// high default timeout because AWS VPN connections are notoriously
	// slow at coming up or going down. There's also no point in checking
	// more frequently than every ten seconds.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    vpnConnectionRefreshFunc(conn, id),
		Timeout:    40 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}

func validateVpnConnectionTunnelPreSharedKey(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if (len(value) < 1) || (len(value) > 64) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length", k))
	}

	if strings.HasPrefix(value, "0") {
		errors = append(errors, fmt.Errorf("%q cannot start with zero character", k))
	}

	if !regexp.MustCompile(`^[0-9a-zA-Z-+&!@#$%^*(),.:_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q can only contain alphanumeric and %q characters", k, "-+&!@#$%^*(),.:_"))
	}

	return
}

func ipsecToMapList(ipsec *computing.NiftyIpsecConfiguration) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if ipsec != nil {
		tmp := make(map[string]interface{})
		tmp["dh"] = *ipsec.DiffieHellmanGroup
		tmp["esp_life_time"] = *ipsec.EncapsulatingSecurityPayloadLifetime
		tmp["ike_life_time"] = *ipsec.InternetKeyExchangeLifetime
		tmp["encryption_algorithm"] = *ipsec.EncryptionAlgorithm
		tmp["hash_algorithm"] = *ipsec.HashingAlgorithm
		tmp["ike_version"] = *ipsec.InternetKeyExchange
		tmp["pre_shared_key"] = *ipsec.PreSharedKey
		result = append(result, tmp)
	}

	return result
}

func tunnelToMapList(tunnel *computing.NiftyTunnel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if tunnel != nil {
		tmp := make(map[string]interface{})
		tmp["type"] = *tunnel.Type
		tmp["mode"] = *tunnel.Mode
		tmp["encapsulation"] = *tunnel.Encapsulation
		tmp["peer_session_id"] = *tunnel.PeerSessionId
		tmp["peer_tunnel_id"] = *tunnel.PeerTunnelId
		tmp["session_id"] = *tunnel.SessionId
		tmp["tunnel_id"] = *tunnel.TunnelId
		tmp["destination_port"] = *tunnel.DestinationPort
		tmp["source_port"] = *tunnel.SourcePort
		result = append(result, tmp)
	}

	return result
}
