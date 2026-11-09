# A Room of Our Own

A cluster is shared. The ingress controller, the DNS server and the cluster's own plumbing already run in it, and soon other teams will too. If two teams both create something called `web`, whose is it?

A **namespace** is a named room inside the cluster. Object names only have to be unique within a namespace, so ten teams can each have a `web`. Namespaces are also where access rules (chapter 11) and resource quotas attach.

```sh
kubectl get namespaces
# NAME              STATUS   AGE
# default           Active   10m
# kube-system       Active   10m    <- the cluster's own components
# ...
kubectl get pods -n kube-system     # -n picks the namespace
kubectl get pods -A                 # -A means every namespace
```

If you don't say `-n`, kubectl uses your **context**'s namespace. A context is the triple kubectl uses for every command: *which cluster, as which user, in which namespace*. kind created one called `kind-elixirbytes`:

```sh
kubectl config get-contexts
kubectl config use-context kind-elixirbytes             # switch cluster
kubectl config set-context --current --namespace=shop   # change the default namespace
```

Not everything lives in a namespace. Nodes, for one, belong to the whole cluster. Deleting a namespace deletes **everything** in it, with no undo.

## Assignment

Everything ElixirBytes builds lives in one namespace, **`shop`**.

1. Create the namespace `shop`.
2. Make sure `kind-elixirbytes` is your current context.
3. Set `shop` as that context's default namespace, so you stop typing `-n shop` on every command for the rest of the course.
4. Check it: `kubectl get pods` should now say `No resources found in shop namespace.`
