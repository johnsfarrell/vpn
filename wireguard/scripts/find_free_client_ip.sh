#!/bin/bash
set -euo pipefail

for i in $(seq 2 254); do
  candidate="10.8.0.$i/32"

  if ! sudo grep -q "AllowedIPs = $candidate" "/etc/wireguard/wg0.conf"; then
    echo "$candidate"
    exit 0
  fi
done

echo "No free IP found for client"
exit 1