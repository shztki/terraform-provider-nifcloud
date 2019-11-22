resource "nifcloud_instance" "example_server_cent_noneglobal" {
  count    = var.instance_db["count"]
  name     = format("%s%03d", var.instance_db["name"], count.index + 101)
  image_id = var.instance_db["imageid"]
  key_name = var.sshkey_example["name"]

  network_interfaces {
    network_name = var.privatelan_example2["name"]
    ipaddress    = "static"
  }

  instance_type     = var.instance_db["server_type"]
  availability_zone = var.default_zone
  accounting_type   = var.charge_type
  security_groups   = [nifcloud_securitygroup.example_firewallgroup_005.name]
  description       = format("%s%03d", var.instance_db["memo"], count.index + 101)
  user_data         = file(format("%s_%03d.sh", var.instance_db["user_data"], count.index + 101))
  depends_on        = [nifcloud_network.example_privatelan_002]
}

#resource "nifcloud_instance" "example_server_win_noneglobal" {
#  count    = var.instance_win["count"]
#  name     = format("%s%03d", var.instance_win["name"], count.index + 101)
#  image_id = var.instance_win["imageid"]
#  admin    = var.admin_user_name
#  password = var.def_pass
#
#  network_interfaces {
#    network_name = var.privatelan_example2["name"]
#    ipaddress    = "static"
#  }
#
#  instance_type     = var.instance_win["server_type"]
#  availability_zone = var.default_zone
#  accounting_type   = var.charge_type
#  security_groups   = [nifcloud_securitygroup.example_firewallgroup_005.name]
#  description       = format("%s%03d", var.instance_win["memo"], count.index + 101)
#  user_data         = file(format("%s_%03d.bat", var.instance_win["user_data"], count.index + 101))
#  depends_on        = [nifcloud_network.example_privatelan_002]
#}