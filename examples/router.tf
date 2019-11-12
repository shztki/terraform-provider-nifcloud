resource "nifcloud_router" "example_router_001" {
  name              = "${lookup(var.router_001, "name")}"
  router_type       = "${lookup(var.router_001, "router_type")}"
  accounting_type   = "${var.charge_type}"
  availability_zone = "${var.default_zone}"
  network_interfaces {
    network_id = "${nifcloud_network.example_privatelan_001.id}" # 値 : net-COMMON_GLOBAL (共通グローバル) | net-COMMON_PRIVATE (共通プライベート) | プライベートLAN のネットワーク ID
    ipaddress  = "${lookup(var.router_001, "ipaddress1")}"
    #dhcp            = false
    #dhcp_options_id = ""
    #dhcp_config_id  = ""
  }
  network_interfaces {
    network_id = "${nifcloud_network.example_privatelan_002.id}" # 値 : net-COMMON_GLOBAL (共通グローバル) | net-COMMON_PRIVATE (共通プライベート) | プライベートLAN のネットワーク ID
    ipaddress  = "${lookup(var.router_001, "ipaddress2")}"
    #dhcp            = false
    #dhcp_options_id = ""
    #dhcp_config_id  = ""
  }
  security_groups = ["${nifcloud_securitygroup.example_firewallgroup_003.name}"]
  description     = "${lookup(var.router_001, "memo")}"
}

resource "nifcloud_route_table" "example_route_table_001" {}

resource "nifcloud_route" "example_route_001_fortable001" {
  destination_cidr_block = "${lookup(var.customer_gateway_001, "lan_side_cidr_block")}"
  route_table_id         = "${nifcloud_route_table.example_route_table_001.id}"
  ip_address             = "${lookup(var.vpn_gateway_001, "private_ip_address")}"
  #network_id             = "${nifcloud_network.example_privatelan_001.id}" # 値 : net-COMMON_GLOBAL (共通グローバル) | net-COMMON_PRIVATE (共通プライベート) | プライベートLAN のネットワーク ID
  #network_name           = "example001"
}

resource "nifcloud_route_table_association" "example_rta_001_forrouter001" {
  router_id      = "${nifcloud_router.example_router_001.id}"
  route_table_id = "${nifcloud_route_table.example_route_table_001.id}"
  depends_on     = ["nifcloud_router.example_router_001"]
}
