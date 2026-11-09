# Base64 Is Not Encryption

Try this on your own cluster:

```sh
kubectl get secret brewd-db -o jsonpath='{.data.password}'
# NjFjZDNhOTk4YjE1...                <- looks scrambled
kubectl get secret brewd-db -o jsonpath='{.data.password}' | base64 -d
# 61cd3a998b15...                     <- your password
```

That's the whole of a default Secret's "security": base64 exists so binary data fits in YAML, not to hide anything. On a default cluster, Secrets are also stored unencrypted in etcd.

So Secrets aren't magic, but they aren't pointless either. Their value is that they're a **separate kind** that the cluster can treat specially.

**Which two things actually protect a Secret? Select both.**
