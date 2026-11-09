# Who Gets Evicted First

A worker node is running out of memory. These three pods are on it:

| pod | resources | QoS class |
| --- | --- | --- |
| `steady` | requests = limits = 64Mi | Guaranteed |
| `brewd-...` | requests 32Mi, limit 128Mi, using 40Mi | Burstable |
| `debug-shell` | nothing set | BestEffort |

The kubelet must free memory now, so it starts **evicting** pods: stopping them so their controllers can recreate them elsewhere.

**Which pod is evicted first?**
