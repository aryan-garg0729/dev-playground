# Requests, Limits and QoS

A node has a fixed amount of CPU and memory. If you don't tell Kubernetes what your containers need, it packs them in blind, and one hungry pod can starve its neighbours.

Each container can declare:

- **requests**: what it needs. The scheduler only places a pod on a node with that much unreserved capacity. It's a reservation, not a cap.
- **limits**: the most it may use. Over the **CPU** limit, it's throttled (slowed down). Over the **memory** limit, it's killed (next lesson).

```yaml
resources:
  requests:
    cpu: 100m        # 100 millicores = 0.1 of a CPU core
    memory: 64Mi     # mebibytes
  limits:
    cpu: 100m
    memory: 64Mi
```

From what you declare, Kubernetes gives each pod a **QoS class**, which decides who gets evicted first when a node runs short of memory:

| QoS class | when |
| --- | --- |
| `Guaranteed` | every container has requests **equal to** limits, for both CPU and memory |
| `Burstable` | at least one request or limit is set, but not Guaranteed |
| `BestEffort` | nothing set at all |

## Assignment

1. Create a Deployment named **`steady`**: 1 replica, label `app: steady`, `nginx:1.27-alpine`, with requests **and** limits both set to **`cpu: 100m`, `memory: 64Mi`**.
2. Look up the QoS class Kubernetes assigned, for a `steady` pod and for a `brewd` pod (which has no resources set yet):

   ```sh
   kubectl get pod <name> -o jsonpath='{.status.qosClass}'
   ```

3. Record both in `answers.txt`:

   ```
   steady_qos=<steady's class>
   brewd_qos=<brewd's class>
   ```

Then predict what brewd's class would become if you gave it only a memory request, and check yourself with the table.
