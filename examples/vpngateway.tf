resource "nifcloud_customer_gateway" "example_customer_gateway_001" {
  name                = "${lookup(var.customer_gateway_001, "name")}"
  ip_address          = "${lookup(var.customer_gateway_001, "ip_address")}"
  lan_side_ip_address = "${lookup(var.customer_gateway_001, "lan_side_ip_address")}"
  lan_side_cidr_block = "${lookup(var.customer_gateway_001, "lan_side_cidr_block")}"
  description         = "${lookup(var.customer_gateway_001, "memo")}"
}

