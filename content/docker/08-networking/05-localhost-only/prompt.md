# For Your Eyes Only

`-p 18089:8080` publishes on **every** network interface of your machine: `0.0.0.0` (all IPv4 addresses) and `[::]` (all IPv6 ones). On a laptop in a café, that's everyone on the Wi-Fi. On a server, it can be the whole internet.

Docker also edits the host's firewall rules to make published ports work, so a "deny all" firewall rule might not protect a published port the way you'd expect.

Put an IP address in front to publish on one interface only:

```sh
docker run -d -p 127.0.0.1:18089:8080 ...   # loopback: reachable only from this machine
docker port <container>
# 8080/tcp -> 127.0.0.1:18089
```

## Assignment

Dex wants an admin copy of brewd for poking at orders, and it must never be reachable from the network.

1. Run **`eb-admin`** from `elixirbytes/brewd:1.6.0`, detached, published on **`127.0.0.1:18089`** only.
2. Submit from `~/elixirbytes`. The checks read `docker port` (no `0.0.0.0` or `[::]` allowed) and call `/healthz` on `127.0.0.1:18089`.
