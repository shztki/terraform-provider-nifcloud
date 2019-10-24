package nifcloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
)

func resourceNifcloudKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudKeyPairCreate,
		Read:   resourceNifcloudKeyPairRead,
		Update: resourceNifcloudKeyPairUpdate,
		Delete: resourceNifcloudKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceAwsKeyPairMigrateState,

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(6, 32),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 40),
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return strings.TrimSpace(v.(string))
					default:
						return ""
					}
				},
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNifcloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	var keyName string
	if v, ok := d.GetOk("key_name"); ok {
		keyName = v.(string)
	} else {
		keyName = resource.UniqueId()
		d.Set("key_name", keyName)
	}

	publicKey := d.Get("public_key").(string)
	req := &computing.ImportKeyPairInput{
		KeyName:           nifcloud.String(keyName),
		PublicKeyMaterial: nifcloud.String(base64Encode([]byte(publicKey))),
		Description:       nifcloud.String(d.Get("description").(string)),
	}
	resp, err := conn.ImportKeyPair(req)
	if err != nil {
		return fmt.Errorf("Error import KeyPair: %s", err)
	}

	d.SetId(*resp.KeyName)
	return resourceNifcloudKeyPairRead(d, meta)
}

func resourceNifcloudKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	if d.HasChange("description") {
		_, err := conn.NiftyModifyKeyPairAttribute(&computing.NiftyModifyKeyPairAttributeInput{
			KeyName:   nifcloud.String(d.Id()),
			Attribute: nifcloud.String("description"),
			Value:     nifcloud.String(d.Get("description").(string)),
		})
		if err != nil {
			return err
		}
	}

	return resourceNifcloudKeyPairRead(d, meta)
}

func resourceNifcloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn
	req := &computing.DescribeKeyPairsInput{
		KeyName: []*string{nifcloud.String(d.Id())},
	}
	resp, err := conn.DescribeKeyPairs(req)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "Client.InvalidParameterNotFound.KeyPair" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving KeyPair: %s", err)
	}

	for _, keyPair := range resp.KeySet {
		if *keyPair.KeyName == d.Id() {
			d.Set("key_name", keyPair.KeyName)
			d.Set("fingerprint", keyPair.KeyFingerprint)
			d.Set("description", keyPair.Description)
			return nil
		}
	}

	return fmt.Errorf("Unable to find key pair within: %#v", resp.KeySet)
}

func resourceNifcloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	_, err := conn.DeleteKeyPair(&computing.DeleteKeyPairInput{
		KeyName: nifcloud.String(d.Id()),
	})
	return err
}
