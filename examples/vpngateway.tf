resource "nifcloud_customer_gateway" "example_customer_gateway_001" {
  name                = "${lookup(var.customer_gateway_001, "name")}"
  ip_address          = "${lookup(var.customer_gateway_001, "ip_address")}"
  lan_side_ip_address = "${lookup(var.customer_gateway_001, "lan_side_ip_address")}"
  lan_side_cidr_block = "${lookup(var.customer_gateway_001, "lan_side_cidr_block")}"
  description         = "${lookup(var.customer_gateway_001, "memo")}"
}

resource "nifcloud_customer_gateway" "example_customer_gateway_002" {
  name                = "${lookup(var.customer_gateway_002, "name")}"
  ip_address          = "${lookup(var.customer_gateway_002, "ip_address")}"
  lan_side_ip_address = "${lookup(var.customer_gateway_002, "lan_side_ip_address")}"
  lan_side_cidr_block = "${lookup(var.customer_gateway_002, "lan_side_cidr_block")}"
  description         = "${lookup(var.customer_gateway_002, "memo")}"
}

resource "nifcloud_vpn_gateway" "example_vpn_gateway_001" {
  name               = "${lookup(var.vpn_gateway_001, "name")}"
  private_ip_address = "${lookup(var.vpn_gateway_001, "private_ip_address")}"
  vpn_gateway_type   = "${lookup(var.vpn_gateway_001, "vpn_gateway_type")}"
  accounting_type    = "${var.charge_type}"
  availability_zone  = "${var.default_zone}"
  network_id         = "${nifcloud_network.example_privatelan_001.id}"
  security_groups    = ["${nifcloud_securitygroup.example_firewallgroup_003.name}"]
  description        = "${lookup(var.vpn_gateway_001, "memo")}"
}

