#!/bin/bash
set -euo pipefail

NAME="${1:?Usage: $0 <device_name> <device_address>}"
ADDR="${2:?Usage: $0 <device_name> <device_address>}"
source ./dns/config.env

sudo tee -a /etc/dnsmasq.hosts > /dev/null <<EOF
$ADDR $NAME.$VPN_DOMAIN $NAME
EOF

sudo systemctl restart dnsmasq
sudo systemctl enable dnsmasq