resource "nifcloud_lb" "example_lb_web" {
  name            = "${lookup(var.lb_001, "name")}"
  accounting_type = "${var.charge_type}"
  network_volume  = "${lookup(var.lb_001, "network_volume")}" # 10,20,30,40,100,200,300,400,500,600,700,800,900,1000,1100,1200,1300,1400,1500,1600,1700,1800,1900,2000
  #ip_version = "v4" # v4,v6
  #policy_type = "standard" # standard,ats

  listener {
    protocol       = "HTTP" # HTTP,HTTPS,FTP,空
    lb_port        = 80     # 1～65535
    instance_port  = 80     # 1～65535
    balancing_type = "1"    # 1 (Round-Robin) | 2 (Least-Connection)
  }

  health_check {
    unhealthy_threshold = 1        # 1 - 10
    target              = "TCP:80" # TCP:宛先ポート | ICMP
    interval            = "30"     # 5 - 300
  }

  #filter_type         = "1" # 1(許可する) | 2(拒否する)
  #filter_ip_addresses = ["1.1.1.1", "1.1.1.2"]

  #session_stickiness_policy_enable            = true
  #session_stickiness_policy_expiration_period = 60 # 3-60

  # TCP:80 only
  sorry_page_enable      = true
  sorry_page_status_code = 503 # 200 | 503

  # SSL only
  #ssl_certificate_id = ""
  #ssl_policy_id = "1"

  instances = ["${nifcloud_instance.example_server_cent[0].name}", "${nifcloud_instance.example_server_cent[1].name}"]
  #depends_on = [""]
  #lifecycle {
  #  ignore_changes = [""]
  #}
}

resource "nifcloud_lb_port" "example_lb_web_443" {
  name = "${nifcloud_lb.example_lb_web.id}"

  listener {
    protocol       = "HTTPS" # HTTP,HTTPS,FTP,空
    lb_port        = 443     # 1～65535
    instance_port  = 443     # 1～65535
    balancing_type = "1"     # 1 (Round-Robin) | 2 (Least-Connection)
  }

  health_check {
    unhealthy_threshold = 1         # 1 - 10
    target              = "TCP:443" # TCP:宛先ポート | ICMP
    interval            = "30"      # 5 - 300
  }

  #filter_type         = "1" # 1(許可する) | 2(拒否する)
  #filter_ip_addresses = ["1.1.1.1", "1.1.1.2"]

  #session_stickiness_policy_enable            = true
  #session_stickiness_policy_expiration_period = 60 # 3-60

  # TCP:80 only
  #sorry_page_enable      = true
  #sorry_page_status_code = 200 # 200 | 503

  # SSL only
  #ssl_certificate_id = ""
  #ssl_policy_id      = "1"

  instances = ["${nifcloud_instance.example_server_cent[0].name}", "${nifcloud_instance.example_server_cent[1].name}"]
  #depends_on = [""]
  #lifecycle {
  #  ignore_changes = [""]
  #}
}
