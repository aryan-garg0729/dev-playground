# The 3am Pager

It's 03:04. Your phone, not Dex's: you're on call now.

> **night-orders is down.** Late-night orders are failing. Last deploy: 22:15, by someone who is asleep.

Set up the incident:

```sh
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: ConfigMap
metadata: { name: night-orders-config, namespace: triage }
data:
  queue_url: "amqp://queue.elixirbytes.internal/orders"
---
apiVersion: apps/v1
kind: Deployment
metadata: { name: night-orders, namespace: triage }
spec:
  replicas: 2
  selector:
    matchLabels: { app: night-orders }
  template:
    metadata:
      labels: { app: night-orders }
    spec:
      containers:
        - name: web
          image: nginx:1.27-alpine
          ports: [{ containerPort: 80 }]
          env:
            - name: QUEUE_URL
              valueFrom:
                configMapKeyRef: { name: night-orders-config, key: queue-url }
          readinessProbe:
            httpGet: { path: /, port: 8080 }
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata: { name: night-orders, namespace: triage }
spec:
  selector: { app: night-orders }
  ports: [{ port: 80, targetPort: 8080 }]
YAML
```

## Assignment

There's more than one thing wrong, and each fault hides the next: you won't see the second until you've fixed the first. Work the way this chapter taught you. Look at the pod's status, then `describe`, then the Service's endpoints, then traffic.

Fix everything **in the cluster** until a customer (the storefront pod) gets the page from `http://night-orders.triage`. Keep `QUEUE_URL` coming from the ConfigMap. Then write the post-incident notes in `incident.txt`, one line per fault:

```
config_cause=<what was wrong with the configuration>
probe_cause=<what was wrong with the readiness check>
service_cause=<what was wrong with the Service>
```

Every check is a hit on the final boss. Defeat it to finish the course.
