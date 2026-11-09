# Claiming Storage

A PersistentVolumeClaim says what you need, not where it comes from:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: orders
spec:
  accessModes:
    - ReadWriteOnce          # mountable read-write by pods on ONE node at a time
  resources:
    requests:
      storage: 1Gi
  # storageClassName: standard   # leave it out to use the default class
```

```sh
kubectl get pvc
# NAME     STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS
# orders   Pending                                      standard
kubectl describe pvc orders
```

## Assignment

1. Create the PVC **`orders`** above in `manifests/orders-pvc.yaml` and apply it. Don't create anything that uses it yet.
2. Look at its status, then read `kubectl describe pvc orders` to find out why it's in that state. You learned the reason in the last lesson.
3. Add one line to `answers.txt`:

   ```
   orders_phase=<the claim's STATUS right now>
   ```

Run `codequest submit` **before** the next lesson: once a pod uses the claim, its phase changes.

A `Pending` PVC isn't always this harmless. The same status, with an event saying the StorageClass doesn't exist, is one of the triage cases in chapter 13.
