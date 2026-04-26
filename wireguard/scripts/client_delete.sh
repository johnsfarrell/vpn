#!/bin/bash
set -euo pipefail

DEVICE_NAME="${1:?Usage: $0 <device_name>}"
CLIENT_DIR="./wireguard/clients/$DEVICE_NAME"
PUBLIC_KEY="$(cat "$CLIENT_DIR/client_public.key")"

# delete peer's section from wg0 config
sudo perl -0pi -e "s/\n?\[Peer\]\nPublicKey = \Q$PUBLIC_KEY\E\nAllowedIPs = [^\n]+\n//" /etc/wireguard/wg0.conf

sudo systemctl restart wg-quick@wg0

rm -rf "$CLIENT_DIR"

