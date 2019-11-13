# https://pfs.nifcloud.com/api/rest/AuthorizeSecurityGroupIngress.htm
#resource "nifcloud_securitygroup" "example_firewallgroup_001" {
#  name        = "${lookup(var.firewallgroup_example_web, "name")}"
#  description = "${lookup(var.firewallgroup_example_web, "memo")}"
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_001" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    from_port   = 8080
#    to_port     = 8080
#    protocol    = "TCP"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "tcp8080"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_002" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    from_port   = 8081
#    to_port     = 8084
#    protocol    = "UDP"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "udp8081-8084"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_003" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "HTTPS"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "HTTPS"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_004" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "http"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "HTTP"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_005" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "ANY"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "ANY"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_006" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "ICMP"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "ICMP"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_007" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "SSH"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "SSH"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_008" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "RDP"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "RDP"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_009" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "GRE"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "GRE"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_010" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "ESP"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "ESP"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_011" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "AH"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "AH"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_012" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "VRRP"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "VRRP"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_013" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "L2TP"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "L2TP"
#    inout       = "IN"
#  }
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_rule_014" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol    = "ICMPV6-all"
#    cidr_blocks = "0.0.0.0/0"
#    #security_groups = "Zenkai"
#    description = "ICMPv6-all"
#    inout       = "IN"
#  }
#}
#
#resource "nifcloud_securitygroup" "example_firewallgroup_002" {
#  name        = "${lookup(var.firewallgroup_example_db, "name")}"
#  description = "${lookup(var.firewallgroup_example_db, "memo")}"
#}
#resource "nifcloud_securitygroup_rule" "example_firewallgroup_002_rule_001" {
#  name = "${nifcloud_securitygroup.example_firewallgroup_002.name}"
#  rules {
#    #from_port       = 80
#    #to_port         = 80
#    protocol = "ANY"
#    #cidr_blocks     = "0.0.0.0/0"
#    security_groups = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
#    description     = "any web"
#    inout           = "IN"
#  }
#}

resource "nifcloud_securitygroup" "example_firewallgroup_003" {
  name                   = "${lookup(var.firewallgroup_example_vpn1, "name")}"
  description            = "${lookup(var.firewallgroup_example_vpn1, "memo")}"
  group_log_limit_update = "${lookup(var.firewallgroup_example_vpn1, "log_limit")}"
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_003_rule_001" {
  name = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "ANY"
    cidr_blocks = "${lookup(var.customer_gateway_001, "ip_address")}"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
    description = "for vpn"
    inout       = "IN"
  }
}

resource "nifcloud_securitygroup" "example_firewallgroup_004" {
  name                   = "${lookup(var.firewallgroup_example_web1, "name")}"
  description            = "${lookup(var.firewallgroup_example_web1, "memo")}"
  group_log_limit_update = "${lookup(var.firewallgroup_example_web1, "log_limit")}"
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_004_rule_001" {
  name = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "HTTP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
    description = "for web1"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_004_rule_002" {
  name = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "HTTPS"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
    description = "for web1"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_004_rule_003" {
  name = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_005.name}"
    description     = "for web1 from db1"
    inout           = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_004_rule_004" {
  name = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
    description     = "for web1 from kanri1"
    inout           = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_004_rule_005" {
  name = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
    description     = "for web1 from vpn1"
    inout           = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_004_rule_006" {
  name = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "SSH"
    cidr_blocks = "${var.allow_cidr_001}"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
    description = "for web1 from office"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_004_rule_007" {
  name = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "RDP"
    cidr_blocks = "${var.allow_cidr_001}"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
    description = "for web1 from office"
    inout       = "IN"
  }
}

resource "nifcloud_securitygroup" "example_firewallgroup_005" {
  name                   = "${lookup(var.firewallgroup_example_db1, "name")}"
  description            = "${lookup(var.firewallgroup_example_db1, "memo")}"
  group_log_limit_update = "${lookup(var.firewallgroup_example_db1, "log_limit")}"
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_005_rule_001" {
  name = "${nifcloud_securitygroup.example_firewallgroup_005.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
    description     = "for db from web1"
    inout           = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_005_rule_002" {
  name = "${nifcloud_securitygroup.example_firewallgroup_005.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
    description     = "for db from vpn1"
    inout           = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_005_rule_003" {
  name = "${nifcloud_securitygroup.example_firewallgroup_005.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
    description     = "for db from kanri1"
    inout           = "IN"
  }
}

resource "nifcloud_securitygroup" "example_firewallgroup_006" {
  name                   = "${lookup(var.firewallgroup_example_kanri1, "name")}"
  description            = "${lookup(var.firewallgroup_example_kanri1, "memo")}"
  group_log_limit_update = "${lookup(var.firewallgroup_example_kanri1, "log_limit")}"
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_006_rule_001" {
  name = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "HTTP"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
    description = "for kanri1"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_006_rule_002" {
  name = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "HTTPS"
    cidr_blocks = "0.0.0.0/0"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
    description = "for kanri1"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_006_rule_003" {
  name = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_005.name}"
    description     = "for kanri1 from db1"
    inout           = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_006_rule_004" {
  name = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_004.name}"
    description     = "for kanri1 from web1"
    inout           = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_006_rule_005" {
  name = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol = "ANY"
    #cidr_blocks     = "0.0.0.0/0"
    security_groups = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
    description     = "for kanri1 from vpn1"
    inout           = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_006_rule_006" {
  name = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "SSH"
    cidr_blocks = "${var.allow_cidr_001}"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
    description = "for kanri1 from office"
    inout       = "IN"
  }
}
resource "nifcloud_securitygroup_rule" "example_firewallgroup_006_rule_007" {
  name = "${nifcloud_securitygroup.example_firewallgroup_006.name}"
  rules {
    #from_port       = 80
    #to_port         = 80
    protocol    = "ANY"
    cidr_blocks = "${var.remote_vpn_cidr}"
    #security_groups = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
    description = "for kanri1 from remote access vpn"
    inout       = "IN"
  }
}