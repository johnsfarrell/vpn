#!/bin/bash
set -euo pipefail

sudo apt install golang-go -y

# redirect HTTP to port 8080 (dashboard)
sudo iptables -t nat -I PREROUTING 1 -i wg0 -p tcp --dport 80 -j REDIRECT --to-ports 8080
sudo iptables -I INPUT 1 -i wg0 -p tcp --dport 80 -j ACCEPT
sudo iptables -I INPUT 1 -i wg0 -p tcp --dport 8080 -j ACCEPT

sudo apt install iptables-persistent -y
sudo netfilter-persistent save