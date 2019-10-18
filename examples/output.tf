output "public_ip_cent" {
  value = "${nifcloud_instance.example_server_cent.*.ip_address}"
}
output "private_ip_cent_noneglobal" {
  value = "${nifcloud_instance.example_server_cent_noneglobal.*.private_ip_address}"
}
output "public_ip_win" {
  value = "${nifcloud_instance.example_server_win.*.ip_address}"
}
output "private_ip_win_noneglobal" {
  value = "${nifcloud_instance.example_server_win_noneglobal.*.private_ip_address}"
}
