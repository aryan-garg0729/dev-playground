# The Menu API Front Door

Priya's ticket:

> **menu-api: internal by default, public on our one door.** The menu API gets 3 pods (`nginx:1.27-alpine` until the real image lands). Other services call it on port **8080**, which is our convention for internal APIs, even though the container listens on 80. Staging needs to reach it too.
>
> Our kind cluster has exactly one public port, 30080, and brewd-np is squatting on it. brewd doesn't need to be public any more (customers will come in through Ingress next sprint), so give the port to the menu API.

## Assignment

1. A Deployment **`menu-api`**: 3 replicas, label/selector **`app: menu-api`**, `nginx:1.27-alpine` on port 80.
2. A ClusterIP Service **`menu-api`** on **port 8080**, forwarding to the pods' port 80.
3. Delete the **`brewd-np`** Service (from your cluster and from `manifests/`).
4. A NodePort Service **`menu-public`** for the `app: menu-api` pods, port 80 to target port 80, **`nodePort: 30080`**.

Check your work from the `visitor` pod in `shop-staging` and from your laptop with `curl`, the way the boss will. Keep `brewd` and `brewd-lb` as they are.

Every check is a hit on the boss. Defeat it to unlock **Ingress**.
