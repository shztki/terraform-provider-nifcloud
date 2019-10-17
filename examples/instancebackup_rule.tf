resource "nifcloud_instancebackup_rule" "example_backup_rule_001" {
  name                      = "${lookup(var.backup_cent_001, "name")}"
  backup_instance_max_count = "${lookup(var.backup_cent_001, "max_count")}"
  time_slot_id              = "${lookup(var.backup_cent_001, "time_slot")}"
  instance_unique_id        = ["${nifcloud_instance.example_server_cent.0.unique_id}"]
  description               = "${lookup(var.backup_cent_001, "memo")}"
  depends_on                = ["nifcloud_volume.example_volume_001"]
}

resource "nifcloud_instancebackup_rule" "example_backup_rule_002" {
  name                      = "${lookup(var.backup_win_001, "name")}"
  backup_instance_max_count = "${lookup(var.backup_win_001, "max_count")}"
  time_slot_id              = "${lookup(var.backup_win_001, "time_slot")}"
  instance_unique_id        = ["${nifcloud_instance.example_server_win_noneglobal.0.unique_id}"]
  description               = "${lookup(var.backup_win_001, "memo")}"
  depends_on                = ["nifcloud_volume.example_volume_002"]
}
