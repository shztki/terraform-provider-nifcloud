resource "nifcloud_network" "example_privatelan_001" {
  name              = var.privatelan_example1["name"]
  cidr_block        = var.privatelan_example1["cidr"]
  availability_zone = var.default_zone
  accounting_type   = var.charge_type
  description       = var.privatelan_example1["memo"]
}

resource "nifcloud_network" "example_privatelan_002" {
  name              = var.privatelan_example2["name"]
  cidr_block        = var.privatelan_example2["cidr"]
  availability_zone = var.default_zone
  accounting_type   = var.charge_type
  description       = var.privatelan_example2["memo"]
}
