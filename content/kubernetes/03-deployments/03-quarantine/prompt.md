# Quarantine a Pod

Dex's favourite trick. One brewd pod is acting strangely: slow, odd log lines. If you delete it, the evidence goes with it. If you leave it, customers keep hitting it.

You now know the ReplicaSet only counts pods that **match its selector**. So change the pod's label:

```sh
kubectl label pod brewd-rs-x7k2p app=quarantine --overwrite
```

What happens, in order:

1. The pod no longer matches `app=brewd-rs`. The ReplicaSet lets go of it (it removes itself from the pod's `ownerReferences`).
2. The ReplicaSet now counts 2 and creates a fresh replacement.
3. The odd pod keeps running, out of rotation, for you to `exec` into and `logs` at your leisure.

Kubernetes tracks ownership in `metadata.ownerReferences`. A pod a controller owns lists it there:

```sh
kubectl get pod brewd-rs-x7k2p -o jsonpath='{.metadata.ownerReferences[*].kind}/{.metadata.ownerReferences[*].name}'
# ReplicaSet/brewd-rs
```

Services (chapter 4) pick pods by label too, so the same relabel takes the pod out of traffic.

## Assignment

1. Pick one `brewd-rs-...` pod and relabel it to **`app=quarantine`**.
2. Confirm the ReplicaSet is back to 3 ready pods, and that your quarantined pod is still `Running` alongside them, with no owner.

Quarantine **exactly one** pod. When you're done investigating a real one, you'd delete it by hand: no controller will.
