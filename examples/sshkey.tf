resource "nifcloud_keypair" "example_ssh_key" {
  key_name    = "${lookup(var.sshkey_example, "name")}"
  public_key  = "${file(var.ssh_pubkey_path)}"
  description = "${lookup(var.sshkey_example, "memo")}"
}
