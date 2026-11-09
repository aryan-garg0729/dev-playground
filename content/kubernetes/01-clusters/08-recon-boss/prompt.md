# Day-One Recon Report

Priya won't approve a migration to a cluster nobody has looked at. Her ticket:

> **Recon before we move anything.** I need four facts about the cluster in a file I can paste into the migration doc, and a second namespace so we can rehearse the move before doing it for real. Read the facts off the cluster. Don't copy them from a tutorial.

## Assignment

In `~/elixirbytes/k8s`, create **`recon.txt`** with four `key=value` lines:

```
control_plane=<the name of the control-plane node>
k8s_version=<that node's Kubernetes version, as kubectl shows it, e.g. v1.31.2>
workers=<how many nodes are NOT control-plane nodes>
dns_namespace=<the namespace the cluster's DNS pods (coredns) run in>
```

Then create a namespace called **`shop-staging`**, next to `shop`.

Hints you've already earned: `get nodes` shows roles and versions, `-A` looks in every namespace, and `-o wide` or `describe` fill in the rest.

Every check is a hit on the boss. The checks read each fact live from the cluster and compare it with your file. Defeat it to unlock **Pods**.
