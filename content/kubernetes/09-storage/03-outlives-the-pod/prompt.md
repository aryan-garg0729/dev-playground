# Data That Outlives the Pod

A pod uses a claim like any other volume, with `persistentVolumeClaim` instead of `emptyDir`:

```yaml
spec:
  volumes:
    - name: data
      persistentVolumeClaim:
        claimName: orders
  containers:
    - name: writer
      volumeMounts:
        - name: data
          mountPath: /data
```

The moment a pod using `orders` is scheduled, the provisioner creates a volume on that node, binds it to the claim, and the claim goes `Bound`. From then on, whichever pod mounts `orders` gets the same data.

## Assignment

1. Create a Deployment named **`ledger-writer`**: 1 replica, label `app: ledger-writer`, `busybox:1.36` running `sleep infinity`, mounting the `orders` claim at **`/data`**.
2. Once it's running, write a line from inside the pod that records which pod wrote it:

   ```sh
   kubectl exec deploy/ledger-writer -- sh -c 'echo "written by $(hostname)" > /data/orders.log'
   ```

3. **Delete the pod** (not the Deployment). Wait for the replacement.
4. Read the file from the new pod: `kubectl exec deploy/ledger-writer -- cat /data/orders.log`. The hostname in it is a pod that no longer exists.

A pod's hostname is its name, so the check can tell the writer and the current pod apart.
