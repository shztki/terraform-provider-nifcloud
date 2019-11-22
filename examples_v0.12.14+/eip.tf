resource "nifcloud_eip" "example_eip_cent" {
  count             = var.instance_cent["count"]
  availability_zone = var.default_zone
  nifty_private_ip  = var.eip_cent["nifty_private_ip"] # true:private | false:public
  description       = format("%s%03d", var.eip_cent["memo"], count.index + 1)
  #description = var.eip_cent["memo"]
  #instance          = element(nifcloud_instance.example_server_cent.*.name, count.index)
}

resource "nifcloud_eip" "example_eip_kanri" {
  count             = var.instance_kanri["count"]
  availability_zone = var.default_zone
  nifty_private_ip  = var.eip_kanri["nifty_private_ip"] # true:private | false:public
  description       = format("%s%03d", var.eip_kanri["memo"], count.index + 1)
  #description = var.eip_kanri["memo"]
  #instance          = nifcloud_instance.example_server_kanri[0].name
}

