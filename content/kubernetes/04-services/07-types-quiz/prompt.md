# Picking a Service Type

The three types stack on top of each other:

| type | reachable from | what you get |
| --- | --- | --- |
| `ClusterIP` | inside the cluster | a virtual IP and a DNS name |
| `NodePort` | `<node>:30000-32767` | a ClusterIP, plus a port on every node |
| `LoadBalancer` | an external address | a NodePort, plus a cloud load balancer in front |

The ElixirBytes **ledger** service records every sale. Only the payments worker, a pod in the same cluster, ever talks to it. It holds money data.

**Which Service type should ledger get?**
