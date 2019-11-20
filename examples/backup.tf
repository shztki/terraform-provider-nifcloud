resource "nifcloud_instancebackup_rule" "example_backup_rule_cent" {
  count                     = "${lookup(var.instance_cent, "count")}"
  name                      = "${format("%s%03d", "${lookup(var.backup_cent, "name")}", count.index + 1)}"
  backup_instance_max_count = "${lookup(var.backup_cent, "max_count")}"
  time_slot_id              = "${lookup(var.backup_cent, "time_slot")}"
  instance_unique_id        = ["${element(nifcloud_instance.example_server_cent.*.unique_id, count.index)}"]
  description               = "${format("%s%03d", "${lookup(var.backup_cent, "memo")}", count.index + 1)}"
  depends_on                = ["nifcloud_volume.example_volume_cent"]
}

resource "nifcloud_instancebackup_rule" "example_backup_rule_db" {
  count                     = "${lookup(var.instance_db, "count")}"
  name                      = "${format("%s%03d", "${lookup(var.backup_db, "name")}", count.index + 1)}"
  backup_instance_max_count = "${lookup(var.backup_db, "max_count")}"
  time_slot_id              = "${lookup(var.backup_db, "time_slot")}"
  instance_unique_id        = ["${element(nifcloud_instance.example_server_cent_noneglobal.*.unique_id, count.index)}"]
  description               = "${format("%s%03d", "${lookup(var.backup_cent, "memo")}", count.index + 1)}"
  depends_on                = ["nifcloud_volume.example_volume_db"]
}

#resource "nifcloud_instancebackup_rule" "example_backup_rule_win" {
#  name                      = "${lookup(var.backup_win_001, "name")}"
#  backup_instance_max_count = "${lookup(var.backup_win_001, "max_count")}"
#  time_slot_id              = "${lookup(var.backup_win_001, "time_slot")}"
#  instance_unique_id        = ["${nifcloud_instance.example_server_win_noneglobal[0].unique_id}"]
#  description               = "${lookup(var.backup_win_001, "memo")}"
#  depends_on                = ["nifcloud_volume.example_volume_002"]
#}
