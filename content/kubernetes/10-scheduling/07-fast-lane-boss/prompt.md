# Fast-Lane Placement

The maintenance team is done with worker2, and it came back with SSDs. Priya's ticket:

> **Use the new disks, and stop putting all our eggs on one node.**
> - Maintenance is over: take the maintenance taint off `elixirbytes-worker2`, and label it `disk=ssd`.
> - **vault-proxy**: 2 replicas (`nginx:1.27-alpine`, label `app: vault-proxy`) that may **only** run on SSD nodes. Use **node affinity** with the `In` operator, not a nodeSelector, so we can add other disk types to the list later.
> - **catalog**: 2 replicas (`nginx:1.27-alpine`, label `app: catalog`) that must **never** share a node.
> - And Dex says that `stubborn` pod has been Pending for days. Don't touch it; if you do the above right, it'll sort itself out.

## Assignment

Make it so. Node affinity looks like this (the rest is from this chapter):

```yaml
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
            - key: <label key>
              operator: In
              values: [<value>, ...]
```

Every check is a hit on the boss. Defeat it to unlock **Access control**.
