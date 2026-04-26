#!/bin/bash
set -euo pipefail

DEVICE_NAME="${1:?Usage: $0 <device_name>}"

CLIENT_DIR="./wireguard/clients/$DEVICE_NAME"
CLIENT_CONF="$CLIENT_DIR/client.conf"
CLIENT_PUBLIC_KEY_FILE="$CLIENT_DIR/client_public.key"
CLIENT_PRIVATE_KEY_FILE="$CLIENT_DIR/client_private.key"

CLIENT_PUBLIC_KEY="$(cat "$CLIENT_PUBLIC_KEY_FILE")"
CLIENT_PRIVATE_KEY="$(cat "$CLIENT_PRIVATE_KEY_FILE")"
CLIENT_IP="$(awk -F ' = ' '/^Address = /{print $2}' "$CLIENT_CONF" | head -n1)"
DNS="$(awk -F ' = ' '/^DNS = /{print $2}' "$CLIENT_CONF" | head -n1)"
ENDPOINT="$(awk -F ' = ' '/^Endpoint = /{print $2}' "$CLIENT_CONF" | head -n1)"

LAST_HANDSHAKE_UNIX="0"
TRANSFER_RX_BYTES="0"
TRANSFER_TX_BYTES="0"

WG_DUMP="$(sudo wg show wg0 dump 2>/dev/null)"
PEER_LINE="$(echo "$WG_DUMP" | awk -F'\t' -v pub="$CLIENT_PUBLIC_KEY" '$1 == pub {print $0; exit}')"
LAST_HANDSHAKE_UNIX="$(echo "$PEER_LINE" | awk -F'\t' '{print $5}')"
TRANSFER_RX_BYTES="$(echo "$PEER_LINE" | awk -F'\t' '{print $6}')"
TRANSFER_TX_BYTES="$(echo "$PEER_LINE" | awk -F'\t' '{print $7}')"

json_escape() {
  local s="$1"
  s="${s//\\/\\\\}"
  s="${s//\"/\\\"}"
  s="${s//$'\n'/\\n}"
  s="${s//$'\r'/\\r}"
  s="${s//$'\t'/\\t}"
  printf '%s' "$s"
}

printf '{'
printf '"name":"%s",' "$(json_escape "$DEVICE_NAME")"
printf '"directory":"%s",' "$(json_escape "$CLIENT_DIR")"
printf '"ip":"%s",' "$(json_escape "$CLIENT_IP")"
printf '"publicKey":"%s",' "$(json_escape "$CLIENT_PUBLIC_KEY")"
printf '"privateKey":"%s",' "$(json_escape "$CLIENT_PRIVATE_KEY")"
printf '"dns":"%s",' "$(json_escape "${DNS:-unknown}")"
printf '"endpoint":"%s",' "$(json_escape "${ENDPOINT:-unknown}")"
printf '"lastHandshakeUnix":%s,' "$LAST_HANDSHAKE_UNIX"
printf '"transferRxBytes":%s,' "$TRANSFER_RX_BYTES"
printf '"transferTxBytes":%s,' "$TRANSFER_TX_BYTES"
printf '}\n'
exit 0