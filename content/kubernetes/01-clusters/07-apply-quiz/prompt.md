# Apply, Create or Replace

You'll describe almost everything in YAML files called **manifests**, kept in `~/elixirbytes/k8s/manifests/`. Three kubectl verbs send a manifest to the cluster:

| command | what it does |
| --- | --- |
| `kubectl create -f x.yaml` | create it; **fail if it already exists** |
| `kubectl replace -f x.yaml` | overwrite it; **fail if it doesn't exist**, and drop fields you left out |
| `kubectl apply -f x.yaml` | make the cluster match the file: create or update |

Two more you'll use constantly:

```sh
kubectl diff -f x.yaml     # what WOULD change? run this before applying to anything you care about
kubectl apply -f manifests/  # a whole folder at once
```

Writing down *what you want* and letting the cluster work out *how* is called **declarative** configuration. Typing one-off commands that do things (`kubectl run`, `kubectl scale`) is **imperative**. Imperative is fine for experiments; declarative is what you can review, commit to Git and replay after a disaster.

**Priya wants one command in the deploy script that works on the very first deploy *and* every deploy after it. Which one?**
