#!/bin/bash
set -euo pipefail

DOMAIN="${1:?Usage: $0 <domain_name>}"

sudo sed -i "/^address=\/$DOMAIN\/0.0.0.0/d" /etc/dnsmasq.d/block.conf

sudo systemctl restart dnsmasq