# Stuck Before It Ever Runs

*archive has not started in 20 minutes.* No crashes, no restarts. It just never begins.

```sh
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: PersistentVolumeClaim
metadata: { name: archive-data, namespace: triage }
spec:
  accessModes: ["ReadWriteOnce"]
  storageClassName: fast-ssd
  resources:
    requests:
      storage: 1Gi
---
apiVersion: apps/v1
kind: Deployment
metadata: { name: archive, namespace: triage }
spec:
  replicas: 1
  selector:
    matchLabels: { app: archive }
  template:
    metadata:
      labels: { app: archive }
    spec:
      containers:
        - name: archive
          image: busybox:1.36
          command: ["/bin/sh", "-c", "echo archive up; while true; do sleep 30; done"]
          volumeMounts:
            - { name: data, mountPath: /var/archive }
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: archive-data
YAML
```

## Assignment

1. The pod is `Pending`, and the scheduler's message mentions an *unbound* claim: it won't place a pod whose volume doesn't exist yet. The claim, in turn, is waiting on something that doesn't exist. Follow the chain: pod → claim → class.
2. Fix it so the claim binds and archive runs, **keeping the claim's name** `archive-data` (the Deployment refers to it).

   Try `kubectl edit pvc archive-data` first and see what the API server says: most of a PVC's spec, including `storageClassName`, is **immutable** after creation. You'll have to do something more drastic to the claim. That's safe here only because this claim never held any data. Say so in your incident note; on a real claim with data, you'd stop and think.

3. Record:

   ```
   archive_cause=<what was wrong>
   ```
