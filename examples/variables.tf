variable "default_region" {
  default = "jp-east-3"
  #default = "jp-east-1"
}

variable "default_zone" {
  default = "east-31"
  #default = "east-12"
}

variable "charge_type" {
  default = "2" # 値 : 1 (月額課金) | 2 (従量課金) 
}

variable "admin_user_name" {
  default = "nifadmin"
}
variable "def_pass" {}

variable "allow_cidr_001" {}

variable "remote_vpn_cidr" {
  default = "10.168.201.0/24"
}

variable "pre_shared_key_001" {}

# Import : https://pfs.nifcloud.com/api/rest/ImportKeyPair.htm
# Modiry : https://pfs.nifcloud.com/api/rest/NiftyModifyKeyPairAttribute.htm
variable "ssh_pubkey_path" {}
variable "sshkey_example" {
  default = {
    name = "example001"
    memo = "example001"
  }
}

# Create : https://pfs.nifcloud.com/api/rest/NiftyCreatePrivateLan.htm
# Modify : https://pfs.nifcloud.com/api/rest/NiftyModifyPrivateLanAttribute.htm
variable "privatelan_example1" {
  default = {
    name = "example001"
    cidr = "192.168.2.0/24"
    memo = "example001"
  }
}
variable "privatelan_example2" {
  default = {
    name = "example002"
    cidr = "192.168.3.0/24"
    memo = "example002"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateSecurityGroup.htm
# Modify: https://pfs.nifcloud.com/api/rest/UpdateSecurityGroup.htm
variable "firewallgroup_example_web1" {
  default = {
    name      = "exampleweb1"
    memo      = "testweb1"
    log_limit = 1000 # 1000,100000
  }
}
variable "firewallgroup_example_db1" {
  default = {
    name      = "exampledb1"
    memo      = "testdb1"
    log_limit = 1000 # 1000,100000
  }
}
variable "firewallgroup_example_vpn1" {
  default = {
    name      = "examplevpn1"
    memo      = "testvpn1"
    log_limit = 1000 # 1000,100000
  }
}
variable "firewallgroup_example_kanri1" {
  default = {
    name      = "examplekanri1"
    memo      = "testkanri1"
    log_limit = 1000 # 1000,100000
  }
}

# Create : https://pfs.nifcloud.com/api/rest/RunInstances.htm
# Modify : https://pfs.nifcloud.com/api/rest/ModifyInstanceAttribute.htm
variable "instance_cent" {
  default = {
    count       = "2"
    name        = "exampleweb"
    imageid     = "183" # 186:win2019std, # 157:win2016std, # 183:cent7.6, # 168:ubuntu18.04
    server_type = "e-small"
    memo        = "exampleweb"
    user_data   = "userdata/cent7"
  }
}
variable "instance_db" {
  default = {
    count       = "1"
    name        = "exampledb"
    imageid     = "183" # 186:win2019std, # 157:win2016std, # 183:cent7.6, # 168:ubuntu18.04
    server_type = "e-small"
    memo        = "exampledb"
    user_data   = "userdata/cent7"
  }
}
variable "instance_kanri" {
  default = {
    count       = "1"
    name        = "examplekanri"
    imageid     = "183" # 186:win2019std, # 157:win2016std, # 183:cent7.6, # 168:ubuntu18.04
    server_type = "e-small"
    memo        = "examplekanri"
    user_data   = "userdata/kanri"
  }
}
variable "instance_win" {
  default = {
    count       = "1"
    name        = "examplewin"
    imageid     = "157" # 186:win2019std, # 157:win2016std, # 183:cent7.6, # 168:ubuntu18.04
    server_type = "e-small"
    memo        = "examplewin"
    user_data   = "userdata/win2019"
  }
}

# Create : https://pfs.nifcloud.com/api/rest/CreateVolume.htm
# Modify : https://pfs.nifcloud.com/api/rest/ModifyVolumeAttribute.htm
variable "volume_cent" {
  default = {
    name      = "centdisk"
    size      = 100
    disk_type = "2" # 2 (標準ディスク) | 3 (高速ディスクA) | 4 (高速ディスクB) | 5(フラッシュドライブ)
    memo      = "centdisk"
  }
}
variable "volume_kanri" {
  default = {
    name      = "kanridisk"
    size      = 100
    disk_type = "2" # 2 (標準ディスク) | 3 (高速ディスクA) | 4 (高速ディスクB) | 5(フラッシュドライブ)
    memo      = "kanridisk"
  }
}
variable "volume_db" {
  default = {
    name      = "dbdisk"
    size      = 100
    disk_type = "2" # 2 (標準ディスク) | 3 (高速ディスクA) | 4 (高速ディスクB) | 5(フラッシュドライブ)
    memo      = "dbdisk"
  }
}
variable "volume_win" {
  default = {
    name      = "windisk"
    size      = 100 # 2 (標準ディスク) | 3 (高速ディスクA) | 4 (高速ディスクB) | 5(フラッシュドライブ)
    disk_type = "2"
    memo      = "windisk"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateInstanceBackupRule.htm
# Modify: https://pfs.nifcloud.com/api/rest/ModifyInstanceBackupRuleAttribute.htm
variable "backup_cent" {
  default = {
    name      = "backupcent"
    max_count = 1   # 1-10
    time_slot = "3" # JST: 1 (0:00-1:59) | 2 (2:00-3:59) | 3 (4:00-5:59) | 4 (6:00-7:59) | 5 (8:00-9:59) | 6 (10:00-11:59) | 7 (12:00-13:59) | 8 (14:00-15:59) | 9 (16:00-17:59) | 10 (18:00-19:59) | 11 (20:00-21:59) | 12 (22:00-23:59)
    memo      = "backupcent"
  }
}
variable "backup_db" {
  default = {
    name      = "backupdb"
    max_count = 7   # 1-10
    time_slot = "3" # JST: 1 (0:00-1:59) | 2 (2:00-3:59) | 3 (4:00-5:59) | 4 (6:00-7:59) | 5 (8:00-9:59) | 6 (10:00-11:59) | 7 (12:00-13:59) | 8 (14:00-15:59) | 9 (16:00-17:59) | 10 (18:00-19:59) | 11 (20:00-21:59) | 12 (22:00-23:59)
    memo      = "backupcent"
  }
}
variable "backup_win" {
  default = {
    name      = "examplebackup2"
    max_count = 1   # 1-10
    time_slot = "4" # JST: 1 (0:00-1:59) | 2 (2:00-3:59) | 3 (4:00-5:59) | 4 (6:00-7:59) | 5 (8:00-9:59) | 6 (10:00-11:59) | 7 (12:00-13:59) | 8 (14:00-15:59) | 9 (16:00-17:59) | 10 (18:00-19:59) | 11 (20:00-21:59) | 12 (22:00-23:59)
    memo      = "backuptest002"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateImage.htm
# Modify: https://pfs.nifcloud.com/api/rest/ModifyImageAttribute.htm
variable "image_001" {
  default = {
    name     = "exampleimage1"
    region   = "west-1"
    zone     = "west-12"
    left     = true
    instance = "testCent6"
    memo     = "testCent6 OS Image"
  }
}

# Create: https://pfs.nifcloud.com/api/rdb/CreateDBParameterGroup.htm
# Modify: https://pfs.nifcloud.com/api/rdb/ModifyDBParameterGroup.htm
variable "db_param_001" {
  default = {
    name   = "exampledb001"
    family = "mariadb10.1" # mysql5.7
    memo   = "exampledb001"
  }
}

# Create: https://pfs.nifcloud.com/api/rdb/CreateDBSecurityGroup.htm
#         https://pfs.nifcloud.com/api/rdb/AuthorizeDBSecurityGroupIngress.htm
variable "db_security_001" {
  default = {
    name = "exampledb001"
    memo = "exampledb001"
  }
}

# Create: https://pfs.nifcloud.com/api/rdb/CreateDBInstance.htm
#         https://pfs.nifcloud.com/api/rdb/CreateDBInstanceReadReplica.htm
#         https://pfs.nifcloud.com/api/rdb/RestoreDBInstanceFromDBSnapshot.htm
# Modify: https://pfs.nifcloud.com/api/rdb/ModifyDBInstance.htm
variable "db_001" {
  default = {
    name              = "testdb"   # db name
    username          = "nifadmin" # db user
    engine            = "MariaDB"  # 値：MySQL | postgres | MariaDB
    engine_version    = "10.1.18"
    allocated_storage = 50
    storage_type      = 0
    identifier        = "exampledb001" # instance name
    #replicate_source_db = "exampledb001"
    #snapshot_identifier = "exampledb-snap001"
    instance_class          = "db.mini"
    backup_retention_period = 3                     # 値：0〜10
    backup_window           = "15:00-16:00"         # UTC
    maintenance_window      = "sun:17:00-sun:18:00" # UTC
    multi_az                = true
    multi_az_type           = 0 # 値：0(データ優先) | 1(性能優先)
    port                    = 3306
    publicly_accessible     = false
    virtual_address         = "192.168.2.200/24"
    master_address          = "192.168.2.199/24"
    slave_address           = "192.168.2.198/24"

    replica_identifier = "exampledb-replica001"
    replica_address    = "192.168.2.197/24"

    #apply_immediately   = true
    #skip_final_snapshot = false
    #final_snapshot_identifier = "final_example_snap001"
  }
}
variable "db_002" {
  default = {
    #name              = "testdb"   # db name
    #username          = "nifadmin" # db user
    #engine            = "MySQL"    # 値：MySQL | postgres | MariaDB
    #engine_version    = "5.7.15"
    #allocated_storage = 50
    #storage_type      = 0
    identifier          = "exampledb-replica003" # instance name
    replicate_source_db = "exampledb001"
    #snapshot_identifier = "exampledb-snap001"
    instance_class = "db.mini"
    #backup_retention_period = 3                     # 値：0〜10
    #backup_window           = "15:00-16:00"         # UTC
    #maintenance_window      = "sun:17:00-sun:18:00" # UTC
    #multi_az                = true
    #multi_az_type           = 1 # 値：0(データ優先) | 1(性能優先)
    #port                    = 3306
    #publicly_accessible     = false
    #virtual_address         = "192.168.2.200/24"
    #master_address          = "192.168.2.199/24"
    #slave_address           = "192.168.2.198/24"

    #replica_identifier = "exampledb-replica001"
    replica_address = "192.168.2.196/24"

    #apply_immediately   = true
    #skip_final_snapshot = false
    #final_snapshot_identifier = "final_example_snap001"
  }
}
variable "db_003" {
  default = {
    #name              = "testdb"   # db name
    #username          = "nifadmin" # db user
    #engine            = "MySQL"    # 値：MySQL | postgres | MariaDB
    #engine_version    = "5.7.15"
    #allocated_storage = 50
    storage_type = 0
    identifier   = "exampledbfromsnap004" # instance name
    #replicate_source_db = "exampledb001"
    snapshot_identifier = "examplesnap001"
    instance_class      = "db.mini"
    #backup_retention_period = 3                     # 値：0〜10
    #backup_window           = "15:00-16:00"         # UTC
    #maintenance_window      = "sun:17:00-sun:18:00" # UTC
    multi_az = false
    #multi_az_type           = 1 # 値：0(データ優先) | 1(性能優先)
    port                = 3306
    publicly_accessible = true
    #virtual_address         = "192.168.2.200/24"
    #master_address          = "192.168.2.199/24"
    #slave_address           = "192.168.2.198/24"

    #replica_identifier = "exampledb-replica001"
    #replica_address    = "192.168.2.197/24"

    #apply_immediately   = true
    #skip_final_snapshot = false
    #final_snapshot_identifier = "final_example_snap001"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateCustomerGateway.htm
# Modify: https://pfs.nifcloud.com/api/rest/NiftyModifyCustomerGatewayAttribute.htm
variable "customer_gateway_001" { # IPSec or IPSec VTI
  default = {
    name                = "examplecg1"
    ip_address          = "1.1.1.1"
    lan_side_ip_address = ""
    lan_side_cidr_block = "192.168.201.0/24"
    memo                = "example customer gateway 001"
  }
}
variable "customer_gateway_002" { # L2TPv3/IPSec
  default = {
    name                = "examplecg2"
    ip_address          = "1.1.1.2"
    lan_side_ip_address = ""
    lan_side_cidr_block = ""
    memo                = "l2tpv3 customer gateway 002"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateVpnGateway.htm
# Modify: https://pfs.nifcloud.com/api/rest/NiftyModifyVpnGatewayAttribute.htm
#         https://pfs.nifcloud.com/api/rest/NiftyUpdateVpnGatewayNetworkInterfaces.htm
variable "vpn_gateway_001" {
  default = {
    name               = "examplevg1"
    private_ip_address = "192.168.2.254"
    vpn_gateway_type   = "small"
    memo               = "example vpn gateway 001"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateVpnConnection.htm
variable "vpn_connection_001" {
  default = {
    type = "IPsec" # IPsec or IPsec VTI 
    memo = "example vpn connection 001"
  }
}
variable "vpn_connection_002" {
  default = {
    type = "L2TPv3 / IPsec"
    memo = "example vpn connection 002"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/NiftyCreateRouter.htm
# Modify: https://pfs.nifcloud.com/api/rest/NiftyModifyRouterAttribute.htm
#         https://pfs.nifcloud.com/api/rest/NiftyUpdateRouterNetworkInterfaces.htm
variable "router_001" {
  default = {
    name        = "examplerouter1"
    router_type = "small"
    memo        = "example router 001"
    ipaddress1  = "192.168.2.250"
    ipaddress2  = "192.168.3.250"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateLoadBalancer.htm
#         https://pfs.nifcloud.com/api/rest/RegisterPortWithLoadBalancer.htm
# Modify: https://pfs.nifcloud.com/api/rest/UpdateLoadBalancer.htm
#         https://pfs.nifcloud.com/api/rest/UpdateLoadBalancerOption.htm
#         https://pfs.nifcloud.com/api/rest/ConfigureHealthCheck.htm
#         https://pfs.nifcloud.com/api/rest/RegisterInstancesWithLoadBalancer.htm
#         https://pfs.nifcloud.com/api/rest/SetFilterForLoadBalancer.htm
#         https://pfs.nifcloud.com/api/rest/SetLoadBalancerListenerSSLCertificate.htm
#         https://pfs.nifcloud.com/api/rest/NiftySetLoadBalancerSSLPoliciesOfListener.htm
variable "lb_001" {
  default = {
    name           = "examplelb1"
    network_volume = 10
    memo           = "example lb 001"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/AllocateAddress.htm
# Modify: https://pfs.nifcloud.com/api/rest/NiftyModifyAddressAttribute.htm
#         https://pfs.nifcloud.com/api/rest/AssociateAddress.htm
variable "eip_cent" {
  default = {
    nifty_private_ip = false # true:private | false:public
    memo             = "example eip cent"
  }
}

variable "eip_kanri" {
  default = {
    nifty_private_ip = false # true:private | false:public
    memo             = "example eip kanri"
  }
}
