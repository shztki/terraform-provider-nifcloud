resource "nifcloud_network" "example_privatelan_001" {
  name              = "${lookup(var.privatelan_example, "name")}"
  cidr_block        = "${lookup(var.privatelan_example, "cidr")}"
  availability_zone = "${var.default_zone}"
  accounting_type   = "${var.charge_type}"
  description       = "${lookup(var.privatelan_example, "memo")}"
}
