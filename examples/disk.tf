resource "nifcloud_volume" "example_volume_001" {
  name            = "${lookup(var.volume_cent, "name")}"
  size            = "${lookup(var.volume_cent, "size")}"
  disk_type       = "${lookup(var.volume_cent, "disk_type")}"
  instance_id     = "${nifcloud_instance.example_server_cent.0.name}"
  accounting_type = "${var.charge_type}"
  description     = "${lookup(var.volume_cent, "memo")}"
}

resource "nifcloud_volume" "example_volume_002" {
  name            = "${lookup(var.volume_win, "name")}"
  size            = "${lookup(var.volume_win, "size")}"
  disk_type       = "${lookup(var.volume_win, "disk_type")}"
  instance_id     = "${nifcloud_instance.example_server_win_noneglobal.0.name}"
  accounting_type = "${var.charge_type}"
  description     = "${lookup(var.volume_win, "memo")}"
}
