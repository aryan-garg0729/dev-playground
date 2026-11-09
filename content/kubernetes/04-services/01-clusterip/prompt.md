# One Name for Many Pods

brewd runs as 4 pods, each with its own IP. Those IPs changed during every rollout last chapter, and they'll change again. The storefront can't keep a list of pod IPs. It needs **one stable name** that always reaches *some* healthy brewd.

That's a **Service**. It has a fixed virtual IP and a DNS name, and it forwards traffic to whatever pods currently match its selector:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: brewd
spec:
  type: ClusterIP        # the default: reachable only inside the cluster
  selector:
    app: brewd           # send traffic to pods with this label
  ports:
    - port: 80           # the port clients connect to on the Service
      targetPort: 80     # the port the container listens on
```

`port` and `targetPort` don't have to match. A Service can listen on 8080 and forward to the pods' 80.

The Service keeps an up-to-date list of the pod IPs it currently sends to, in objects called **EndpointSlices**:

```sh
kubectl get svc brewd
# NAME    TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)
# brewd   ClusterIP   10.96.41.170   <none>        80/TCP
kubectl get endpointslices -l kubernetes.io/service-name=brewd
# NAME          ADDRESSTYPE   PORTS   ENDPOINTS
# brewd-k8x2q   IPv4          80      10.244.1.7,10.244.2.4,10.244.1.8 + 1 more...
```

Older material uses `kubectl get endpoints brewd`. The Endpoints API still works but is deprecated; EndpointSlices replaced it.

## Assignment

1. Write **`manifests/brewd-svc.yaml`** with the Service above and apply it.
2. Compare the EndpointSlice's addresses with `kubectl get pods -l app=brewd -o wide`. Every ready pod's IP should be there.
3. Delete one brewd pod and look again: the list updates by itself.
