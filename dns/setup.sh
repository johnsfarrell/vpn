#!/bin/bash
set -euo pipefail

VPN_DOMAIN="${VPN_DOMAIN:-vpn.local}"

sudo apt install -y dnsmasq dnsutils

sudo tee /etc/dnsmasq.d/vpn.conf > /dev/null <<EOF
interface=wg0
listen-address=10.8.0.1
bind-interfaces

local=/$VPN_DOMAIN/

addn-hosts=/etc/dnsmasq.hosts

server=1.1.1.1
server=8.8.8.8
EOF

sudo tee /etc/dnsmasq.hosts > /dev/null <<EOF
10.8.0.1 $VPN_DOMAIN
EOF

sudo systemctl restart dnsmasq
sudo systemctl enable dnsmasq