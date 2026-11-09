# Restarting a Wedged Container

Readiness takes a sick pod out of traffic and waits for it to recover. Some things never recover by themselves: a deadlock, a worker thread stuck forever. The process is alive, so restartPolicy never kicks in. A **liveness probe** covers this: if it fails `failureThreshold` times in a row, the kubelet kills the container and restarts it.

`exec` probes run a command inside the container. Exit code 0 is healthy, anything else is a failure:

```yaml
livenessProbe:
  exec:
    command: ["cat", "/tmp/healthy"]
  initialDelaySeconds: 5
  periodSeconds: 5
  failureThreshold: 3
```

## Assignment

Create a pod named **`zombie`** (`busybox:1.36`) that is healthy for 20 seconds and then wedges: it deletes its health file but keeps running.

```yaml
command: ["sh", "-c", "touch /tmp/healthy; sleep 20; rm /tmp/healthy; sleep 3600"]
```

Give it the liveness probe above (`cat /tmp/healthy`, every 5 seconds). Then watch:

```sh
kubectl get pod zombie -w            # RESTARTS goes up about once a minute
kubectl describe pod zombie          # Events: Liveness probe failed ... will be restarted
```

Add one line to `answers.txt`:

```
zombie_exit=<the exit code of zombie's last terminated run>
```

Then work out what it means: the process never exited by itself, so who stopped it, and with which signal? (Exit codes above 128 are 128 + the signal number.)

Why about a minute, when the probe fails after about 35 seconds? The kubelet first sends SIGTERM and waits up to `terminationGracePeriodSeconds` (30 by default) for a clean exit. A shell running as PID 1 ignores SIGTERM, so the kubelet waits the full 30 seconds and then sends SIGKILL. It's the PID 1 lesson from the Docker course again, this time in Kubernetes.
