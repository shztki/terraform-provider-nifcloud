resource "nifcloud_instance" "example_server_cent_noneglobal" {
  count    = "${lookup(var.instance_cent, "count")}"
  name     = "${format("%s%03d", "${lookup(var.instance_cent, "name")}", count.index + 11)}"
  image_id = "${lookup(var.instance_cent, "imageid")}"
  key_name = "${lookup(var.sshkey_example, "name")}"

  network_interfaces {
    network_name = "${lookup(var.privatelan_example, "name")}"
    ipaddress    = "static"
  }

  instance_type     = "${lookup(var.instance_cent, "server_type")}"
  availability_zone = "${var.default_zone}"
  accounting_type   = "${var.charge_type}"
  security_groups   = ["${nifcloud_securitygroup.example_firewallgroup_002.name}"]
  description     = "${format("%s%03d", "${lookup(var.instance_cent, "memo")}", count.index + 11)}"
  user_data       = "${base64encode(file(format("%s_%03d.sh", "${lookup(var.instance_cent, "user_data")}", count.index + 11)))}"
  depends_on      = ["nifcloud_network.example_privatelan_001"]
}

resource "nifcloud_instance" "example_server_win_noneglobal" {
  count    = "${lookup(var.instance_win, "count")}"
  name     = "${format("%s%03d", "${lookup(var.instance_win, "name")}", count.index + 11)}"
  image_id = "${lookup(var.instance_win, "imageid")}"
  admin    = "${var.admin_user_name}"
  password = "${var.def_pass}"

  network_interfaces {
    network_name = "${lookup(var.privatelan_example, "name")}"
    ipaddress    = "static"
  }

  instance_type     = "${lookup(var.instance_win, "server_type")}"
  availability_zone = "${var.default_zone}"
  accounting_type   = "${var.charge_type}"
  security_groups   = ["${nifcloud_securitygroup.example_firewallgroup_002.name}"]
  description     = "${format("%s%03d", "${lookup(var.instance_win, "memo")}", count.index + 11)}"
  user_data       = "${base64encode(file(format("%s_%03d.bat", "${lookup(var.instance_win, "user_data")}", count.index + 11)))}"
  depends_on      = ["nifcloud_network.example_privatelan_001"]
}
