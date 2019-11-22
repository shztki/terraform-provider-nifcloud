resource "nifcloud_customer_gateway" "example_customer_gateway_001" {
  name                = var.customer_gateway_001["name"]
  ip_address          = var.customer_gateway_001["ip_address"]
  lan_side_ip_address = var.customer_gateway_001["lan_side_ip_address"]
  lan_side_cidr_block = var.customer_gateway_001["lan_side_cidr_block"]
  description         = var.customer_gateway_001["memo"]
}

#resource "nifcloud_customer_gateway" "example_customer_gateway_002" {
#  name                = var.customer_gateway_002["name"]
#  ip_address          = var.customer_gateway_002["ip_address"]
#  lan_side_ip_address = var.customer_gateway_002["lan_side_ip_address"]
#  lan_side_cidr_block = var.customer_gateway_002["lan_side_cidr_block"]
#  description         = var.customer_gateway_002["memo"]
#}

resource "nifcloud_vpn_gateway" "example_vpn_gateway_001" {
  name               = var.vpn_gateway_001["name"]
  private_ip_address = var.vpn_gateway_001["private_ip_address"]
  vpn_gateway_type   = var.vpn_gateway_001["vpn_gateway_type"]
  accounting_type    = var.charge_type
  availability_zone  = var.default_zone
  network_id         = nifcloud_network.example_privatelan_001.id
  security_groups    = [nifcloud_securitygroup.example_firewallgroup_003.name]
  description        = var.vpn_gateway_001["memo"]
  depends_on         = [nifcloud_router.example_router_001]
}

resource "nifcloud_vpn_connection" "example_vpn_connection_001" {
  vpn_gateway_id      = nifcloud_vpn_gateway.example_vpn_gateway_001.id
  customer_gateway_id = nifcloud_customer_gateway.example_customer_gateway_001.id
  type                = var.vpn_connection_001["type"] # IPsec | L2TPv3 /IPsec | IPsec VTI
  description         = var.vpn_connection_001["memo"]
  ipsec {
    dh                   = 2        # 2 (1024-bit MODP group) | 5 (1536-bit MODP group) | 14 (2048-bit MODP group) | 15 (3072-bit MODP group) | 16 (4096-bit MODP group) | 17 (6144-bit MODP group) | 18 (8192-bit MODP group) | 19 (256-bit random ECP group) | 20 (384-bit random ECP group) | 21 (512-bit random ECP group) | 22 (1024-bit MODP Group with 160-bit Prime Order Subgroup) | 23 (2048-bit MODP Group with 224-bit Prime Order Subgroup) | 24 (2048-bit MODP Group with 256-bit Prime Order Subgroup) | 25 (192-bit Random ECP Group) | 26 (224-bit Random ECP Group)
    esp_life_time        = 3600     # 30～86400
    ike_life_time        = 28800    # 30～86400
    encryption_algorithm = "AES256" # AES128 | AES256 | 3DES
    hash_algorithm       = "SHA256" # SHA1 | MD5 | SHA256 | SHA384 | SHA512
    ike_version          = "IKEv2"  # IKEv1 | IKEv2
    pre_shared_key       = var.pre_shared_key_001
  }
  depends_on = [nifcloud_vpn_gateway.example_vpn_gateway_001, nifcloud_customer_gateway.example_customer_gateway_001]
  lifecycle {
    ignore_changes = [ipsec, tunnel]
  }
}

#resource "nifcloud_vpn_connection" "example_vpn_connection_002" {
#  vpn_gateway_id      = nifcloud_vpn_gateway.example_vpn_gateway_001.id
#  customer_gateway_id = nifcloud_customer_gateway.example_customer_gateway_002.id
#  type                = var.vpn_connection_002["type"] # IPsec | L2TPv3 /IPsec | IPsec VTI
#  description         = var.vpn_connection_002["memo"]
#  ipsec {
#    #    dh                   = 2
#    #    esp_life_time        = 3600     # 30～86400
#    #    ike_life_time        = 28800    # 30～86400
#    #    encryption_algorithm = "AES256" # AES128 | AES256 | 3DES
#    #    hash_algorithm       = "SHA1"   # SHA1 | MD5 | SHA256 | SHA384 | SHA512
#    #    ike_version          = "IKEv2"  # IKEv1 | IKEv2
#    #    pre_shared_key       = "test1023"
#  }
#  tunnel {
#    type             = "L2TPv3"    # L2TPv3
#    mode             = "Unmanaged" # Unmanaged | Managed
#    encapsulation    = "UDP"       # IP (NiftyTunnel.ModeがUnmanagedの場合のみ指定可) | UDP
#    mtu              = "1450"
#    peer_session_id  = "1111" # modeがUnmanagedの場合
#    peer_tunnel_id   = "2222" # modeがUnmanagedの場合
#    session_id       = "1111" # modeがUnmanagedの場合
#    tunnel_id        = "2222" # modeがUnmanagedの場合
#    destination_port = "1702" # modeがUnmanagedかつencapsulationがUDPの場合
#    source_port      = "1702" # modeがUnmanagedかつencapsulationがUDPの場合
#  }
#  #depends_on = [nifcloud_vpn_gateway.example_vpn_gateway_001"[nifcloud_customer_gateway.example_customer_gateway_002]
#  lifecycle {
#    ignore_changes = [ipsec, tunnel]
#  }
#}

resource "nifcloud_route_table" "example_route_table_002" {}

resource "nifcloud_route" "example_route_002_fortable002" {
  destination_cidr_block = var.privatelan_example2["cidr"]
  route_table_id         = nifcloud_route_table.example_route_table_002.id
  ip_address             = var.router_001["ipaddress1"]
  #network_id             = nifcloud_network.example_privatelan_001.id # 値 : net-COMMON_GLOBAL (共通グローバル) | net-COMMON_PRIVATE (共通プライベート) | プライベートLAN のネットワーク ID
  #network_name           = "example001"
}

resource "nifcloud_route_table_association_with_vpn_gateway" "example_rta_002_forvpngateway001" {
  vpn_gateway_id = nifcloud_vpn_gateway.example_vpn_gateway_001.id
  route_table_id = nifcloud_route_table.example_route_table_002.id
  depends_on     = [nifcloud_vpn_gateway.example_vpn_gateway_001, nifcloud_vpn_connection.example_vpn_connection_001]
}
