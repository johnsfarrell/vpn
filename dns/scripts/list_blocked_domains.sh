#!/bin/bash
set -euo pipefail

sudo touch /etc/dnsmasq.d/block.conf

sudo grep -o "^address=/[^/]\+/0.0.0.0$" /etc/dnsmasq.d/block.conf | cut -d'/' -f2 || true