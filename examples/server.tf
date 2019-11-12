resource "nifcloud_instance" "example_server_cent" {
  count    = "${lookup(var.instance_cent, "count")}"
  name     = "${format("%s%03d", "${lookup(var.instance_cent, "name")}", count.index + 1)}"
  image_id = "${lookup(var.instance_cent, "imageid")}"
  key_name = "${lookup(var.sshkey_example, "name")}"

  #ip_type             = "static" # static | elastic | none
  #public_ip           = "" # elastic
  network_interfaces {
    network_id = "net-COMMON_GLOBAL" # net-COMMON_GLOBAL | net-COMMON_PRIVATE
  }
  network_interfaces {
    network_name = "${lookup(var.privatelan_example2, "name")}"
    ipaddress    = "static"
  }

  instance_type     = "${lookup(var.instance_cent, "server_type")}"
  availability_zone = "${var.default_zone}"
  accounting_type   = "${var.charge_type}"
  security_groups   = ["${nifcloud_securitygroup.example_firewallgroup_004.name}"]
  description       = "${format("%s%03d", "${lookup(var.instance_cent, "memo")}", count.index + 1)}"
  user_data         = "${file(format("%s_%03d.sh", "${lookup(var.instance_cent, "user_data")}", count.index + 1))}"
  depends_on        = ["nifcloud_network.example_privatelan_002"]
}

resource "nifcloud_instance" "example_server_kanri" {
  count    = "${lookup(var.instance_kanri, "count")}"
  name     = "${format("%s%03d", "${lookup(var.instance_kanri, "name")}", count.index + 1)}"
  image_id = "${lookup(var.instance_kanri, "imageid")}"
  key_name = "${lookup(var.sshkey_example, "name")}"

  #ip_type             = "static" # static | elastic | none
  #public_ip           = "" # elastic
  network_interfaces {
    network_id = "net-COMMON_GLOBAL" # net-COMMON_GLOBAL | net-COMMON_PRIVATE
  }
  network_interfaces {
    network_name = "${lookup(var.privatelan_example1, "name")}"
    ipaddress    = "static"
  }

  instance_type     = "${lookup(var.instance_kanri, "server_type")}"
  availability_zone = "${var.default_zone}"
  accounting_type   = "${var.charge_type}"
  security_groups   = ["${nifcloud_securitygroup.example_firewallgroup_006.name}"]
  description       = "${format("%s%03d", "${lookup(var.instance_kanri, "memo")}", count.index + 1)}"
  user_data         = "${file(format("%s_%03d.sh", "${lookup(var.instance_kanri, "user_data")}", count.index + 1))}"
  depends_on        = ["nifcloud_network.example_privatelan_001"]
}

#resource "nifcloud_instance" "example_server_win" {
#  count    = "${lookup(var.instance_win, "count")}"
#  name     = "${format("%s%03d", "${lookup(var.instance_win, "name")}", count.index + 1)}"
#  image_id = "${lookup(var.instance_win, "imageid")}"
#  admin    = "${var.admin_user_name}"
#  password = "${var.def_pass}"
#
#  #ip_type             = "static" # static | elastic | none
#  #public_ip           = "" # elastic
#  network_interfaces {
#    network_id = "net-COMMON_GLOBAL" # net-COMMON_GLOBAL | net-COMMON_PRIVATE
#  }
#  network_interfaces {
#    //network_id = "net-COMMON_PRIVATE"
#    network_name = "${lookup(var.privatelan_example2, "name")}"
#    ipaddress    = "static"
#  }
#
#  instance_type     = "${lookup(var.instance_win, "server_type")}"
#  availability_zone = "${var.default_zone}"
#  accounting_type   = "${var.charge_type}"
#  security_groups   = ["${nifcloud_securitygroup.example_firewallgroup_004.name}"]
#  description       = "${format("%s%03d", "${lookup(var.instance_win, "memo")}", count.index + 1)}"
#  user_data         = "${file(format("%s_%03d.bat", "${lookup(var.instance_win, "user_data")}", count.index + 1))}"
#  depends_on        = ["nifcloud_network.example_privatelan_002"]
#}