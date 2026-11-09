# Where Service Traffic Goes

Mo, poking around in the storefront pod:

```sh
kubectl exec storefront -- wget -qO- -T 3 http://brewd    # works
kubectl exec storefront -- ping -c 2 -W 2 brewd
# PING brewd (10.96.41.170): 56 data bytes
# --- brewd ping statistics ---
# 2 packets transmitted, 0 packets received, 100% packet loss
```

*"Half-broken Service?"*

Some background. Two different things handle cluster networking:

- The **CNI plugin** (kindnet here; Calico or Cilium elsewhere) gives every pod a real, routable IP, across nodes.
- **kube-proxy** runs on every node and implements Services. When a Service or its endpoints change, it rewrites that node's packet-forwarding rules: "a packet for `10.96.41.170:80/TCP` goes to one of these pod IPs instead."

**Why does `wget` work while `ping` gets nothing back?**
