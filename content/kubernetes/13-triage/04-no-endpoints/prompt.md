# Healthy Pods, Empty Service

*status-board is down*, says the alert. The pods disagree: they're all `1/1 Running`.

```sh
kubectl apply -f - <<'YAML'
apiVersion: apps/v1
kind: Deployment
metadata: { name: status-board, namespace: triage }
spec:
  replicas: 2
  selector:
    matchLabels: { app: status-board }
  template:
    metadata:
      labels: { app: status-board }
    spec:
      containers:
        - name: web
          image: nginx:1.27-alpine
          ports: [{ containerPort: 80 }]
---
apiVersion: v1
kind: Service
metadata: { name: status-board, namespace: triage }
spec:
  selector: { app: statusboard }
  ports: [{ port: 80, targetPort: 80 }]
YAML
```

## Assignment

The pods are fine, so look at what stands **between** the Service and its pods. One command tells you the Service is broken; comparing two more tells you why.

Fix it on the **Service** side, check it from the storefront pod (`wget -qO- http://status-board.triage`), and record:

```
board_cause=<what was wrong>
```
