# Running Is Not Ready

Last Thursday a brewd pod sat there for forty minutes with a wedged worker. `kubectl get pods` said `Running` the whole time, and the Service kept sending it a quarter of all orders.

`Running` means the container's process exists. It says nothing about whether it works. Kubernetes can't know what "works" means for your app, so you tell it, with **probes**: checks the kubelet runs against each container, over and over.

A **readiness probe** answers *"should this pod get traffic right now?"*:

```yaml
containers:
  - name: brewd
    readinessProbe:
      httpGet:              # also: tcpSocket, exec (run a command), grpc
        path: /
        port: 80
      initialDelaySeconds: 2   # wait this long after start before the first check
      periodSeconds: 5         # then check every 5 seconds
      failureThreshold: 3      # 3 failures in a row = not ready
```

While a pod is not ready:

- it shows `0/1` in the READY column,
- it's **removed from the Service's endpoints**, so it gets no traffic,
- it is **not** restarted. Readiness only controls traffic.

It also makes rollouts honest. With `maxUnavailable: 0`, the Deployment waits for each new pod to be *ready* before stopping an old one. Without a probe, "ready" just means "started", so a broken release counts as ready the moment its process launches.

## Assignment

Add the readiness probe above to brewd in `manifests/brewd.yaml` (path `/`, port 80, `periodSeconds: 5`), apply it, and watch the rollout with `kubectl get pods -w`: each new pod spends a moment at `0/1 Running` before going `1/1`.
