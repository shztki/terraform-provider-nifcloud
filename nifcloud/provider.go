package nifcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The access key for API operations.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The secret key for API operations.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The endpoint for API operations.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region where Nifcloud operations will take place.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			// "nifcloud_instance": dataSourceInstance(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nifcloud_instance":                                 resourceNifcloudInstance(),
			"nifcloud_network":                                  resourceNifcloudNetwork(),
			"nifcloud_volume":                                   resourceNifcloudVolume(),
			"nifcloud_securitygroup":                            resourceNifcloudSecurityGroup(),
			"nifcloud_securitygroup_rule":                       resourceNifcloudSecurityGroupRule(),
			"nifcloud_keypair":                                  resourceNifcloudKeyPair(),
			"nifcloud_instancebackup_rule":                      resourceNifcloudInstanceBackupRule(),
			"nifcloud_image":                                    resourceNifcloudImage(),
			"nifcloud_customer_gateway":                         resourceNifcloudCustomerGateway(),
			"nifcloud_vpn_gateway":                              resourceNifcloudVpnGateway(),
			"nifcloud_vpn_connection":                           resourceNifcloudVpnConnection(),
			"nifcloud_db_parameter_group":                       resourceNifcloudDbParameterGroup(),
			"nifcloud_db_security_group":                        resourceNifcloudDbSecurityGroup(),
			"nifcloud_db_instance":                              resourceNifcloudDbInstance(),
			"nifcloud_router":                                   resourceNifcloudRouter(),
			"nifcloud_route_table":                              resourceNifcloudRouteTable(),
			"nifcloud_route":                                    resourceNifcloudRoute(),
			"nifcloud_route_table_association":                  resourceNifcloudRouteTableAssociation(),
			"nifcloud_route_table_association_with_vpn_gateway": resourceNifcloudRouteTableAssociationWithVpnGateway(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
		Endpoint:  d.Get("endpoint").(string),
		Region:    d.Get("region").(string),
	}

	return config.Client()
}

// This is a global MutexKV for use within this plugin.
var awsMutexKV = mutexkv.NewMutexKV()