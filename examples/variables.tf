variable "default_region" {
  default = "jp-east-3"
}

variable "default_zone" {
  default = "east-31"
}

variable "charge_type" {
  default = "2" # 値 : 1 (月額課金) | 2 (従量課金) 
}

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
    cidr = "192.168.0.0/24"
    memo = "example001"
  }
}

# https://pfs.nifcloud.com/api/rest/CreateSecurityGroup.htm
variable "firewallgroup_example" {
  default = {
    name = "exampleweb"
    memo = "test"
  }
}

variable "admin_user_name" {
  default = "nifadmin"
}

variable "def_pass" {}

# Create : https://pfs.nifcloud.com/api/rest/RunInstances.htm
# Modify : https://pfs.nifcloud.com/api/rest/ModifyInstanceAttribute.htm
variable "instance_cent" {
  default = {
    count       = "1"
    name        = "examplecent"
    imageid     = "183" # 186:win2019std, 183:cent7.6
    server_type = "e-small"
    memo        = "examplecent"
    user_data   = "userdata/cent7"
  }
}

variable "instance_win" {
  default = {
    count       = "1"
    name        = "examplewin"
    imageid     = "186" # 186:win2019std, 183:cent7.6
    server_type = "e-small"
    memo        = "examplewin"
    user_data   = "userdata/win2019"
  }
}

variable "firewall_group_web" {
  default = ["Zenkai"]
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
