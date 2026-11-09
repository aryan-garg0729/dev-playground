# Never on the Same Node

Run `kubectl get pods -o wide` and look at brewd: it's quite possible that most of its replicas share one node. "4 replicas" is a lie if one node failing takes out all 4.

**Pod anti-affinity** says "don't put me near pods with these labels", where "near" is a **topology key**: a node label that defines the domain. `kubernetes.io/hostname` means "the same node"; `topology.kubernetes.io/zone` would mean "the same zone".

```yaml
spec:
  template:
    metadata:
      labels: { app: spread }
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchLabels: { app: spread }
              topologyKey: kubernetes.io/hostname
```

`required...` is a hard rule: if it can't be met, the pod stays Pending. `preferredDuringScheduling...` is a soft one, with a `weight`. The `IgnoredDuringExecution` part means pods already running aren't moved if things change later.

## Assignment

Create a Deployment **`spread`**: 2 replicas, label `app: spread`, `nginx:1.27-alpine`, with the required anti-affinity above.

The obvious version gets one pod running and one stuck `Pending`. **Work out why** before you fix it: count the nodes this pod is allowed to use right now (the control plane is tainted, and so is worker2). `required` has no fallback.

The fix is one field you used two lessons ago. Don't remove any taints.
