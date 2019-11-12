#!/bin/bash
yum install -y httpd
systemctl enable httpd
systemctl start httpd

echo "DEVICE=ens192
ONBOOT=yes
BOOTPROTO=static
NETMASK=255.255.255.0
IPADDR=192.168.3.2" > /etc/sysconfig/network-scripts/ifcfg-ens192
echo "192.168.2.0/24 via 192.168.3.250
192.168.201.0/24 via 192.168.3.250" > /etc/sysconfig/network-scripts/route-ens192
ifdown ens192; ifup ens192

## MASQUERADE
#iptables -F INPUT
#iptables -F OUTPUT
#iptables -F FORWARD
#iptables -P INPUT ACCEPT
#iptables -P OUTPUT ACCEPT
#iptables -P FORWARD ACCEPT
#echo 1 > /proc/sys/net/ipv4/ip_forward
#echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
#iptables -t nat -A POSTROUTING -o ens160 -j MASQUERADE
#iptables-save > /etc/sysconfig/iptables
#systemctl enable iptables

## nifcloud ではディスクがあとからアタッチされるので、以下のような処理はuserdata不可
#for i in $(find /sys/class/scsi_host -name 'scan') $(find /sys/devices -name 'scan') ;do echo "- - -" > $i ; done
#echo "n
#p
#1
#
#
#p
#w" | sudo fdisk /dev/sdb
#partprobe
#mkfs.xfs /dev/sdb1
#mkdir /add_disk1
#mount /dev/sdb1 /add_disk1
#ID=`blkid /dev/sdb1 | sed -e "s/^.*UUID=\"\(.*\)\".*TYPE.*$/\1/g"`
#echo "UUID=$ID /add_disk1              xfs     defaults        0 0" >> /etc/fstab
