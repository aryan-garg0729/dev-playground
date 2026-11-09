# No More Zombie Pods

After this week, Dex has written down rules for every service that takes money. Priya's ticket applies them to checkout:

> **checkout: production-grade from day one.**
> - 3 replicas of `nginx:1.27-alpine` (stand-in), label `app: checkout`, port 80, behind a ClusterIP Service `checkout` on port 80.
> - **Readiness**: HTTP `GET /` on port 80. Pods that aren't ready get no traffic.
> - **Liveness**: also HTTP `GET /` on port 80 (a cheap local check, nothing that touches a database), starting **10 seconds** after the container starts.
> - **Resources**: request `50m` CPU and `32Mi` memory. Cap memory at `64Mi`. **No CPU limit**: we'd rather checkout run hot than be throttled during the sale.
> - It must come up clean: zero restarts.

## Assignment

Write `manifests/checkout.yaml` with the Deployment and the Service, and apply it.

Before you apply, predict the QoS class from the requests and limits above. Then check yourself.

If a pod restarts while you're experimenting, delete the pod (the Deployment replaces it) so the restart counts start from zero again.

Every check is a hit on the boss. Defeat it to unlock **Jobs, CronJobs and DaemonSets**.
