# Who Is This Pod?

The status page wants to show which ElixirBytes services are up, so it needs to **list pods** through the Kubernetes API. Who is it, as far as the API server is concerned?

Every request to the API server is **authenticated** (who are you?) and then **authorised** (are you allowed to do this?). You authenticate with the certificate kind put in your kubeconfig. Pods authenticate as a **ServiceAccount**: an identity that lives in a namespace and exists for programs, not people.

- Every namespace has one called `default`, and a pod that doesn't name one runs as it.
- Kubernetes mounts a short-lived token for the pod's ServiceAccount at `/var/run/secrets/kubernetes.io/serviceaccount/token`.
- In RBAC rules and audit logs, a ServiceAccount's username is `system:serviceaccount:<namespace>:<name>`.

```sh
kubectl create serviceaccount menu-reader
kubectl get serviceaccounts
kubectl get pod storefront -o jsonpath='{.spec.serviceAccountName}'
```

Giving every app its own ServiceAccount means you can grant each one exactly what it needs, and see in the logs who did what. Sharing `default` means granting one app something grants it to all of them.

## Assignment

1. Create a ServiceAccount named **`menu-reader`** in `shop`.
2. Add two lines to `answers.txt`:

   ```
   storefront_sa=<the ServiceAccount the storefront pod runs as>
   menu_reader_user=<menu-reader's full username, as RBAC sees it>
   ```
