# Pods With Names and Their Own Disks

The secrets vault ElixirBytes is evaluating runs as a small cluster of its own: every member keeps its own data, and the members find each other **by name**. A Deployment can't do that. Its pods are interchangeable: random names, all sharing one claim (if any), started and stopped in any order.

A **StatefulSet** gives each pod a stable identity:

- **Names by ordinal**: `vault-0`, `vault-1`... A replacement for `vault-0` is called `vault-0` again.
- **Its own storage**: `volumeClaimTemplates` creates one PVC per pod (`data-vault-0`, `data-vault-1`), and a pod always gets *its* claim back.
- **Order**: pods start one at a time, `0` first, and are removed in reverse.
- **Its own DNS name**, through a **headless Service** (next lesson).

A headless Service is a Service with `clusterIP: None`: no virtual IP and no load balancing. Instead, DNS returns the pods themselves.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: vault
spec:
  clusterIP: None
  selector: { app: vault }
  ports: [{ port: 80 }]
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: vault
spec:
  serviceName: vault            # the headless Service that gives pods their DNS names
  replicas: 2
  selector:
    matchLabels: { app: vault }
  template:
    metadata:
      labels: { app: vault }
    spec:
      containers:
        - name: vault
          image: busybox:1.36
          command: ["sh", "-c", "[ -f /data/id ] || echo \"I am $(hostname)\" > /data/id; exec sleep infinity"]
          volumeMounts:
            - name: data
              mountPath: /data
  volumeClaimTemplates:          # one PVC per pod, named <template>-<pod>
    - metadata:
        name: data
      spec:
        accessModes: [ReadWriteOnce]
        resources:
          requests:
            storage: 1Gi
```

## Assignment

Save both objects above as `manifests/vault.yaml` and apply. In a second terminal, run `kubectl get pods -l app=vault -w` and watch the order: `vault-1` doesn't start until `vault-0` is ready. Then `kubectl get pvc`.
