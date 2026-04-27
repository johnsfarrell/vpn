#!/bin/bash
set -euo pipefail

DEVICE_NAME="${1:?Usage: $0 <device_name>}"
CLIENT_DIR="./wireguard/clients/$DEVICE_NAME"
PUBLIC_KEY="$(tr -d '\n' < "$CLIENT_DIR/client_public.key")"

# remove this peer from saved config
sudo env PUBLIC_KEY="$PUBLIC_KEY" perl -0pi -e '
  my $key = quotemeta($ENV{"PUBLIC_KEY"});
  s/\n?\[Peer\]\nPublicKey = $key\nAllowedIPs = [^\n]+\n//g;
' /etc/wireguard/wg0.conf

# fast replacement for: sudo systemctl restart wg-quick@wg0
sudo wg set wg0 peer "$PUBLIC_KEY" remove

sudo bash ./dns/scripts/client_delete.sh "$DEVICE_NAME"

rm -rf "$CLIENT_DIR"
