# Ship It, Then Unship It

brewd 1.28 is ready (played here by `nginx:1.28-alpine`). Change the image and the Deployment rolls it out: new ReplicaSet up, old one down, one pod at a time.

```sh
kubectl set image deploy/brewd brewd=nginx:1.28-alpine   # container=image
kubectl rollout status deploy/brewd     # waits, and prints progress, until it's done
kubectl get rs -l app=brewd             # two ReplicaSets now; the old one at 0
```

`kubectl set image` is imperative again. In real life you'd change the image in `manifests/brewd.yaml` and apply it; the rollout is the same either way.

Every rollout gets a **revision** number, and the Deployment keeps the old ReplicaSets so it can go back:

```sh
kubectl rollout history deploy/brewd
# REVISION  CHANGE-CAUSE
# 1         <none>
# 2         <none>
kubectl rollout undo deploy/brewd                  # back to the previous revision
kubectl rollout undo deploy/brewd --to-revision=1  # or a specific one
```

`rollout undo` is the fastest fix for a bad release: it doesn't rebuild anything, it just scales the old ReplicaSet back up. Then you fix the code properly.

## Assignment

1. Roll `brewd` out to **`nginx:1.28-alpine`** and watch it with `kubectl rollout status`.
2. Dex reports that 1.28 breaks potion checkout. **Roll back**.
3. Read `kubectl rollout history deploy/brewd` and `kubectl get rs -l app=brewd`.

Finish on `nginx:1.27-alpine` with 4 ready pods.

`rollout undo` prints a warning that it doesn't update the `last-applied-configuration` annotation. That's kubectl reminding you the file is the source of truth: re-apply `manifests/brewd.yaml` (it still says 1.27) and confirm `kubectl diff` is empty.
