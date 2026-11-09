# What a Pod Shares

Why does Kubernetes wrap containers in pods at all? Because some containers only make sense together: a web server and a helper that refreshes its files, an app and a proxy that handles its TLS. These helpers are called **sidecars**.

Containers in the same pod:

- are always scheduled onto the **same node**, and start and stop together
- share **one network**: one IP address, one set of ports
- can share **volumes** (folders) if the pod declares them (next lesson)

They still have their own filesystems, images and processes.

```yaml
spec:
  containers:
    - name: brewd
      image: nginx:1.27-alpine
    - name: cache
      image: redis:7-alpine      # listens on port 6379
```

**The `brewd` container in this pod wants to talk to `cache`. Which address works?**
