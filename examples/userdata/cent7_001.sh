#!/bin/bash
yum install -y httpd
systemctl enable httpd
systemctl start httpd

echo "DEVICE=ens192
ONBOOT=yes
BOOTPROTO=static
NETMASK=255.255.255.0
IPADDR=192.168.0.10" > /etc/sysconfig/network-scripts/ifcfg-ens192
ifdown ens192; ifup ens192

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
