#!/bin/bash
set -euo pipefail

DOMAIN="${1:?Usage: $0 <domain_name>}"

sudo tee -a /etc/dnsmasq.d/block.conf > /dev/null <<EOF
address=/$DOMAIN/0.0.0.0
EOF

sudo systemctl restart dnsmasq