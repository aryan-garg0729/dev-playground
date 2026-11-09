# An Ingress Is Just Rules

Mo copied an Ingress from a blog post into a brand-new test cluster:

```sh
kubectl apply -f blog-ingress.yaml
# ingress.networking.k8s.io/blog created
kubectl get ingress
# NAME   CLASS   HOSTS   ADDRESS   PORTS   AGE
# blog   nginx   *                 80      10m
curl http://<the cluster>/blog
# curl: (7) Failed to connect
```

The `ADDRESS` column is empty. On your cluster it fills in with the controller's address a few seconds after you create an Ingress.

**What's most likely wrong?**
