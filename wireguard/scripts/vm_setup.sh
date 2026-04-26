#!/bin/bash
set -euo pipefail

sudo apt update && sudo apt upgrade -y

sudo apt install wireguard qrencode -y

# enable IP forwarding

echo "net.ipv4.ip_forward=1" | sudo tee /etc/sysctl.d/vpn.conf

sudo sysctl --system

# setup WireGuard configuration

wg genkey | sudo tee /etc/wireguard/server_private.key | wg pubkey | sudo tee /etc/wireguard/server_public.key

sudo bash -c 'cat > /etc/wireguard/wg0.conf <<EOF
[Interface]
Address = 10.8.0.1/24
ListenPort = 51820
PrivateKey = $(cat /etc/wireguard/server_private.key)

PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -A FORWARD -o wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o ens3 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -D FORWARD -o wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o ens3 -j MASQUERADE
EOF'

sudo systemctl enable wg-quick@wg0
sudo systemctl start wg-quick@wg0
sudo wg

### UPDATE IP PACKET FILTER RULES

# allow input traffic on port 51820 and persist the rules
sudo iptables -I INPUT 1 -p udp --dport 51820 -j ACCEPT

# remove rule for rejecting forwarded traffic
sudo iptables -D FORWARD -j REJECT --reject-with icmp-host-prohibited 2>/dev/null || true

# allow traffic incoming/outgoing from wireguard
sudo iptables -I FORWARD 1 -i wg0 -j ACCEPT
sudo iptables -I FORWARD 1 -o wg0 -j ACCEPT

# rewrite source IP to look like came from VPN
sudo iptables -t nat -C POSTROUTING -o ens3 -j MASQUERADE 2>/dev/null || sudo iptables -t nat -A POSTROUTING -o ens3 -j MASQUERADE

sudo apt install iptables-persistent -y
sudo netfilter-persistent save