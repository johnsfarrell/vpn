#!/bin/bash
set -euo pipefail

NAME="${1:?Usage: $0 <device_name>}"
source ./dns/config.env

sudo sed -i "/^.*[[:space:]]$NAME\.$VPN_DOMAIN[[:space:]]$NAME$/d" /etc/dnsmasq.hosts

sudo systemctl restart dnsmasq