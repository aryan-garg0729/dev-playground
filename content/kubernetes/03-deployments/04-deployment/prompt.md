# Deployments

ReplicaSets keep the count. Now Priya asks for brewd 1.28. A ReplicaSet can't help: change its template and it does... nothing. Existing pods match the selector, so the count is fine and they keep running the old image. You'd have to delete pods by hand and hope.

A **Deployment** manages ReplicaSets for you. Change its pod template and it:

1. creates a **new** ReplicaSet for the new template,
2. scales the new one up and the old one down, a few pods at a time,
3. keeps the old ReplicaSet (at 0) so you can go back.

The manifest looks almost the same as a ReplicaSet, plus a **strategy**:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: brewd
spec:
  replicas: 3
  selector:
    matchLabels:
      app: brewd
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1          # at most 1 pod ABOVE replicas during an update
      maxUnavailable: 0    # never fewer than replicas ready
  template:
    ...
```

`maxSurge: 1, maxUnavailable: 0` means "start one new pod, wait until it's ready, then stop one old pod, repeat". Slower, but the shop never has fewer than 3 pods serving. (The defaults are 25% each.) The other strategy, `Recreate`, stops everything first; it's for apps that can't run two versions at once.

From now on you'll almost never write a bare ReplicaSet. You write Deployments, and they write ReplicaSets.

## Assignment

1. Write **`manifests/brewd.yaml`**: a Deployment named **`brewd`**, 3 replicas, selector and pod label **`app: brewd`**, one container `brewd` running **`nginx:1.27-alpine`** on port 80, and the `maxSurge: 1` / `maxUnavailable: 0` strategy above.
2. Apply it, then look at what it built: `kubectl get deploy,rs,pods -l app=brewd`.
3. Delete the old `brewd-rs` ReplicaSet and the quarantined pod from the last lesson; we're done with them.
