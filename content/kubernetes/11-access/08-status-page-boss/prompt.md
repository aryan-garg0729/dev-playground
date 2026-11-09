# Least-Privilege Status Page

Priya's ticket, after the `status-page-quickfix` incident:

> **status-page, done properly.**
> - Its own ServiceAccount, `status-page`. Not `default`, not `menu-reader`.
> - It may **read** pods, Deployments and Services in `shop`, and nothing else: no changes, no Secrets, no other namespaces.
> - It runs as a Deployment, `status-page` (1 replica, `nginx:1.27-alpine`, label `app: status-page`).
> - Network: only the ingress controller (namespace `ingress-nginx`) may connect to it. Not the storefront, not anyone else.

## Assignment

Build it: a ServiceAccount, a Role, a RoleBinding, the Deployment, and a NetworkPolicy named **`status-page-ingress`**.

You haven't written a rule for another API group yet. Deployments are in `apps`, so they need their own rule (or their own entry in `apiGroups`). `kubectl api-resources` lists every kind's group.

To allow a whole namespace in a NetworkPolicy, use a `namespaceSelector`. Every namespace carries a built-in label `kubernetes.io/metadata.name: <its name>`.

Test with `auth can-i --as=system:serviceaccount:shop:status-page` before the boss does.

Every check is a hit on the boss. Defeat it to unlock **Autoscaling and CRDs**.
