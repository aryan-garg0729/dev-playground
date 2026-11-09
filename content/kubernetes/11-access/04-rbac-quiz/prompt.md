# There Is No Deny

During a review, Dex runs:

```sh
kubectl auth can-i delete pods --as=system:serviceaccount:shop:menu-reader -n shop
# yes
```

`pod-reader` doesn't grant `delete`. He digs further and finds that someone also created this, months ago, "just to get the status page working":

```yaml
kind: RoleBinding
metadata: { name: status-page-quickfix }
subjects: [{ kind: ServiceAccount, name: menu-reader, namespace: shop }]
roleRef: { kind: ClusterRole, name: edit, apiGroup: rbac.authorization.k8s.io }
```

(`edit` is a built-in ClusterRole that can change almost everything in a namespace.)

**How do you take `delete pods` away from menu-reader?**
