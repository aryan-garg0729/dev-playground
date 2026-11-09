# Roles and Bindings

A ServiceAccount on its own can do nothing. **RBAC** (role-based access control) grants it things, using two kinds of object:

- A **Role** is a list of rules: which **verbs** (`get`, `list`, `watch`, `create`, `update`, `patch`, `delete`) on which **resources**, in one namespace.
- A **RoleBinding** gives a Role to **subjects**: users, groups or ServiceAccounts.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
rules:
  - apiGroups: [""]              # "" is the core group: pods, services, configmaps, secrets...
    resources: ["pods", "pods/log"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: menu-reader-pods
subjects:
  - kind: ServiceAccount
    name: menu-reader
    namespace: shop
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

`pods/log` is a **subresource**: reading logs is a separate permission from reading the pod object. Deployments live in the `apps` group, so a rule for them says `apiGroups: ["apps"]`; `kubectl api-resources` shows every kind's group.

**ClusterRoles** and **ClusterRoleBindings** are the same idea for the whole cluster. You don't need them here.

You can test permissions without deploying anything, by **impersonating** the ServiceAccount:

```sh
kubectl auth can-i list pods --as=system:serviceaccount:shop:menu-reader -n shop
# yes
kubectl auth can-i --list --as=system:serviceaccount:shop:menu-reader -n shop
```

## Assignment

Save the Role and RoleBinding above in `manifests/rbac.yaml` and apply them. Then try `auth can-i` for `list pods`, `delete pods` and `get secrets` as menu-reader. You want `yes`, `no` and `no`.
