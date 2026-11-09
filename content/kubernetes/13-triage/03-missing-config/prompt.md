# Starts, Dies, Repeats

Next alert: *digest-worker restarting*. Set it up:

```sh
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: ConfigMap
metadata: { name: digest-config, namespace: triage }
data:
  target: "https://elixirbytes.internal/digest"
---
apiVersion: apps/v1
kind: Deployment
metadata: { name: digest-worker, namespace: triage }
spec:
  replicas: 1
  selector:
    matchLabels: { app: digest-worker }
  template:
    metadata:
      labels: { app: digest-worker }
    spec:
      containers:
        - name: worker
          image: busybox:1.36
          command: ["/bin/sh", "-c"]
          args:
            - |
              if [ -z "$DIGEST_TARGET" ]; then
                echo "FATAL: DIGEST_TARGET is not set, refusing to start"
                exit 1
              fi
              echo "digest worker up, posting to $DIGEST_TARGET"
              while true; do sleep 30; done
YAML
```

## Assignment

1. Diagnose it. The container tells you exactly what's wrong right before it dies, if you read the right run's logs.
2. Fix it. There's already a ConfigMap in the namespace holding what it wants. Wire the value in **from the ConfigMap**: don't paste the URL into the Deployment.
3. Record the cause in `incident.txt`:

   ```
   digest_cause=<what was wrong>
   ```
