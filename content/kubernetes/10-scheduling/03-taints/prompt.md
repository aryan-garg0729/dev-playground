# A Node in Maintenance

`elixirbytes-worker2` is due for maintenance. No new work should land there, but the maintenance team's own diagnostics pod has to.

A **taint** marks a node so pods stay away unless they explicitly accept it. It has a key, an optional value, and an **effect**:

| effect | means |
| --- | --- |
| `NoSchedule` | new pods without a matching toleration aren't scheduled here; running pods stay |
| `PreferNoSchedule` | try to avoid it, but it's allowed |
| `NoExecute` | also **evicts** running pods that don't tolerate it |

```sh
kubectl taint node elixirbytes-worker2 maintenance=true:NoSchedule
kubectl taint node elixirbytes-worker2 maintenance-       # remove it (trailing minus)
```

A pod accepts a taint with a toleration:

```yaml
tolerations:
  - key: maintenance
    operator: Equal
    value: "true"
    effect: NoSchedule
```

You saw one already: the DaemonSet in chapter 8 tolerated the control plane's taint.

## Assignment

1. Taint **`elixirbytes-worker2`** with **`maintenance=true:NoSchedule`**.
2. Create a Deployment **`tolerant`** (1 replica, label `app: tolerant`, `busybox:1.36` running `sleep infinity`) whose pod runs **on `elixirbytes-worker2`**, in spite of the taint.

A toleration alone isn't enough: it lets the pod onto worker2, but doesn't send it there. You also need something that pulls it there. Every node has a `kubernetes.io/hostname` label.

Leave the taint in place; the next lessons use it.
