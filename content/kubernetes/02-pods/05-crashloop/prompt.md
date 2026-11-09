# The Pod That Won't Stay Up

Some containers don't crash once. They crash, get restarted, crash again, forever. You'll see this more than any other failure in your Kubernetes life, so let's make one on purpose.

When a container exits, the kubelet looks at the pod's **`restartPolicy`**:

| restartPolicy | restarts the container when it... |
| --- | --- |
| `Always` | exits for any reason (the default) |
| `OnFailure` | exits with a non-zero code |
| `Never` | never |

A container that keeps dying gets restarted with a growing delay: 10 s, 20 s, 40 s... up to five minutes. While it waits, the pod shows **`CrashLoopBackOff`**. That's not the error; it's Kubernetes saying "this keeps failing, so I'm backing off". The error is in the logs.

```sh
kubectl get pods
# NAME       READY   STATUS             RESTARTS     AGE
# cauldron   0/1     CrashLoopBackOff   3 (20s ago)  80s
kubectl logs cauldron              # the current run (often empty: it just started)
kubectl logs cauldron --previous   # the run that crashed: read THIS one
kubectl describe pod cauldron      # Last State: Terminated, Reason, Exit Code
```

## Assignment

Save this as **`manifests/cauldron.yaml`** and apply it. It picks a random batch number and exit code the first time it starts (and keeps them in an `emptyDir`, which survives container restarts), then crashes with them every time. Your numbers won't match anyone else's.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cauldron
  labels:
    app: cauldron
spec:
  volumes:
    - name: state
      emptyDir: {}
  containers:
    - name: cauldron
      image: busybox:1.36
      volumeMounts:
        - name: state
          mountPath: /state
      command:
        - sh
        - -c
        - |
          [ -f /state/n ] || od -An -N2 -tu1 /dev/urandom > /state/n
          set -- $(cat /state/n)
          sleep 3
          echo "cauldron cracked at batch $(( $1 % 90 + 10 ))"
          exit $(( $2 % 50 + 3 ))
```

Notice there's no `restartPolicy` in it. Let it crash a couple of times, then investigate and add three lines to `answers.txt`:

```
cauldron_exit=<the exit code it crashes with>
cauldron_policy=<the restartPolicy the pod actually has>
cauldron_said=<the full line it prints before crashing>
```

If you delete and recreate the pod, it picks new numbers, so check `answers.txt` again.
