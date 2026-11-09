# Install an Ingress Controller

ElixirBytes now has brewd, the menu API, and two new mini-shops coming this week. One NodePort or LoadBalancer each means five-digit ports, or a cloud load balancer (and bill) per service. Priya wants **one front door** that routes by URL: `/potions` to one service, `/scrolls` to another.

That front door has two parts:

- An **Ingress** is an object holding routing rules: "requests for this host and path go to that Service."
- An **ingress controller** is a real program (a pod) that reads every Ingress and actually proxies the traffic. Without one, Ingress objects are just notes nobody reads.

Kubernetes doesn't ship a controller. You install one. We'll use **ingress-nginx**, which runs nginx and rewrites its config whenever an Ingress changes. kind publishes a build of it set up for kind: it runs on the node labelled `ingress-ready=true` (the control plane, from your `kind.yaml`) and listens on that node's ports 80 and 443, which you mapped to your laptop.

One honest note: the ingress-nginx project was **retired in March 2026**. It no longer gets new features, but it's still everywhere in production, and the Ingress API it implements is stable and isn't going away. Its successor, the Gateway API, comes up at the end of this chapter.

## Assignment

1. Install the controller and wait for it:

   ```sh
   kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.15.1/deploy/static/provider/kind/deploy.yaml
   kubectl -n ingress-nginx wait --for=condition=ready pod -l app.kubernetes.io/component=controller --timeout=180s
   ```

2. Look at what it created: `kubectl get all -n ingress-nginx` and `kubectl get ingressclass`. An **IngressClass** names a controller, so an Ingress can say which one should handle it.
3. `curl -i localhost`. A `404 Not Found` from nginx is success: you reached the controller, and it has no rules yet.
