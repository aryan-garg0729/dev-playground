# Scaling a StatefulSet Down

The vault was scaled to 3 replicas for a load test: `vault-0`, `vault-1` and `vault-2`, each with its own claim. The test is over:

```sh
kubectl scale sts vault --replicas=2
```

**What happens to the pods and their storage?**
