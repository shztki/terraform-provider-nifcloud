# https://pfs.nifcloud.com/api/rest/AuthorizeSecurityGroupIngress.htm
resource "nifcloud_securitygroup" "example_firewallgroup_001" {
  name        = "${lookup(var.firewallgroup_example, "name")}"
  description = "${lookup(var.firewallgroup_example, "memo")}"

  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "TCP"
    cidr_blocks     = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description     = "webt"
  }
  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "UDP"
    cidr_blocks     = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description     = "webu"
  }
  ingress {
    #from_port       = 80
    #to_port         = 80
    protocol        = "https"
    cidr_blocks     = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description     = "https"
  }
  ingress {
    #from_port       = 80
    #to_port         = 80
    protocol        = "SSH"
    cidr_blocks     = "1.1.1.1"
    #security_groups = "Zenkai"
    description     = "ssh"
  }
  ingress {
    #from_port       = 80
    #to_port         = 80
    protocol        = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "Zenkai"
    description     = "group-zenkai"
  }
  ingress {
    #from_port       = 80
    #to_port         = 80
    protocol        = "RDP"
    cidr_blocks     = "1.1.1.0/28"
    #security_groups = "Zenkai"
    description     = "rdp"
  }

}
