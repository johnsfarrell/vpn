#!/bin/bash
set -euo pipefail

DEVICE_NAME="${1:?Usage: $0 <device_name>}"
CLIENT_DIR="./wireguard/clients/$DEVICE_NAME"
PUBLIC_IP_ADDRESS=$(curl -fsS ifconfig.me)
CLIENT_IP="$("./wireguard/scripts/find_free_client_ip.sh")"


mkdir -p "$CLIENT_DIR"

wg genkey | tee "$CLIENT_DIR/client_private.key" | wg pubkey > "$CLIENT_DIR/client_public.key"

sudo bash -c 'cat >> /etc/wireguard/wg0.conf <<EOF

[Peer]
PublicKey = '"$(cat "$CLIENT_DIR/client_public.key")"'
AllowedIPs = '"$CLIENT_IP"'
EOF'

sudo systemctl restart wg-quick@wg0

cat > "$CLIENT_DIR/client.conf" <<EOF
[Interface]
PrivateKey = $(cat "$CLIENT_DIR/client_private.key")
Address = $CLIENT_IP
DNS = 1.1.1.1

[Peer]
PublicKey = $(sudo cat /etc/wireguard/server_public.key)
Endpoint = $PUBLIC_IP_ADDRESS:51820
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
EOF

qrencode -t ansiutf8 < "$CLIENT_DIR/client.conf"
