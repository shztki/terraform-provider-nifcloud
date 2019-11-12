resource "nifcloud_volume" "example_volume_001" {
  count           = "${lookup(var.instance_cent, "count")}"
  name            = "${format("%s%03d", "${lookup(var.volume_cent, "name")}", count.index + 1)}"
  size            = "${lookup(var.volume_cent, "size")}"
  disk_type       = "${lookup(var.volume_cent, "disk_type")}"
  instance_id     = "${element(nifcloud_instance.example_server_cent.*.name, count.index)}"
  accounting_type = "${var.charge_type}"
  description     = "${format("%s%03d", "${lookup(var.volume_cent, "memo")}", count.index + 1)}"
}

#resource "nifcloud_volume" "example_volume_002" {
#  name            = "${lookup(var.volume_win, "name")}"
#  size            = "${lookup(var.volume_win, "size")}"
#  disk_type       = "${lookup(var.volume_win, "disk_type")}"
#  instance_id     = "${nifcloud_instance.example_server_win_noneglobal.0.name}"
#  accounting_type = "${var.charge_type}"
#  description     = "${lookup(var.volume_win, "memo")}"
#}
