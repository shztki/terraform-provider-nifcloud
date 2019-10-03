output "public_ip_cent" {
  value = "${nifcloud_instance.example_server_cent.*.ip_address}"
}
output "public_ip_win" {
  value = "${nifcloud_instance.example_server_win.*.ip_address}"
}
