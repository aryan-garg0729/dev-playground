# Hard Rules Have No Fallback

A cluster has **3** nodes that all accept the `search` pods (no taints, enough resources). The search Deployment is scaled to **5** replicas, with:

```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchLabels: { app: search }
        topologyKey: kubernetes.io/hostname
```

**What happens?**
