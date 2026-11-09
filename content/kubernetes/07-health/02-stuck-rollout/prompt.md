# The Rollout That Stopped Itself

The next brewd release adds a proper health endpoint, `/healthz`, and points the readiness probe at it. Except the endpoint didn't make it into the build. (Our stand-in nginx has no `/healthz` either, so it returns 404, and the probe fails.)

Watch what the probe and `maxUnavailable: 0` do together:

```sh
# in manifests/brewd.yaml, change the readiness probe path to /healthz, then:
kubectl apply -f manifests/brewd.yaml
kubectl rollout status deploy/brewd --timeout=60s
# Waiting for deployment "brewd" rollout to finish: 1 out of 4 new replicas have been updated...
# error: timed out waiting for the condition
kubectl get pods -l app=brewd
kubectl describe pod <the new one>     # Events: Readiness probe failed: HTTP probe failed with statuscode: 404
```

The Deployment created **one** new pod (`maxSurge: 1`), waited for it to become ready, and it never did. So it never stopped a single old pod. The shop kept running on the old version the whole time, and customers noticed nothing.

A Deployment that makes no progress for `progressDeadlineSeconds` (default 600) gets marked as failed in its status (`ProgressDeadlineExceeded`), but it doesn't roll itself back. You do that.

## Assignment

1. Change the readiness probe path to **`/healthz`**, apply, and watch the rollout stall.
2. While it's stuck, count and record in `answers.txt`:

   ```
   stuck_new_ready=<how many pods of the NEW ReplicaSet are ready>
   stuck_old_ready=<how many pods of the OLD ReplicaSet are ready>
   ```

   (`kubectl get rs -l app=brewd` shows both.)
3. Roll back with `kubectl rollout undo deploy/brewd`, and put the probe path back to `/` in your file.
