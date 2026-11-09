# Finding brewd by Name

Nobody types `10.96.41.170`. Every cluster runs a DNS server (CoreDNS, the pods you found in `kube-system` in the chapter 1 boss), and every Service gets a name:

```
<service>.<namespace>.svc.cluster.local
```

Pods get a DNS search list, so shorter forms work too:

| from a pod in... | name that works |
| --- | --- |
| the same namespace (`shop`) | `brewd` |
| another namespace | `brewd.shop` (or the full name) |
| anywhere, always | `brewd.shop.svc.cluster.local` |

You can only test this from **inside** the cluster: your laptop doesn't use the cluster's DNS. Use a pod you already have, or start a throwaway one that deletes itself on exit:

```sh
kubectl exec storefront -- wget -qO- http://brewd
kubectl run tmp --rm -it --image=busybox:1.36 --restart=Never -- sh
# / # nslookup brewd
# / # wget -qO- http://brewd
```

busybox's `nslookup` tries every entry in the search list, so expect a few `NXDOMAIN` lines (and exit code 1) around the answer. The line you want starts with `Name:`.

## Assignment

1. From `storefront`, fetch `http://brewd` and look up `brewd` with `nslookup`.
2. Add the fully qualified name `nslookup` gives you to `answers.txt`:

   ```
   brewd_fqdn=<the full name of the brewd Service>
   ```

3. Create a long-running pod called **`visitor`** in the **`shop-staging`** namespace (`busybox:1.36`, running `sleep 86400`). From inside it, check that `http://brewd` **fails** and then find the name that **works**.
