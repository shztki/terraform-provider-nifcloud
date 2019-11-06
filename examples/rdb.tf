resource "nifcloud_db_parameter_group" "example_db_parameter_group_001" {
  name        = "${lookup(var.db_param_001, "name")}"
  family      = "${lookup(var.db_param_001, "family")}"
  description = "${lookup(var.db_param_001, "memo")}"

  parameter {
    name  = "character_set_client"
    value = "utf8mb4"
    #apply_method = "immediate" # default:immediate
  }
  parameter {
    name  = "character_set_connection"
    value = "utf8mb4"
    #apply_method = "immediate" # default:immediate
  }
  parameter {
    name  = "character_set_database"
    value = "utf8mb4"
    #apply_method = "immediate" # default:immediate
  }
  parameter {
    name  = "character_set_results"
    value = "utf8mb4"
    #apply_method = "immediate" # default:immediate
  }
  parameter {
    name  = "character_set_server"
    value = "utf8mb4"
    #apply_method = "immediate" # default:immediate
  }
  parameter {
    name  = "collation_connection"
    value = "utf8mb4_general_ci"
    #apply_method = "immediate" # default:immediate
  }
  parameter {
    name  = "collation_server"
    value = "utf8mb4_general_ci"
    #apply_method = "immediate" # default:immediate
  }
}

resource "nifcloud_db_security_group" "example_db_security_group_001" {
  name              = "${lookup(var.db_security_001, "name")}"
  description       = "${lookup(var.db_security_001, "memo")}"
  availability_zone = "${var.default_zone}"

  ingress {
    cidr = "192.168.2.0/24"
    #security_group_name= "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  }
  ingress {
    #cidr  = "192.168.2.0/24"
    security_group_name = "${nifcloud_securitygroup.example_firewallgroup_001.name}"
  }
}

resource "nifcloud_db_instance" "example_db_001" {
  #replicate_source_db = "${lookup(var.db_001, "replicate_source_db")}"
  #snapshot_identifier = "${lookup(var.db_001, "snapshot_identifier")}"
  identifier        = "${lookup(var.db_001, "identifier")}"
  allocated_storage = "${lookup(var.db_001, "allocated_storage")}"
  storage_type      = "${lookup(var.db_001, "storage_type")}"
  password          = "${var.def_pass}"

  name           = "${lookup(var.db_001, "name")}"     # DB Name
  username       = "${lookup(var.db_001, "username")}" # DB Master User
  engine         = "${lookup(var.db_001, "engine")}"
  engine_version = "${lookup(var.db_001, "engine_version")}"

  availability_zone   = "${var.default_zone}"
  publicly_accessible = "${lookup(var.db_001, "publicly_accessible")}"

  instance_class          = "${lookup(var.db_001, "instance_class")}"
  backup_retention_period = "${lookup(var.db_001, "backup_retention_period")}"
  backup_window           = "${lookup(var.db_001, "backup_window")}"
  maintenance_window      = "${lookup(var.db_001, "maintenance_window")}"
  port                    = "${lookup(var.db_001, "port")}"

  security_group_names = ["${nifcloud_db_security_group.example_db_security_group_001.name}"]
  parameter_group_name = "${nifcloud_db_parameter_group.example_db_parameter_group_001.name}"

  multi_az      = "${lookup(var.db_001, "multi_az")}"
  multi_az_type = "${lookup(var.db_001, "multi_az_type")}"

  network_id      = "${nifcloud_network.example_privatelan_001.id}"
  virtual_address = "${lookup(var.db_001, "virtual_address")}"
  master_address  = "${lookup(var.db_001, "master_address")}"
  slave_address   = "${lookup(var.db_001, "slave_address")}"
  #replica_identifier = "${lookup(var.db_001, "replica_identifier")}"
  #replica_address    = "${lookup(var.db_001, "replica_address")}"

  #apply_immediately = "${lookup(var.db_001, "apply_immediately")}"
  #skip_final_snapshot = "${lookup(var.db_001, "skip_final_snapshot")}"
  #final_snapshot_identifier = "${lookup(var.db_001, "final_snapshot_identifier")}"
}

resource "nifcloud_db_instance" "example_db_002" {
  replicate_source_db = "${nifcloud_db_instance.example_db_001.id}" # "${lookup(var.db_002, "replicate_source_db")}"
  #snapshot_identifier = "${lookup(var.db_002, "snapshot_identifier")}"
  identifier = "${lookup(var.db_002, "identifier")}"
  #allocated_storage = "${lookup(var.db_002, "allocated_storage")}"
  #storage_type      = "${lookup(var.db_002, "storage_type")}"
  #password          = "${var.def_pass}"

  #name           = "${lookup(var.db_002, "name")}"     # DB Name
  #username       = "${lookup(var.db_002, "username")}" # DB Master User
  #engine         = "${lookup(var.db_002, "engine")}"
  #engine_version = "${lookup(var.db_002, "engine_version")}"

  #availability_zone   = "${var.default_zone}"
  #publicly_accessible = "${lookup(var.db_002, "publicly_accessible")}"

  instance_class = "${lookup(var.db_002, "instance_class")}"
  #backup_retention_period = "${lookup(var.db_002, "backup_retention_period")}"
  #backup_window           = "${lookup(var.db_002, "backup_window")}"
  #maintenance_window      = "${lookup(var.db_002, "maintenance_window")}"
  #port                    = "${lookup(var.db_002, "port")}"

  #security_group_names = ["${nifcloud_db_security_group.example_db_security_group_001.name}"]
  #parameter_group_name = "${nifcloud_db_parameter_group.example_db_parameter_group_001.name}"

  #multi_az      = "${lookup(var.db_002, "multi_az")}"
  #multi_az_type = "${lookup(var.db_002, "multi_az_type")}"

  #network_id      = "${nifcloud_network.example_privatelan_001.id}"
  #virtual_address = "${lookup(var.db_002, "virtual_address")}"
  #master_address  = "${lookup(var.db_002, "master_address")}"
  #slave_address = "${lookup(var.db_002, "slave_address")}"
  #replica_identifier = "${lookup(var.db_002, "replica_identifier")}"
  replica_address = "${lookup(var.db_002, "replica_address")}"

  #apply_immediately = "${lookup(var.db_002, "apply_immediately")}"
  #skip_final_snapshot = "${lookup(var.db_002, "skip_final_snapshot")}"
  #final_snapshot_identifier = "${lookup(var.db_002, "final_snapshot_identifier")}"
  lifecycle {
    ignore_changes = [availability_zone, engine, engine_version, name, network_id, security_group_names, username]
  }
}

#resource "nifcloud_db_instance" "example_db_003" {
#  #replicate_source_db = "${nifcloud_db_instance.example_db_001.id}" # "${lookup(var.db_003, "replicate_source_db")}"
#  snapshot_identifier = "${lookup(var.db_003, "snapshot_identifier")}"
#  identifier          = "${lookup(var.db_003, "identifier")}"
#  #allocated_storage = "${lookup(var.db_003, "allocated_storage")}"
#  #storage_type      = "${lookup(var.db_003, "storage_type")}"
#  #password          = "${var.def_pass}"
#
#  #name           = "${lookup(var.db_003, "name")}"     # DB Name
#  #username       = "${lookup(var.db_003, "username")}" # DB Master User
#  #engine         = "${lookup(var.db_003, "engine")}"
#  #engine_version = "${lookup(var.db_003, "engine_version")}"
#
#  availability_zone   = "${var.default_zone}"
#  publicly_accessible = "${lookup(var.db_003, "publicly_accessible")}"
#
#  instance_class = "${lookup(var.db_003, "instance_class")}"
#  #backup_retention_period = "${lookup(var.db_003, "backup_retention_period")}"
#  #backup_window           = "${lookup(var.db_003, "backup_window")}"
#  #maintenance_window      = "${lookup(var.db_003, "maintenance_window")}"
#  port = "${lookup(var.db_003, "port")}"
#
#  security_group_names = ["${nifcloud_db_security_group.example_db_security_group_001.name}"]
#  parameter_group_name = "${nifcloud_db_parameter_group.example_db_parameter_group_001.name}"
#
#  multi_az = "${lookup(var.db_003, "multi_az")}"
#  #multi_az_type = "${lookup(var.db_003, "multi_az_type")}"
#
#  #network_id      = "${nifcloud_network.example_privatelan_001.id}"
#  #virtual_address = "${lookup(var.db_003, "virtual_address")}"
#  #master_address  = "${lookup(var.db_003, "master_address")}"
#  #slave_address = "${lookup(var.db_003, "slave_address")}"
#  #replica_identifier = "${lookup(var.db_003, "replica_identifier")}"
#  #replica_address    = "${lookup(var.db_003, "replica_address")}"
#
#  #apply_immediately = "${lookup(var.db_003, "apply_immediately")}"
#  #skip_final_snapshot = "${lookup(var.db_003, "skip_final_snapshot")}"
#  #final_snapshot_identifier = "${lookup(var.db_003, "final_snapshot_identifier")}"
#  lifecycle {
#    ignore_changes = [engine, engine_version, name, username]
#  }
#}

