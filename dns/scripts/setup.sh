#!/bin/bash
set -euo pipefail

source ./dns/config.env

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
10.8.0.1 dashboard.$VPN_DOMAIN dashboard
EOF

sudo systemctl restart dnsmasq
sudo systemctl enable dnsmasq

sudo iptables -I INPUT 1 -i wg0 -p udp --dport 53 -j ACCEPT
sudo iptables -I INPUT 1 -i wg0 -p tcp --dport 53 -j ACCEPT

sudo apt install iptables-persistent -y
sudo netfilter-persistent save