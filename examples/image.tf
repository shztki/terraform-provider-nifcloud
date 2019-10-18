resource "nifcloud_image" "example_image_001" {
  region_name       = "${lookup(var.image_001, "region")}"
  availability_zone = "${lookup(var.image_001, "zone")}"
  name              = "${lookup(var.image_001, "name")}"
  left_instance     = "${lookup(var.image_001, "left")}"
  instance_id       = "${lookup(var.image_001, "instance")}"
  description       = "${lookup(var.image_001, "memo")}"
  #instance_id       = "${nifcloud_instance.example_server_cent.0.name}"
  #depends_on        = ["nifcloud_instance.example_server_cent.0"]
}
