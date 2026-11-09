# Policies Add Up

`vault-proxy` pods are protected by your `vault-proxy-ingress` policy. Then the monitoring team adds their own:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-monitoring
spec:
  podSelector:
    matchLabels: { app: vault-proxy }
  policyTypes: [Ingress]
  ingress:
    - from:
        - namespaceSelector:
            matchLabels: { kubernetes.io/metadata.name: monitoring }
```

**Who can now reach the vault-proxy pods?**
