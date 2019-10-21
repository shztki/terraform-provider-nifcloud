variable "default_region" {
  default = "jp-east-3"
}

variable "default_zone" {
  default = "east-31"
}

variable "charge_type" {
  default = "2" # 値 : 1 (月額課金) | 2 (従量課金) 
}

variable "admin_user_name" {
  default = "nifadmin"
}
variable "def_pass" {}

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
variable "privatelan_example" {
  default = {
    name = "example001"
    cidr = "192.168.2.0/24"
    memo = "example001"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateSecurityGroup.htm
# Modify: https://pfs.nifcloud.com/api/rest/UpdateSecurityGroup.htm
variable "firewallgroup_example_web" {
  default = {
    name = "exampleweb"
    memo = "testweb"
  }
}
variable "firewallgroup_example_db" {
  default = {
    name = "exampledb"
    memo = "testdb"
  }
}
variable "firewallgroup_example_vpn" {
  default = {
    name = "examplevpn"
    memo = "testvpn"
  }
}

# Create : https://pfs.nifcloud.com/api/rest/RunInstances.htm
# Modify : https://pfs.nifcloud.com/api/rest/ModifyInstanceAttribute.htm
variable "instance_cent" {
  default = {
    count       = "1"
    name        = "examplecent"
    imageid     = "183" # 186:win2019std, # 157:win2016std, # 183:cent7.6, # 168:ubuntu18.04
    server_type = "e-small"
    memo        = "examplecent"
    user_data   = "userdata/cent7"
  }
}
variable "instance_win" {
  default = {
    count       = "1"
    name        = "examplewin"
    imageid     = "157" # 186:win2019std, # 157:win2016std, # 183:cent7.6, # 168:ubuntu18.04
    server_type = "e-medium8"
    memo        = "examplewin"
    user_data   = "userdata/win2019"
  }
}

# Create : https://pfs.nifcloud.com/api/rest/CreateVolume.htm
# Modify : https://pfs.nifcloud.com/api/rest/ModifyVolumeAttribute.htm
variable "volume_cent" {
  default = {
    name      = "example001"
    size      = 100
    disk_type = "2"
    memo      = "example001"
  }
}
variable "volume_win" {
  default = {
    name      = "example002"
    size      = 100
    disk_type = "2"
    memo      = "example002"
  }
}

# Create: https://pfs.nifcloud.com/api/rest/CreateInstanceBackupRule.htm
# Modify: https://pfs.nifcloud.com/api/rest/ModifyInstanceBackupRuleAttribute.htm
variable "backup_cent_001" {
  default = {
    name      = "examplebackup1"
    max_count = 7   # 1-10
    time_slot = "3" # JST: 1 (0:00-1:59) | 2 (2:00-3:59) | 3 (4:00-5:59) | 4 (6:00-7:59) | 5 (8:00-9:59) | 6 (10:00-11:59) | 7 (12:00-13:59) | 8 (14:00-15:59) | 9 (16:00-17:59) | 10 (18:00-19:59) | 11 (20:00-21:59) | 12 (22:00-23:59)
    memo      = "backuptest001"
  }
}

variable "backup_win_001" {
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

# Create: https://pfs.nifcloud.com/api/rest/CreateCustomerGateway.htm
# Modify: https://pfs.nifcloud.com/api/rest/NiftyModifyCustomerGatewayAttribute.htm
variable "customer_gateway_001" { # IPSec or IPSec VTI
  default = {
    name                = "examplecg1"
    ip_address          = "0.0.0.0"
    lan_side_ip_address = ""
    lan_side_cidr_block = "192.168.201.0/24"
    memo                = "example customer gateway 001"
  }
}
variable "customer_gateway_002" { # L2TPv3/IPSec
  default = {
    name                = "examplecg2"
    ip_address          = "0.0.0.0"
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
