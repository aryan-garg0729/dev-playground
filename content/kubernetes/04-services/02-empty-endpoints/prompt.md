# The Service That Serves Nothing

Dex, at the end of his shift: *"potion-api is up, both pods are green, and every request to it times out. I'm going home."*

This is the most common Service bug there is, so let's reproduce it. Apply this exactly as written:

```sh
kubectl apply -f - <<'YAML'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: potion-api
spec:
  replicas: 2
  selector:
    matchLabels: { app: potion-api }
  template:
    metadata:
      labels: { app: potion-api }
    spec:
      containers:
        - name: api
          image: nginx:1.27-alpine
          ports: [{ containerPort: 80 }]
---
apiVersion: v1
kind: Service
metadata:
  name: potion-api
spec:
  selector: { app: potion_api }
  ports: [{ port: 80, targetPort: 80 }]
YAML
```

(`-f -` reads the manifest from stdin. `---` separates two objects in one file.)

The pods are fine. The Service is fine *as an object*: the API server happily accepts a selector that matches nothing. Here's what that looks like:

```sh
kubectl get endpointslices -l kubernetes.io/service-name=potion-api
# NAME               ADDRESSTYPE   PORTS     ENDPOINTS
# potion-api-9zq4d   IPv4          <unset>   <unset>
```

**No endpoints means the selector matches no ready pods.** Burn that pattern in; you'll use it in the triage chapter. To find out why, compare the Service's selector with the pods' labels:

```sh
kubectl get svc potion-api -o jsonpath='{.spec.selector}'
kubectl get pods --show-labels
```

## Assignment

Find the mismatch and fix it **on the Service**. Other things already depend on the pods' labels, so don't touch the Deployment.

Both ready pods should then show up as the Service's endpoints.
