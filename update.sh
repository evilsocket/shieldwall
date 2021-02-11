#!/usr/bin/env bash
service shieldwall-agent stop || true
cp shieldwall-agent /usr/bin/
cp shieldwall-agent.service /etc/systemd/system/
systemctl daemon-reload
service shieldwall-agent start