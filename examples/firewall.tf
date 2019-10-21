# https://pfs.nifcloud.com/api/rest/AuthorizeSecurityGroupIngress.htm
resource "nifcloud_securitygroup" "example_firewallgroup_001" {
  name        = "${lookup(var.firewallgroup_example_web, "name")}"
  description = "${lookup(var.firewallgroup_example_web, "memo")}"

  #  rules {
  #    from_port       = 80
  #    to_port         = 80
  #    protocol        = "TCP"
  #    cidr_blocks     = "0.0.0.0/0"
  #    #security_groups = "Zenkai"
  #    description     = "webt"
  #    inout           = "IN"
  #  }
  #  rules {
  #    #from_port       = 80
  #    #to_port         = 80
  #    protocol        = "ANY"
  #    cidr_blocks     = "0.0.0.0/0"
  #    #security_groups = "Zenkai"
  #    description     = "all"
  #    inout           = "OUT"
  #  }
}

resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_001" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    from_port   = 80
    to_port     = 80
    protocol    = "TCP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "tcp80"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_002" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    from_port   = 80
    to_port     = 80
    protocol    = "TCP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "tcp80"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_003" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "HTTPS"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "HTTPS"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_004" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "http"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "HTTP"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_005" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "ANY"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "ANY"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_006" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "ICMP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "ICMP"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_007" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "SSH"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "SSH"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_008" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "RDP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "RDP"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_009" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "GRE"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "GRE"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_010" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "ESP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "ESP"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_011" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "AH"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "AH"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_012" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "VRRP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "VRRP"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_013" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "L2TP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "L2TP"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_014" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "ICMPV6-all"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "ICMPv6-all"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_015" {
  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  rules {
    from_port   = 80
    to_port     = 80
    protocol    = "UDP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "udp80"
    inout       = "IN"
  }
}

resource "nifcloud_securitygroup" "example_firewallgroup_002" {
  name        = "${lookup(var.firewallgroup_example_db, "name")}"
  description = "${lookup(var.firewallgroup_example_db, "memo")}"
}

resource "nifcloud_securitygroup_rule" "example_firewallgroup_002_rule_001" {
  name = "${nifcloud_securitygroup.example_firewallgroup_002.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
    description     = "any web"
    inout           = "IN"
  }
}

resource "nifcloud_securitygroup" "example_firewallgroup_003" {
  name        = "${lookup(var.firewallgroup_example_vpn, "name")}"
  description = "${lookup(var.firewallgroup_example_vpn, "memo")}"

  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "ANY"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "Zenkai"
    description = "any"
    inout       = "IN"
  }
}
