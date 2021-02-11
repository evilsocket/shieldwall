#!/usr/bin/env bash
service shieldwall-agent stop || true
cp shieldwall-agent /usr/bin/
mkdir -p /etc/shieldwall/
test -s /etc/shieldwall/config.yaml || cp agent.example.yaml /etc/shieldwall/config.yaml
cp shieldwall-agent.service /etc/systemd/system/
systemctl daemon-reload
service shieldwall-agent start