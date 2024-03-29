#!/bin/bash
yum install -y httpd mod_ssl
systemctl enable httpd
systemctl start httpd

echo "DEVICE=ens160
ONBOOT=yes
BOOTPROTO=static
GATEWAY=192.168.3.1
NETMASK=255.255.255.0
IPADDR=192.168.3.101" > /etc/sysconfig/network-scripts/ifcfg-ens160
echo "192.168.2.0/24 via 192.168.3.250
192.168.201.0/24 via 192.168.3.250" > /etc/sysconfig/network-scripts/route-ens160
ifdown ens160; ifup ens160

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
