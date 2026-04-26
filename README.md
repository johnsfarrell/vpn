# VPN

## Creating a VM instance

We chose to use [Oracle Free Tier](https://www.oracle.com/cloud/free/) for cloud computing needs. The goal is to maintain a free VPN service.

Launch a VM instance:

| Property | Value |
|------ |------ |
| Shape | VM.Standard.E2.1.Micro |
| Image | Canonical-Ubuntu-22.04-2026.03.31-0 |
| Public IPv4 | Yes |

## VM setup

One time setup script:

```bash
sudo bash ./wireguard/scripts/vm_setup.sh
```

Adds WireGuard private and public keys to `/etc/wireguard/server_private.key` and `/etc/wireguard/server_public.key`, respectively.

Creating clients:

```bash
./wireguard/scripts/client_setup.sh <device_name>
```
