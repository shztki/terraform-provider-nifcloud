package nifcloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/service/rdb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNifcloudDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudDbInstanceCreate,
		Read:   resourceNifcloudDbInstanceRead,
		Update: resourceNifcloudDbInstanceUpdate,
		Delete: resourceNifcloudDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNifcloudDbInstanceImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
//				ForceNew:     true,
				ValidateFunc: validateDbName,
			},

			"username": {
				Type:         schema.TypeString,
				Optional:     true,
//				ForceNew:     true,
				ValidateFunc: validateUserName,
			},

			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validatePassword,
			},

			"engine": {
				Type:     schema.TypeString,
				Optional: true,
//				ForceNew: true,
				StateFunc: func(v interface{}) string {
					value := v.(string)
					return strings.ToLower(value)
				},
			},

			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
//				ForceNew: true,
			},

			"allocated_storage": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newInt, err := strconv.Atoi(new)

					if err != nil {
						return false
					}

					oldInt, err := strconv.Atoi(old)

					if err != nil {
						return false
					}

					// Allocated is higher than the configuration
					if oldInt > newInt {
						return true
					}

					return false
				},
			},

			"storage_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"identifier": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRdbIdentifier,
			},

			"replica_identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRdbIdentifier,
			},

			"replicate_source_db": {
				Type:     schema.TypeString,
				Optional: true,
			},

//			"replicas": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Elem:     &schema.Schema{Type: schema.TypeString},
//			},

			"snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
//				ForceNew: true,
			},

			"backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"backup_window": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateOnceADayWindowFormat,
			},

			"license_model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"maintenance_window": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(v interface{}) string {
					if v != nil {
						value := v.(string)
						return strings.ToLower(value)
					}
					return ""
				},
				ValidateFunc: validateOnceAWeekWindowFormat,
			},

			"multi_az": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"multi_az_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"publicly_accessible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"security_group_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"parameter_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"master_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"slave_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"replica_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"virtual_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"endpoint_public": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"skip_final_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			
			"final_snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if (len(value) < 1) || (len(value) > 255) {
						es = append(es, fmt.Errorf("%q must be between 1 and 255 characters in length", k))
					}
					if !regexp.MustCompile(`^[0-9A-Za-z-]+$`).MatchString(value) {
						es = append(es, fmt.Errorf("only alphanumeric characters and hyphens allowed in %q", k))
					}
					if !regexp.MustCompile(`^[a-z]`).MatchString(value) {
						es = append(es, fmt.Errorf("first character of %q must be a letter", k))
					}
					if regexp.MustCompile(`--`).MatchString(value) {
						es = append(es, fmt.Errorf("%q cannot contain two consecutive hyphens", k))
					}
					if regexp.MustCompile(`-$`).MatchString(value) {
						es = append(es, fmt.Errorf("%q cannot end in a hyphen", k))
					}
					return
				},
			},

			"ca_cert_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// apply_immediately is used to determine when the update modifications
			// take place.
			// See https://pfs.nifcloud.com/api/rdb/ModifyDBInstance.htm
			"apply_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
//				Computed: true,
			},
		},
	}
}

func resourceNifcloudDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).rdbconn

	// Some API calls (e.g. CreateDBInstanceReadReplica and
	// RestoreDBInstanceFromDBSnapshot do not support all parameters to
	// correctly apply all settings in one pass. For missing parameters or
	// unsupported configurations, we may need to call ModifyDBInstance
	// afterwards to prevent Terraform operators from API errors or needing
	// to double apply.
	var requiresModifyDbInstance bool
	modifyDbInstanceInput := &rdb.ModifyDBInstanceInput{
		ApplyImmediately: nifcloud.Bool(true),
	}

	// Some ModifyDBInstance parameters (e.g. DBParameterGroupName) require
	// a database instance reboot to take affect. During resource creation,
	// we expect everything to be in sync before returning completion.
	var requiresRebootDbInstance bool

	var identifier string
	if v, ok := d.GetOk("identifier"); ok {
		identifier = v.(string)
	}

	if v, ok := d.GetOk("replicate_source_db"); ok {
		opts := rdb.CreateDBInstanceReadReplicaInput{
			DBInstanceIdentifier:       nifcloud.String(identifier),
			DBInstanceClass:            nifcloud.String(d.Get("instance_class").(string)),
			SourceDBInstanceIdentifier: nifcloud.String(v.(string)),
		}

		if attr, ok := d.GetOk("replica_address"); ok {
			opts.NiftyReadReplicaPrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("storage_type"); ok {
			opts.NiftyStorageType = nifcloud.Int64(int64(attr.(int)))
			requiresModifyDbInstance = true
		}

		// after modify parameter
		if attr, ok := d.GetOk("allocated_storage"); ok {
			modifyDbInstanceInput.AllocatedStorage = nifcloud.Int64(int64(attr.(int)))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("backup_retention_period"); ok {
			modifyDbInstanceInput.BackupRetentionPeriod = nifcloud.Int64(int64(attr.(int)))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("parameter_group_name"); ok {
			modifyDbInstanceInput.DBParameterGroupName = nifcloud.String(attr.(string))
			requiresModifyDbInstance = true
			requiresRebootDbInstance = true
		}
		if attr, ok := d.GetOk("backup_window"); ok {
			modifyDbInstanceInput.PreferredBackupWindow = nifcloud.String(attr.(string))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("maintenance_window"); ok {
			modifyDbInstanceInput.PreferredMaintenanceWindow = nifcloud.String(attr.(string))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("password"); ok {
			modifyDbInstanceInput.MasterUserPassword = nifcloud.String(attr.(string))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("security_group_names"); ok {
			if attr := attr.(*schema.Set); attr.Len() > 0 {
				modifyDbInstanceInput.DBSecurityGroups = expandStringSet(attr)
			}
			requiresModifyDbInstance = true
		}

		log.Printf("[DEBUG] DB Instance Replica create configuration: %#v", opts)
		_, err := conn.CreateDBInstanceReadReplica(&opts)
		if err != nil {
			return fmt.Errorf("Error creating DB Instance: %s", err)
		}
	
	} else if _, ok := d.GetOk("snapshot_identifier"); ok {
		if _, ok := d.GetOk("availability_zone"); !ok {
			return fmt.Errorf(`provider.nifcloud: nifcloud_db_instance: %s: "availability_zone": required field is not set`, d.Get("identifier").(string))
		}
		opts := rdb.RestoreDBInstanceFromDBSnapshotInput{
			DBInstanceClass:         nifcloud.String(d.Get("instance_class").(string)),
			DBInstanceIdentifier:    nifcloud.String(identifier),
			DBSnapshotIdentifier:    nifcloud.String(d.Get("snapshot_identifier").(string)),
			PubliclyAccessible:      nifcloud.Bool(d.Get("publicly_accessible").(bool)),
			AvailabilityZone:        nifcloud.String(d.Get("availability_zone").(string)),
		}

//		if attr, ok := d.GetOk("name"); ok {
//			// "Note: This parameter [DBName] doesn't apply to the MySQL, PostgreSQL, or MariaDB engines."
//			// https://docs.nifcloud.amazon.com/AmazonRDS/latest/APIReference/API_RestoreDBInstanceFromDBSnapshot.html
//			switch strings.ToLower(d.Get("engine").(string)) {
//			case "mysql", "postgres", "mariadb":
//				// skip
//			default:
//				opts.DBName = nifcloud.String(attr.(string))
//			}
//		}

		if attr, ok := d.GetOk("multi_az"); ok {
			opts.MultiAZ = nifcloud.Bool(attr.(bool))
		}

		if attr, ok := d.GetOk("multi_az_type"); ok {
			opts.NiftyMultiAZType = nifcloud.Int64(int64(attr.(int)))
		}

		if attr, ok := d.GetOk("network_id"); ok {
			opts.NiftyNetworkId = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("replica_identifier"); ok {
			opts.NiftyReadReplicaDBInstanceIdentifier = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("master_address"); ok {
			opts.NiftyMasterPrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("slave_address"); ok {
			opts.NiftySlavePrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("replica_address"); ok {
			opts.NiftyReadReplicaPrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("virtual_address"); ok {
			opts.NiftyVirtualPrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("parameter_group_name"); ok {
			opts.NiftyDBParameterGroupName = nifcloud.String(attr.(string))
		}

		if attr := d.Get("security_group_names").(*schema.Set); attr.Len() > 0 {
			var s []*string
			for _, v := range attr.List() {
				s = append(s, nifcloud.String(v.(string)))
			}
			opts.NiftyDBSecurityGroups = s
		}

		if attr, ok := d.GetOk("storage_type"); ok {
			opts.NiftyStorageType = nifcloud.Int64(int64(attr.(int)))
		}

		if attr, ok := d.GetOk("port"); ok {
			opts.Port = nifcloud.Int64(int64(attr.(int)))
		}

		// after modify parameter
		if attr, ok := d.GetOk("allocated_storage"); ok {
			modifyDbInstanceInput.AllocatedStorage = nifcloud.Int64(int64(attr.(int)))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("backup_retention_period"); ok {
			modifyDbInstanceInput.BackupRetentionPeriod = nifcloud.Int64(int64(attr.(int)))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("backup_window"); ok {
			modifyDbInstanceInput.PreferredBackupWindow = nifcloud.String(attr.(string))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("maintenance_window"); ok {
			modifyDbInstanceInput.PreferredMaintenanceWindow = nifcloud.String(attr.(string))
			requiresModifyDbInstance = true
		}
		if attr, ok := d.GetOk("password"); ok {
			modifyDbInstanceInput.MasterUserPassword = nifcloud.String(attr.(string))
			requiresModifyDbInstance = true
		}

		log.Printf("[DEBUG] DB Instance restore from snapshot configuration: %s", opts)
		_, err := conn.RestoreDBInstanceFromDBSnapshot(&opts)

		if err != nil {
			if isNifcloudErr(err, "SerializationError", "failed decoding Query response") {
				// nothing
				log.Printf("[DEBUG] DB Instance restore from snapshot error: %s", err)
			} else if isNifcloudErr(err, "InternalFailure", "System Error") {
				// nothing
				log.Printf("[DEBUG] DB Instance restore from snapshot error: %s", err)
			} else {
				return fmt.Errorf("Error creating DB Instance: %s", err)
			}
		}

	} else {
		if _, ok := d.GetOk("availability_zone"); !ok {
			return fmt.Errorf(`provider.nifcloud: nifcloud_db_instance: %s: "availability_zone": required field is not set`, d.Get("identifier").(string))
		}
		if _, ok := d.GetOk("allocated_storage"); !ok {
			return fmt.Errorf(`provider.nifcloud: nifcloud_db_instance: %s: "allocated_storage": required field is not set`, d.Get("identifier").(string))
		}
		if _, ok := d.GetOk("name"); !ok {
			return fmt.Errorf(`provider.nifcloud: nifcloud_db_instance: %s: "name": required field is not set`, d.Get("identifier").(string))
		}
		if _, ok := d.GetOk("engine"); !ok {
			return fmt.Errorf(`provider.nifcloud: nifcloud_db_instance: %s: "engine": required field is not set`, d.Get("identifier").(string))
		}
		if _, ok := d.GetOk("password"); !ok {
			return fmt.Errorf(`provider.nifcloud: nifcloud_db_instance: %s: "password": required field is not set`, d.Get("identifier").(string))
		}
		if _, ok := d.GetOk("username"); !ok {
			return fmt.Errorf(`provider.nifcloud: nifcloud_db_instance: %s: "username": required field is not set`, d.Get("identifier").(string))
		}
		opts := rdb.CreateDBInstanceInput{
			AllocatedStorage:        nifcloud.Int64(int64(d.Get("allocated_storage").(int))),
			AvailabilityZone:        nifcloud.String(d.Get("availability_zone").(string)),
			DBName:                  nifcloud.String(d.Get("name").(string)),
			DBInstanceClass:         nifcloud.String(d.Get("instance_class").(string)),
			DBInstanceIdentifier:    nifcloud.String(identifier),
			MasterUsername:          nifcloud.String(d.Get("username").(string)),
			MasterUserPassword:      nifcloud.String(d.Get("password").(string)),
			Engine:                  nifcloud.String(d.Get("engine").(string)),
			EngineVersion:           nifcloud.String(d.Get("engine_version").(string)),
			PubliclyAccessible:      nifcloud.Bool(d.Get("publicly_accessible").(bool)),
		}

		attr := d.Get("backup_retention_period")
		opts.BackupRetentionPeriod = nifcloud.Int64(int64(attr.(int)))
		if attr, ok := d.GetOk("multi_az"); ok {
			opts.MultiAZ = nifcloud.Bool(attr.(bool))
		}

		if attr, ok := d.GetOk("multi_az_type"); ok {
			opts.NiftyMultiAZType = nifcloud.Int64(int64(attr.(int)))
		}

		if attr, ok := d.GetOk("network_id"); ok {
			opts.NiftyNetworkId = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("replica_identifier"); ok {
			opts.NiftyReadReplicaDBInstanceIdentifier = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("master_address"); ok {
			opts.NiftyMasterPrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("slave_address"); ok {
			opts.NiftySlavePrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("replica_address"); ok {
			opts.NiftyReadReplicaPrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("virtual_address"); ok {
			opts.NiftyVirtualPrivateAddress = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("maintenance_window"); ok {
			opts.PreferredMaintenanceWindow = nifcloud.String(attr.(string))
		}

		if attr, ok := d.GetOk("backup_window"); ok {
			opts.PreferredBackupWindow = nifcloud.String(attr.(string))
		}

//		if attr, ok := d.GetOk("license_model"); ok {
//			opts.LicenseModel = nifcloud.String(attr.(string))
//		}

		if attr, ok := d.GetOk("parameter_group_name"); ok {
			opts.DBParameterGroupName = nifcloud.String(attr.(string))
		}

		if attr := d.Get("security_group_names").(*schema.Set); attr.Len() > 0 {
			var s []*string
			for _, v := range attr.List() {
				s = append(s, nifcloud.String(v.(string)))
			}
			opts.DBSecurityGroups = s
		}
		if attr, ok := d.GetOk("storage_type"); ok {
			opts.NiftyStorageType = nifcloud.Int64(int64(attr.(int)))
		}

		if attr, ok := d.GetOk("port"); ok {
			opts.Port = nifcloud.Int64(int64(attr.(int)))
		}

		log.Printf("[DEBUG] DB Instance create configuration: %#v", opts)
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.CreateDBInstance(&opts)
			if err != nil {
				if isNifcloudErr(err, "Client.ResourceIncorrectState.DBSecurityGroup.Processing", "") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if isResourceTimeoutError(err) {
			_, err = conn.CreateDBInstance(&opts)
		}
		if err != nil {
//			if isNifcloudErr(err, "InvalidParameterValue", "") {
//				opts.MasterUserPassword = nifcloud.String("********")
//				return fmt.Errorf("Error creating DB Instance: %s, %+v", err, opts)
//			}
			return fmt.Errorf("Error creating DB Instance: %s", err)
		}
	}

	d.SetId(d.Get("identifier").(string))

	stateConf := &resource.StateChangeConf{
		Pending:    resourceNifcloudDbInstanceCreatePendingStates,
		Target:     []string{"available"},
		Refresh:    resourceNifcloudDbInstanceStateRefreshFunc(d.Id(), conn),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}

	log.Printf("[INFO] Waiting for DB Instance (%s) to be available", d.Id())
	_, err := stateConf.WaitForState()
	if err != nil {
		return err
	}

	if requiresModifyDbInstance {
		modifyDbInstanceInput.DBInstanceIdentifier = nifcloud.String(d.Id())

		log.Printf("[INFO] DB Instance (%s) configuration requires ModifyDBInstance: %s", d.Id(), modifyDbInstanceInput)
		_, err := conn.ModifyDBInstance(modifyDbInstanceInput)
		if err != nil {
			return fmt.Errorf("error modifying DB Instance (%s): %s", d.Id(), err)
		}

		log.Printf("[INFO] Waiting for DB Instance (%s) to be available", d.Id())
		err = waitUntilNifcloudDbInstanceIsAvailableAfterUpdate(d.Id(), conn, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for DB Instance (%s) to be available: %s", d.Id(), err)
		}
	}

	if requiresRebootDbInstance {
		rebootDbInstanceInput := &rdb.RebootDBInstanceInput{
			DBInstanceIdentifier: nifcloud.String(d.Id()),
		}

		log.Printf("[INFO] DB Instance (%s) configuration requires RebootDBInstance: %s", d.Id(), rebootDbInstanceInput)
		_, err := conn.RebootDBInstance(rebootDbInstanceInput)
		if err != nil {
			return fmt.Errorf("error rebooting DB Instance (%s): %s", d.Id(), err)
		}

		log.Printf("[INFO] Waiting for DB Instance (%s) to be available", d.Id())
		err = waitUntilNifcloudDbInstanceIsAvailableAfterUpdate(d.Id(), conn, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for DB Instance (%s) to be available: %s", d.Id(), err)
		}
	}

	return resourceNifcloudDbInstanceRead(d, meta)
}

func resourceNifcloudDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	v, err := resourceNifcloudDbInstanceRetrieve(d.Id(), meta.(*NifcloudClient).rdbconn)

	if err != nil {
		return err
	}
	if v == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", v.DBName)
	d.Set("identifier", v.DBInstanceIdentifier)
	d.Set("username", v.MasterUsername)
	d.Set("engine", v.Engine)
	d.Set("engine_version", v.EngineVersion)
	d.Set("allocated_storage", v.AllocatedStorage)
	d.Set("storage_type", v.NiftyStorageType)
	d.Set("instance_class", v.DBInstanceClass)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("backup_retention_period", v.BackupRetentionPeriod)
	d.Set("backup_window", v.PreferredBackupWindow)
	d.Set("license_model", v.LicenseModel)
	d.Set("maintenance_window", v.PreferredMaintenanceWindow)
	d.Set("publicly_accessible", v.PubliclyAccessible)
	d.Set("multi_az", v.MultiAZ)
	d.Set("network_id", v.NiftyNetworkId)

	if len(v.DBParameterGroups) > 0 {
		d.Set("parameter_group_name", v.DBParameterGroups[0].DBParameterGroupName)
	}

	if v.Endpoint != nil {
		d.Set("port", v.Endpoint.Port)
		if v.Endpoint.Address != nil && v.Endpoint.Port != nil {
			d.Set("endpoint_public",
				fmt.Sprintf("%s:%d", *v.Endpoint.Address, *v.Endpoint.Port))
		}
		if v.Endpoint.NiftyPrivateAddress != nil && v.Endpoint.Port != nil {
			d.Set("endpoint",
				fmt.Sprintf("%s:%d", *v.Endpoint.NiftyPrivateAddress, *v.Endpoint.Port))
		}
	}

	d.Set("status", v.DBInstanceStatus)
//	if v.OptionGroupMemberships != nil {
//		d.Set("option_group_name", v.OptionGroupMemberships[0].OptionGroupName)
//	}

	// Create an empty schema.Set to hold all security group names
	sgn := &schema.Set{
		F: schema.HashString,
	}
	for _, v := range v.DBSecurityGroups {
		sgn.Add(*v.DBSecurityGroupName)
	}
	d.Set("security_group_names", sgn)

	// replica things
//	var replicas []string
//	for _, v := range v.ReadReplicaDBInstanceIdentifiers {
//		replicas = append(replicas, *v.ReadReplicaDBInstanceIdentifier)
//	}
//	if replicas != nil {
//		d.Set("replica_identifier", replicas[0])
//	}

	d.Set("replicate_source_db", v.ReadReplicaSourceDBInstanceIdentifier)
	d.Set("ca_cert_identifier", v.CACertificateIdentifier)

	return nil
}

func resourceNifcloudDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).rdbconn

	// multi az replica delete
	if attr, ok := d.GetOk("replica_identifier"); ok {
		log.Printf("[DEBUG] Replica DB Instance destroy: %v", attr)
		opts := rdb.DeleteDBInstanceInput{DBInstanceIdentifier: nifcloud.String(attr.(string))}
		opts.SkipFinalSnapshot = nifcloud.Bool(true)
		log.Printf("[DEBUG] Replica DB Instance destroy configuration: %v", opts)
		_, err := conn.DeleteDBInstance(&opts)
		if err != nil && !isNifcloudErr(err, "Client.InvalidParameterNotFound.DBInstance", "is already being deleted") {
			return fmt.Errorf("error deleting Replica Database Instance %q: %s", attr, err)
		}
		err = waitUntilNifcloudDbInstanceIsDeleted(attr.(string), conn, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] DB Instance destroy: %v", d.Id())

	opts := rdb.DeleteDBInstanceInput{DBInstanceIdentifier: nifcloud.String(d.Id())}

	skipFinalSnapshot := d.Get("skip_final_snapshot").(bool)
	opts.SkipFinalSnapshot = nifcloud.Bool(skipFinalSnapshot)

	if !skipFinalSnapshot {
		if name, present := d.GetOk("final_snapshot_identifier"); present {
			opts.FinalDBSnapshotIdentifier = nifcloud.String(name.(string))
		} else {
			return fmt.Errorf("DB Instance FinalSnapshotIdentifier is required when a final snapshot is required")
		}
	}

	log.Printf("[DEBUG] DB Instance destroy configuration: %v", opts)
	_, err := conn.DeleteDBInstance(&opts)

	// InvalidDBInstanceState: Instance XXX is already being deleted.
	if err != nil && !isNifcloudErr(err, "Client.InvalidParameterNotFound.DBInstance", "is already being deleted") {
		return fmt.Errorf("error deleting Database Instance %q: %s", d.Id(), err)
	}

	log.Println("[INFO] Waiting for DB Instance to be destroyed")
	return waitUntilNifcloudDbInstanceIsDeleted(d.Id(), conn, d.Timeout(schema.TimeoutDelete))
}

func waitUntilNifcloudDbInstanceIsAvailableAfterUpdate(id string, conn *rdb.Rdb, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    resourceNifcloudDbInstanceUpdatePendingStates,
		Target:     []string{"available", "storage-optimization"},
		Refresh:    resourceNifcloudDbInstanceStateRefreshFunc(id, conn),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}
	_, err := stateConf.WaitForState()
	return err
}

func waitUntilNifcloudDbInstanceIsDeleted(id string, conn *rdb.Rdb, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    resourceNifcloudDbInstanceDeletePendingStates,
		Target:     []string{},
		Refresh:    resourceNifcloudDbInstanceStateRefreshFunc(id, conn),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}
	_, err := stateConf.WaitForState()
	return err
}

func resourceNifcloudDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).rdbconn

	d.Partial(true)

	req := &rdb.ModifyDBInstanceInput{
		ApplyImmediately:     nifcloud.Bool(d.Get("apply_immediately").(bool)),
		DBInstanceIdentifier: nifcloud.String(d.Id()),
	}

	d.SetPartial("apply_immediately")

	if !nifcloud.BoolValue(req.ApplyImmediately) {
		log.Println("[INFO] Only settings updating, instance changes will be applied in next maintenance window")
	}

	requestUpdate := false
	if d.HasChange("allocated_storage") {
		d.SetPartial("allocated_storage")
		req.AllocatedStorage = nifcloud.Int64(int64(d.Get("allocated_storage").(int)))
		requestUpdate = true
	}
	if d.HasChange("backup_retention_period") {
		d.SetPartial("backup_retention_period")
		req.BackupRetentionPeriod = nifcloud.Int64(int64(d.Get("backup_retention_period").(int)))
		requestUpdate = true
	}
	if d.HasChange("instance_class") {
		d.SetPartial("instance_class")
		req.DBInstanceClass = nifcloud.String(d.Get("instance_class").(string))
		requestUpdate = true
	}
	if d.HasChange("parameter_group_name") {
		d.SetPartial("parameter_group_name")
		req.DBParameterGroupName = nifcloud.String(d.Get("parameter_group_name").(string))
		requestUpdate = true
	}
	if d.HasChange("backup_window") {
		d.SetPartial("backup_window")
		req.PreferredBackupWindow = nifcloud.String(d.Get("backup_window").(string))
		requestUpdate = true
	}
	if d.HasChange("maintenance_window") {
		d.SetPartial("maintenance_window")
		req.PreferredMaintenanceWindow = nifcloud.String(d.Get("maintenance_window").(string))
		requestUpdate = true
	}
	if d.HasChange("password") {
		d.SetPartial("password")
		req.MasterUserPassword = nifcloud.String(d.Get("password").(string))
		requestUpdate = true
	}
	if d.HasChange("multi_az") {
		d.SetPartial("multi_az")
		req.MultiAZ = nifcloud.Bool(d.Get("multi_az").(bool))
		requestUpdate = true
	}
	if d.HasChange("multi_az_type") {
		d.SetPartial("multi_az_type")
		req.NiftyMultiAZType = nifcloud.Int64(int64(d.Get("multi_az_type").(int)))
		requestUpdate = true
	}
	if d.HasChange("identifier") {
		d.SetPartial("identifier")
		req.NewDBInstanceIdentifier = nifcloud.String(d.Get("identifier").(string))
		requestUpdate = true
	}
	if d.HasChange("security_group_names") {
		if attr := d.Get("security_group_names").(*schema.Set); attr.Len() > 0 {
			req.DBSecurityGroups = expandStringSet(attr)
		}
		requestUpdate = true
	}

	log.Printf("[DEBUG] Send DB Instance Modification request: %t", requestUpdate)
	if requestUpdate {
		log.Printf("[DEBUG] DB Instance Modification request: %s", req)

		err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			_, err := conn.ModifyDBInstance(req)

			// Retry for ...
			if isNifcloudErr(err, "Client.ResourceIncorrectState.DBParameterGroup.Applying", "") {
				return resource.RetryableError(err)
			}

			if err != nil {
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if isResourceTimeoutError(err) {
			_, err = conn.ModifyDBInstance(req)
		}

		if err != nil {
			return fmt.Errorf("Error modifying DB Instance %s: %s", d.Id(), err)
		}
		
		d.SetId(d.Get("identifier").(string))

		log.Printf("[DEBUG] Waiting for DB Instance (%s) to be available", d.Id())
		err = waitUntilNifcloudDbInstanceIsAvailableAfterUpdate(d.Id(), conn, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for DB Instance (%s) to be available: %s", d.Id(), err)
		}
	}

	// separate request to promote a database
//	if d.HasChange("replica_identifier") {
//		o, n := d.GetChange("replica_identifier")
//		if o != nil {
//			return fmt.Errorf("cannot elect new source database for replication n = %v", n)
//		} else {
//			return fmt.Errorf("cannot elect new source database for replication")
//		}
//		if d.Get("replicate_source_db").(string) == "" {
//			// promote
//			opts := rdb.CreateDBInstanceReadReplicaInput{
//				DBInstanceIdentifier: nifcloud.String(d.Id()),
//			}
//			attr := d.Get("backup_retention_period")
//			opts.BackupRetentionPeriod = nifcloud.Int64(int64(attr.(int)))
//			if attr, ok := d.GetOk("backup_window"); ok {
//				opts.PreferredBackupWindow = nifcloud.String(attr.(string))
//			}
//			_, err := conn.PromoteReadReplica(&opts)
//			if err != nil {
//				return fmt.Errorf("Error promoting database: %#v", err)
//			}
//			d.Set("replicate_source_db", "")
//		} else {
//			return fmt.Errorf("cannot elect new source database for replication")
//		}
//	}

	d.Partial(false)

	return resourceNifcloudDbInstanceRead(d, meta)
}

// resourceNifcloudDbInstanceRetrieve fetches DBInstance information from the AWS
// API. It returns an error if there is a communication problem or unexpected
// error with nifcloud. When the DBInstance is not found, it returns no error and a
// nil pointer.
func resourceNifcloudDbInstanceRetrieve(id string, conn *rdb.Rdb) (*rdb.DBInstance, error) {
	opts := rdb.DescribeDBInstancesInput{
		DBInstanceIdentifier: nifcloud.String(id),
	}

	log.Printf("[DEBUG] DB Instance describe configuration: %#v", opts)

	resp, err := conn.DescribeDBInstances(&opts)
	if err != nil {
		if isNifcloudErr(err, "Client.InvalidParameterNotFound.DBInstance", "") {
			return nil, nil
		}
		return nil, fmt.Errorf("Error retrieving DB Instances: %s", err)
	}

	if len(resp.DBInstances) != 1 || resp.DBInstances[0] == nil || nifcloud.StringValue(resp.DBInstances[0].DBInstanceIdentifier) != id {
		return nil, nil
	}

	return resp.DBInstances[0], nil
}

func resourceNifcloudDbInstanceImport(
	d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Neither skip_final_snapshot nor final_snapshot_identifier can be fetched
	// from any API call, so we need to default skip_final_snapshot to true so
	// that final_snapshot_identifier is not required
	d.Set("skip_final_snapshot", true)
	return []*schema.ResourceData{d}, nil
}

func resourceNifcloudDbInstanceStateRefreshFunc(id string, conn *rdb.Rdb) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := resourceNifcloudDbInstanceRetrieve(id, conn)

		if err != nil {
			log.Printf("Error on retrieving DB Instance when waiting: %s", err)
			return nil, "", err
		}

		if v == nil {
			return nil, "", nil
		}

		if v.DBInstanceStatus != nil {
			log.Printf("[DEBUG] DB Instance status for instance %s: %s", id, *v.DBInstanceStatus)
		}

		return v, *v.DBInstanceStatus, nil
	}
}

// Database instance status: https://pfs.nifcloud.com/spec/rdb/server.htm
var resourceNifcloudDbInstanceCreatePendingStates = []string{
	"backing-up",
	"creating",
	"modifying",
	"rebooting",
	"renaming",
}

var resourceNifcloudDbInstanceDeletePendingStates = []string{
	"available",
	"failed",
	"backing-up",
	"creating",
	"deleting",
	"incompatible-parameters",
	"modifying",
	"storage-full",
}

var resourceNifcloudDbInstanceUpdatePendingStates = []string{
	"backing-up",
	"creating",
	"modifying",
	"rebooting",
	"renaming",
	"storage-full",
}

func validateDbName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 64) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and underscore allowed in %q", k))
	}
	return
}

func validateUserName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 16) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 16 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and underscore allowed in %q", k))
	}
	return
}

func validatePassword(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 41) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 41 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters allowed in %q", k))
	}
	return
}

func validateRdbIdentifier(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 1) || (len(value) > 63) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 63 characters in length", k))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q", k))
	}
	if !regexp.MustCompile(`^[a-zA-Z]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"first character of %q must be a letter", k))
	}
	if regexp.MustCompile(`--`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot contain two consecutive hyphens", k))
	}
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen", k))
	}
	return
}

func validateOnceADayWindowFormat(v interface{}, k string) (ws []string, errors []error) {
	// valid time format is "hh24:mi"
	validTimeFormat := "([0-1][0-9]|2[0-3]):([0-5][0-9])"
	validTimeFormatConsolidated := "^(" + validTimeFormat + "-" + validTimeFormat + "|)$"

	value := v.(string)
	if !regexp.MustCompile(validTimeFormatConsolidated).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must satisfy the format of \"hh24:mi-hh24:mi\"", k))
	}
	return
}

func validateOnceAWeekWindowFormat(v interface{}, k string) (ws []string, errors []error) {
	// valid time format is "ddd:hh24:mi"
	validTimeFormat := "(sun|mon|tue|wed|thu|fri|sat):([0-1][0-9]|2[0-3]):([0-5][0-9])"
	validTimeFormatConsolidated := "^(" + validTimeFormat + "-" + validTimeFormat + "|)$"

	value := strings.ToLower(v.(string))
	if !regexp.MustCompile(validTimeFormatConsolidated).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must satisfy the format of \"ddd:hh24:mi-ddd:hh24:mi\"", k))
	}
	return
}

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