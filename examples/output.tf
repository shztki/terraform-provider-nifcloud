output "public_ip_cent" {
  value = "${nifcloud_instance.example_server_cent.*.ip_address}"
}
output "public_ip_kanri" {
  value = "${nifcloud_instance.example_server_kanri.*.ip_address}"
}
output "private_ip_cent_noneglobal" {
  value = "${nifcloud_instance.example_server_cent_noneglobal.*.private_ip_address}"
}
output "public_ip_lb" {
  value = "${nifcloud_lb.example_lb_web.dns_name}"
}
#output "public_ip_win" {
#  value = "${nifcloud_instance.example_server_win.*.ip_address}"
#}
#output "private_ip_win_noneglobal" {
#  value = "${nifcloud_instance.example_server_win_noneglobal.*.private_ip_address}"
#}
output "puclic_eip_cent" {
  value = "${nifcloud_eip.example_eip_cent.*.public_ip}"
  #value = "${nifcloud_eip.example_eip_cent.*.private_ip}"
}
output "puclic_eip_kanri" {
  value = "${nifcloud_eip.example_eip_kanri.*.public_ip}"
  #value = "${nifcloud_eip.example_eip_kanri.*.private_ip}"
}
