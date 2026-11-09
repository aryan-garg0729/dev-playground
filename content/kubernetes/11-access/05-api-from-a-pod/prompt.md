# Calling the API From a Pod

`auth can-i` asks the API server a hypothetical question. Now let's make a real request, the way the status page will.

Inside every pod (unless you turn it off with `automountServiceAccountToken: false`):

- `/var/run/secrets/kubernetes.io/serviceaccount/token`: a token for the pod's ServiceAccount
- `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt`: the cluster's certificate
- `https://kubernetes.default.svc`: the API server, as an ordinary Service

The API is plain HTTPS and JSON. `kubectl get pods -n shop` is really `GET /api/v1/namespaces/shop/pods`.

## Assignment

1. Run a pod named **`api-caller`** as the **`menu-reader`** ServiceAccount, using an image that has `curl`:

   ```yaml
   spec:
     serviceAccountName: menu-reader
     containers:
       - name: curl
         image: curlimages/curl:8.10.1
         command: ["sleep", "infinity"]
   ```

2. From inside it, list pods with the mounted token:

   ```sh
   kubectl exec -it api-caller -- sh
   TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
   curl -sk -H "Authorization: Bearer $TOKEN" https://kubernetes.default.svc/api/v1/namespaces/shop/pods | head -20
   ```

   (`-k` skips certificate checks to keep this short. Real code uses `--cacert` with the mounted `ca.crt`.)

3. Now try a delete, as a **dry run** so nothing can actually be deleted even if RBAC were wrong:

   ```sh
   curl -sk -X DELETE -H "Authorization: Bearer $TOKEN" \
     "https://kubernetes.default.svc/api/v1/namespaces/shop/pods/storefront?dryRun=All"
   ```

   Read the 403's `message`: it names the user, the verb and the resource. That one line is usually the whole RBAC debugging session.

4. Add the username from that message to `answers.txt`:

   ```
   denied_user=<the user the 403 message names>
   ```
