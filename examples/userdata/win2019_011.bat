@REM init script
@echo off
netsh interface ipv4 set add name="Ethernet0" source=static addr="192.168.0.111" mask="255.255.255.0" gateway="192.168.0.10"
netsh interface ipv4 set dns name="Ethernet0" source=static addr="8.8.8.8" register=non validate=no
netsh interface ipv4 add dns name="Ethernet0" addr="8.8.4.4" index=2 validate=no