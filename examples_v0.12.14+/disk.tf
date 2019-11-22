resource "nifcloud_volume" "example_volume_cent" {
  count           = var.instance_cent["count"]
  name            = format("%s%03d", var.volume_cent["name"], count.index + 1)
  size            = var.volume_cent["size"]
  disk_type       = var.volume_cent["disk_type"]
  instance_id     = element(nifcloud_instance.example_server_cent.*.name, count.index)
  accounting_type = var.charge_type
  description     = format("%s%03d", var.volume_cent["memo"], count.index + 1)
}

resource "nifcloud_volume" "example_volume_kanri" {
  count           = var.instance_kanri["count"]
  name            = format("%s%03d", var.volume_kanri["name"], count.index + 1)
  size            = var.volume_kanri["size"]
  disk_type       = var.volume_kanri["disk_type"]
  instance_id     = element(nifcloud_instance.example_server_kanri.*.name, count.index)
  accounting_type = var.charge_type
  description     = format("%s%03d", var.volume_kanri["memo"], count.index + 1)
}

resource "nifcloud_volume" "example_volume_db" {
  count           = var.instance_db["count"]
  name            = format("%s%03d", var.volume_db["name"], count.index + 1)
  size            = var.volume_db["size"]
  disk_type       = var.volume_db["disk_type"]
  instance_id     = element(nifcloud_instance.example_server_cent_noneglobal.*.name, count.index)
  accounting_type = var.charge_type
  description     = format("%s%03d", var.volume_db["memo"], count.index + 1)
}

#resource "nifcloud_volume" "example_volume_win" {
#  name            = var.volume_win["name"]
#  size            = var.volume_win["size"]
#  disk_type       = var.volume_win["disk_type"]
#  instance_id     = nifcloud_instance.example_server_win_noneglobal[0].name
#  accounting_type = var.charge_type
#  description     = var.volume_win["memo"]
#}
