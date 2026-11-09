# Read the Nodes

Before Dex lets anyone deploy to a cluster, he asks two questions: *how many machines can take work right now, and what's actually running the containers?* Guessing doesn't count. The cluster will tell you.

`kubectl get` lists objects of a kind. Nodes are objects like any other:

```sh
kubectl get nodes
# NAME                        STATUS   ROLES           AGE   VERSION
# elixirbytes-control-plane   Ready    control-plane   5m    v1.3x.x
# ...
```

Add `-o wide` for more columns, including `CONTAINER-RUNTIME`:

```sh
kubectl get nodes -o wide
```

Every object is really a big structured document. `-o yaml` shows all of it, and `-o jsonpath` pulls out one field. Paths start at the object's root, and `.items[0]` is the first item of a list:

```sh
kubectl get nodes -o jsonpath='{.items[0].metadata.name}'
# elixirbytes-control-plane
kubectl get node elixirbytes-worker -o jsonpath='{.status.nodeInfo.kubeletVersion}'
# v1.3x.x
```

`kubectl describe node <name>` gives the human-friendly version, with the node's recent **Events** at the bottom. You'll live in `describe` for the rest of this course.

## Assignment

Our migration runbook needs the facts, written down. In `~/elixirbytes/k8s`, create **`answers.txt`** with two `key=value` lines:

```
nodes=<how many nodes are Ready>
runtime=<the container runtime and version, exactly as Kubernetes reports it>
```

1. Count the nodes whose `STATUS` is `Ready`.
2. Find the runtime string. It looks like `containerd://2.1.3`: scheme, `://`, version. Copy it exactly, including the `containerd://` part.

The checks read the real values from the cluster and compare them with your file, so the file must match *your* cluster, not this page.

You'll keep adding lines to `answers.txt` all course. One `key=value` per line, no spaces around `=`.
