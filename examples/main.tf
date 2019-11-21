provider "nifcloud" {
  region = "${var.default_region}"
}

terraform {
  required_version = "<= 0.12.13"
}
