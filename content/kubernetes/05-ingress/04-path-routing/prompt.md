# Routing by Path

An Ingress lists rules. Each rule has paths, and each path points at a Service and port:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: shop
spec:
  ingressClassName: nginx        # which controller handles this Ingress
  rules:
    - http:
        paths:
          - path: /potions
            pathType: Prefix     # /potions, /potions/, /potions/elixir...
            backend:
              service:
                name: potions
                port:
                  number: 80
```

`pathType: Prefix` matches whole path segments: `/potions` matches `/potions/elixir` but not `/potionsandmore`. `Exact` matches only the exact path.

Apply that, `curl localhost/potions`, and you'll get a **404**, but a different one from the last lesson. Read it carefully: it isn't from the controller this time. The request reached the potions pod, and the pod was asked for **`/potions`**. Its httpd only has `/index.html`.

The controller passes the original path through. To strip it, ingress-nginx supports an **annotation** on the Ingress:

```yaml
metadata:
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
```

Now `/potions` (and `/potions/anything`) arrives at the pod as `/`. Annotations like this are controller-specific: another controller would ignore this one, which is one of the things the Gateway API was designed to fix.

## Assignment

1. Write **`manifests/ingress.yaml`**: an Ingress named **`shop`**, class `nginx`, with two `Prefix` paths: **`/potions`** to Service `potions` and **`/scrolls`** to Service `scrolls`, both on port 80.
2. Apply it, see the 404, and work out where it came from (`kubectl logs deploy/potions` helps).
3. Add the rewrite annotation and apply again.

`curl localhost/potions` and `curl localhost/scrolls` should now show each shop's page.
