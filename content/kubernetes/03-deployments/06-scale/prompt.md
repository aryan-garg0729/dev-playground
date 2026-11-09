# Scaling, and Keeping the File Honest

Black Moon is Friday. Mo: *"Can we just have more brewds?"*

Yes. Scaling only changes `replicas`, not the pod template, so the Deployment keeps the same ReplicaSet and just changes its count:

```sh
kubectl scale deploy brewd --replicas=6
kubectl get rs -l app=brewd      # same ReplicaSet, DESIRED 6
```

That was imperative. The cluster now says 6, your file still says 3, and the next person to run `kubectl apply -f manifests/` quietly scales the shop back to 3 in the middle of the sale. This drift between file and cluster is how real outages start.

`kubectl diff` shows the gap before it bites:

```sh
kubectl diff -f manifests/brewd.yaml
# -  replicas: 6
# +  replicas: 3
```

`diff` exits `0` when they agree and `1` when they differ, so scripts can check it too.

## Assignment

1. Scale `brewd` to 6 with `kubectl scale` and watch the new pods arrive.
2. Run `kubectl diff -f manifests/brewd.yaml` and read the drift.
3. Priya settles on **4** replicas for the sale. Make that true **in the file** (`manifests/brewd.yaml`), apply it, and finish with `kubectl diff` printing nothing.
