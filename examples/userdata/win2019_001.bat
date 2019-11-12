@REM init script
@echo off
netsh interface ipv4 set add name="Ethernet1" source=static addr="192.168.3.51" mask="255.255.255.0" gateway=""
netsh interface ipv4 set dns name="Ethernet0" source=static addr="8.8.8.8" register=non validate=no
netsh interface ipv4 add dns name="Ethernet0" addr="8.8.4.4" index=2 validate=no
route -p add 192.168.2.0 mask 255.255.255.0 192.168.3.250
route -p add 192.168.201.0 mask 255.255.255.0 192.168.3.250
