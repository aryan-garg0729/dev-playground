# One on Every Node

Some things have to run on **every** machine: a log shipper reading each node's log files, a metrics agent, a network plugin. With a Deployment you'd set `replicas: 3`, hope they spread out, and then add a node and forget to bump the count.

A **DaemonSet** runs exactly one pod per node. There's no `replicas` field: the node count *is* the replica count. Add a node and a pod appears on it; remove one and its pod goes too.

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-scout
spec:
  selector:
    matchLabels: { app: node-scout }
  template:
    metadata:
      labels: { app: node-scout }
    spec:
      containers:
        - name: scout
          image: busybox:1.36
          env:
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName     # the "downward API": facts about the pod itself
          command: ["sh", "-c", "echo scouting $NODE_NAME; exec sleep infinity"]
```

Apply that and count the pods: 2, not 3. The control plane has a **taint**, a mark that says "don't schedule here unless you explicitly accept it":

```sh
kubectl describe node elixirbytes-control-plane | grep Taints
# Taints: node-role.kubernetes.io/control-plane:NoSchedule
```

A pod accepts a taint with a **toleration** (more on both in the Scheduling chapter):

```yaml
    spec:
      tolerations:
        - key: node-role.kubernetes.io/control-plane
          operator: Exists
          effect: NoSchedule
```

## Assignment

Create the DaemonSet **`node-scout`** above, see it skip the control plane, then add the toleration so it runs on **all three** nodes. Check with `kubectl get pods -l app=node-scout -o wide`.
