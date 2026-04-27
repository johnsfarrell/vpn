#!/bin/bash
set -euo pipefail

DEVICE_NAME="${1:?Usage: $0 <device_name>}"
VPN_DOMAIN="${VPN_DOMAIN:-vpn.local}"

sudo sed -i "/^.*[[:space:]]$DEVICE_NAME\.$VPN_DOMAIN$/d" /etc/dnsmasq.hosts

sudo systemctl restart dnsmasq