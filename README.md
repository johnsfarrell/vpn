# VPN

## VM instance

We chose to use [Oracle Free Tier](https://www.oracle.com/cloud/free/) for cloud computing needs. The goal is to maintain a free VPN service.

Launch a VM instance:

| Property | Value |
|------ |------ |
| Shape | VM.Standard.E2.1.Micro |
| Image | Canonical-Ubuntu-22.04-2026.03.31-0 |
| Public IPv4 | Yes |

Add security rule:

| Stateless | Source | IP Protocol | Source Port Range | Destination Port Range | Type and Code | Allows |
| --- | --- | --- | --- | --- | --- | --- |
| No | 0.0.0.0/0 | UDP | All | 51820 | | UDP traffic for ports: 51820 |


## VM setup

One time setup script:

```bash
sudo bash ./wireguard/scripts/setup.sh
```

Adds WireGuard private and public keys to `/etc/wireguard/server_private.key` and `/etc/wireguard/server_public.key`, respectively.

The VPS IP address is `10.8.0.1/24`.

Creating clients:

```bash
./wireguard/scripts/client_create.sh <device_name>
```

## Dashboard

One time setup script:

```bash
sudo bash ./dashboard/scripts/setup.sh
```

Run dashboard:

```bash
go run ./dashboard/*.go
```

## WireGuard for clients

Install [WireGuard app](https://www.wireguard.com/install/) on your device.

## DNS

One time setup script:

```bash
# export VPN_DOMAIN=custom.domain.internal
sudo bash ./dns/setup.sh
```

If `VPN_DOMAIN` is not set, scripts default to `vpn.local`.

