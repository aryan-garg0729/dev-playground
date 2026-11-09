# It Keeps Coming Back

Nobody at ElixirBytes is going to hand-recreate pods at 3am. We need something that *keeps* pods running: a **controller**.

A **ReplicaSet** is the simplest one. You give it a pod template, a count and a selector, and it runs the reconciliation loop from chapter 1 forever:

1. Count the pods that match my selector.
2. Too few? Create one from the template. Too many? Delete one.
3. Repeat.

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: brewd-rs
spec:
  replicas: 3
  selector:
    matchLabels:
      app: brewd-rs          # which pods I count...
  template:                  # ...and how I make more
    metadata:
      labels:
        app: brewd-rs        # must match the selector, or the API server rejects it
    spec:
      containers:
        - name: brewd
          image: nginx:1.27-alpine
          ports:
            - containerPort: 80
```

The `template` is a pod manifest without `apiVersion` and `kind`. The pods get generated names like `brewd-rs-x7k2p`.

## Assignment

1. Save the ReplicaSet above as `manifests/brewd-rs.yaml` and apply it.
2. Open a second terminal and run `kubectl get pods -w` (watch mode).
3. In the first terminal, delete one of the `brewd-rs-...` pods. Watch the second terminal: a replacement appears almost before the old one is gone.
4. Try deleting two at once. Same thing.

You can't win this fight by deleting pods. To actually remove them, delete (or scale down) the ReplicaSet.
