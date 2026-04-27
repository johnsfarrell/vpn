#!/bin/bash
set -euo pipefail

sudo apt install golang-go -y

# redirect HTTP on 10.8.0.1 to port 8080 (dashboard)
sudo iptables -t nat -D PREROUTING -i wg0 -p tcp --dport 80 -j REDIRECT --to-ports 8080
sudo iptables -t nat -I PREROUTING 1 -i wg0 -d 10.8.0.1 -p tcp --dport 80 -j REDIRECT --to-ports 8080
sudo iptables -I INPUT 1 -i wg0 -p tcp --dport 80 -j ACCEPT
sudo iptables -I INPUT 1 -i wg0 -p tcp --dport 8080 -j ACCEPT

sudo apt install iptables-persistent -y
sudo netfilter-persistent save