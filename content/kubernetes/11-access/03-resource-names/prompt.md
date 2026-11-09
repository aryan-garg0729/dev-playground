# Just That One ConfigMap

The status page should also show the current shop banner, which lives in the `brewd-config` ConfigMap. Mo suggests adding `configmaps` to the Role. That grants **every** ConfigMap in the namespace, and someone will eventually put something sensitive in one.

A rule can name the exact objects it covers with **`resourceNames`**:

```yaml
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["brewd-config"]
    verbs: ["get"]
```

One catch: `resourceNames` works for verbs that act on one named object (`get`, `update`, `delete`...) but not for `list` or `watch`, which ask for "all of them". That's fine here: the status page knows the name it wants.

Roles are **additive**: a subject's permissions are the union of every rule in every Role bound to it. There's no "deny" rule, so you restrict by not granting.

## Assignment

Add a rule to the **`pod-reader`** Role so menu-reader can **`get`** the ConfigMap **`brewd-config`**, and nothing else about ConfigMaps. Keep the pod rules as they are.

Check with `auth can-i` as menu-reader: `get configmap/brewd-config` should say `yes`, and `get configmap/kube-root-ca.crt` (another ConfigMap every namespace has) should say `no`.
