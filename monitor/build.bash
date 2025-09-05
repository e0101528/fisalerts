#!/bin/bash
cd src/cmd
go build -o monitor
systemctl stop monitor
cp monitor /usr/local/monitor/monitor
journalctl -umonitor --rotate
journalctl -umonitor --vacuum-time=10s

cd ../../resources/production/
cat config_bare.yaml check/*.check > config.yaml

systemctl start monitor

