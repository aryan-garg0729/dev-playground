# Choosing a Node

So far the scheduler has put pods wherever it liked. It works in two phases for every pod:

1. **Filter**: which nodes *can* take it? Enough unreserved CPU and memory for its requests, matching any node rules, no taint it doesn't tolerate.
2. **Score**: of those, which is *best*? It prefers spreading pods out, nodes that already have the image, and so on.

Sometimes you need a say. One of our workers has fast local disks (pretend, it's a container on your laptop), and the search index wants them.

Nodes have labels like everything else. Some come built in:

```sh
kubectl get nodes --show-labels
# kubernetes.io/hostname=elixirbytes-worker, kubernetes.io/os=linux, ...
kubectl label node elixirbytes-worker tier=fast
```

The simplest rule is a **nodeSelector**: "only nodes with these labels".

```yaml
spec:
  template:
    spec:
      nodeSelector:
        tier: fast
```

Node labels and taints belong to the cluster, not to your manifests: `kubectl apply -f manifests/` won't recreate them after a rebuild. Write the commands down.

## Assignment

1. Label **`elixirbytes-worker`** (not `worker2`) with **`tier=fast`**.
2. Create a Deployment **`pinned`**: 2 replicas, label `app: pinned`, `nginx:1.27-alpine`, with a `nodeSelector` of `tier: fast`.
3. Check where its pods went: `kubectl get pods -l app=pinned -o wide`.
