resource "nifcloud_db_parameter_group" "example_db_parameter_group_001" {
  name        = "${lookup(var.db_param_mariadb_10_1, "name")}"
  family      = "${lookup(var.db_param_mariadb_10_1, "family")}"
  description = "${lookup(var.db_param_mariadb_10_1, "memo")}"

  parameter {
    name  = "binlog_cache_size"
    value = "65536"
    #apply_method = "immediate" # default:immediate
  }
  parameter {
    name  = "character_set_client"
    value = "utf8mb4"
    #apply_method = "immediate" # default:immediate
  }
}

resource "nifcloud_db_security_group" "example_db_security_group_001" {
  name              = "${lookup(var.db_security_001, "name")}"
  description       = "${lookup(var.db_security_001, "memo")}"
  availability_zone = "${var.default_zone}"

  ingress {
    cidr = "192.168.2.0/24"
    #security_group_name= "${nifcloud_securitygroup.example_firewallgroup_003.name}"
  }
  ingress {
    #cidr  = "192.168.2.0/24"
    security_group_name = "${nifcloud_securitygroup.example_firewallgroup_003.name}"
  }
}

