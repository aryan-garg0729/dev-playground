# Reading a Pending Pod

A pod the scheduler can't place stays `Pending` forever, with no restarts, no logs, and no error in `kubectl get`. The explanation is in `describe`, written by the scheduler itself:

```sh
kubectl describe pod <name>
# Events:
#   Warning  FailedScheduling  default-scheduler
#   0/3 nodes are available: 1 Insufficient memory, 2 node(s) had untolerated taint(s). ...
```

Read it as arithmetic: **every node is accounted for**, each with the first filter that ruled it out. The numbers add up to the node count. "Insufficient memory" means the pod's requests don't fit; "untolerated taint" means a taint it doesn't tolerate; "didn't match Pod's node affinity/selector" means a nodeSelector or affinity rule.

The tail, `preemption: ... Preemption is not helpful`, means the scheduler also checked whether evicting lower-priority pods would help, and it wouldn't.

## Assignment

Create a pod named **`stubborn`** (`busybox:1.36`, `sleep infinity`) with a `nodeSelector` of **`kubernetes.io/hostname: elixirbytes-worker2`** and **no** tolerations. It will sit `Pending`.

Read its `describe` output, and record the scheduler's arithmetic in `answers.txt`:

```
stubborn_tainted=<how many nodes were ruled out by a taint>
stubborn_mismatch=<how many were ruled out by the node selector>
```

Then work out *which* nodes those are. Leave `stubborn` Pending: it comes back in the boss.
