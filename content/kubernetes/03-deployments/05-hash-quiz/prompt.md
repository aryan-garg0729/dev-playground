# Reading a Pod Name

Run this and look closely at the names:

```sh
kubectl get rs,pods -l app=brewd --show-labels
# NAME                               DESIRED   CURRENT   READY
# replicaset.apps/brewd-6d4cf56db6   3         3         3
#
# NAME                         READY   STATUS    LABELS
# pod/brewd-6d4cf56db6-x2kq9   1/1     Running   app=brewd,pod-template-hash=6d4cf56db6
# ...
```

The pod's name has three parts: `brewd`, `6d4cf56db6` and `x2kq9`. The middle part also shows up as the `pod-template-hash` label, and the ReplicaSet's selector includes it, so two ReplicaSets of the same Deployment never fight over pods.

**What is `6d4cf56db6`?**
