resource "nifcloud_network" "example_privatelan_001" {
  name              = "${lookup(var.privatelan_example1, "name")}"
  cidr_block        = "${lookup(var.privatelan_example1, "cidr")}"
  availability_zone = "${var.default_zone}"
  accounting_type   = "${var.charge_type}"
  description       = "${lookup(var.privatelan_example1, "memo")}"
}

resource "nifcloud_network" "example_privatelan_002" {
  name              = "${lookup(var.privatelan_example2, "name")}"
  cidr_block        = "${lookup(var.privatelan_example2, "cidr")}"
  availability_zone = "${var.default_zone}"
  accounting_type   = "${var.charge_type}"
  description       = "${lookup(var.privatelan_example2, "memo")}"
}
