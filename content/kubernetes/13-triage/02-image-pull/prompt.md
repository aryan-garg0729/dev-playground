# Pods That Never Start

You come back from lunch to an alert: *pager-api has 0 of 2 pods ready*. Whoever deployed it has gone home. Set up the incident:

```sh
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: Namespace
metadata: { name: triage }
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pager-api
  namespace: triage
spec:
  replicas: 2
  selector:
    matchLabels: { app: pager-api }
  template:
    metadata:
      labels: { app: pager-api }
    spec:
      containers:
        - name: api
          image: nginx:1.27-alpne
          ports: [{ containerPort: 80 }]
YAML
```

Everything in this chapter lives in the **`triage`** namespace, so add `-n triage` to your commands (or switch your context's namespace for the chapter and switch back later).

## Assignment

1. Find out why the pods never start. `kubectl -n triage get pods`, then `describe` one. The Events quote the registry's own words; read them literally.
2. Fix it **in the cluster** so both pods run.
3. Write down what was wrong in **`incident.txt`** (in `~/elixirbytes/k8s`, next to `answers.txt`), in your own words:

   ```
   pager_cause=<what was wrong>
   ```

   The check looks for a key word that shows you found the cause, not for a sentence it wrote. Name the thing that was wrong.
