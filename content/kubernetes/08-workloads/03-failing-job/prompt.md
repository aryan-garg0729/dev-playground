# A Job That Fails Properly

The supplier import keeps failing at 2am, and by the time anyone looks, the evidence is gone. With `restartPolicy: Never`, every failed attempt stays behind as a pod, logs and all.

## Assignment

Create a Job named **`bad-batch`** (`busybox:1.36`) with **`backoffLimit: 2`** and **`restartPolicy: Never`**, whose container does this:

```yaml
command: ["sh", "-c", "echo 'ERROR: supplier feed unreachable'; exit 1"]
```

Watch it fail. Retries are spaced out (10 s, then 20 s), so give it about a minute:

```sh
kubectl get job bad-batch -w
kubectl get pods -l job-name=bad-batch
kubectl describe job bad-batch      # Events, and the Failed condition at the end
```

When the Job has given up, add one line to `answers.txt`:

```
bad_batch_pods=<how many pods it created in total>
```

Predict the number before you count.
