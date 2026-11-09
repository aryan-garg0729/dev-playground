# Letting the Outside In

A ClusterIP is invisible from outside the cluster. Priya wants the shop reachable from a browser without anyone running `port-forward`.

A **NodePort** Service is a ClusterIP plus one more thing: it opens the same port (30000-32767) on **every node**. Traffic to `<any node>:<nodePort>` goes to the Service, then to a pod.

```yaml
spec:
  type: NodePort
  selector:
    app: brewd
  ports:
    - port: 80
      targetPort: 80
      nodePort: 30080     # leave it out and Kubernetes picks a free one
```

Our kind nodes are containers, so "every node" isn't on your laptop's network. That's what the `extraPortMappings` in `kind.yaml` were for: your laptop's port 30080 maps to the control-plane node's 30080.

```
laptop:30080 -> node:30080 -> Service brewd-np -> one brewd pod:80
```

On a real cluster you'd hit a node's IP directly, and you'd rarely hand NodePorts to customers: five-digit ports and node IPs that change aren't a nice address. They're mostly a building block for the next lesson.

## Assignment

1. Create a Service named **`brewd-np`** of type `NodePort`, selecting `app: brewd`, port 80 to target port 80, with **`nodePort: 30080`**. Save it in `manifests/`.
2. From your laptop, `curl localhost:30080`.

Keep the ClusterIP Service `brewd` too. A pod can be behind any number of Services.
