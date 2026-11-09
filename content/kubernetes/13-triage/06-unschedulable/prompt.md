# The Pod Nobody Can Fit

*report-builder: 0 of 1 ready.*

```sh
kubectl apply -f - <<'YAML'
apiVersion: apps/v1
kind: Deployment
metadata: { name: report-builder, namespace: triage }
spec:
  replicas: 1
  selector:
    matchLabels: { app: report-builder }
  template:
    metadata:
      labels: { app: report-builder }
    spec:
      containers:
        - name: builder
          image: busybox:1.36
          command: ["/bin/sh", "-c", "echo building; while true; do sleep 30; done"]
          resources:
            requests:
              memory: "64Gi"
              cpu: "16"
YAML
```

## Assignment

1. `describe` gives you the scheduler's own arithmetic. Compare what the pod asks for with what a node actually has: `kubectl describe node elixirbytes-worker`, under **Allocatable** and **Allocated resources**.
2. Fix it with requests that are realistic for a report builder (a few hundred millicores and a few hundred MiB is plenty) and fit on a worker. **Keep requests set**: deleting them would make it schedule, but as BestEffort, first in line for eviction (chapter 7).
3. Record:

   ```
   report_cause=<what was wrong>
   ```
