# A Deleted Pod Stays Deleted

Try this on your cluster (the pod is called `doomed` so nothing you need gets hurt):

```sh
kubectl run doomed --image=busybox:1.36 -- sleep 3600
kubectl get pods
kubectl delete pod doomed
kubectl get pods -w        # watch for 30 seconds. Ctrl-C to stop.
```

A pod you create yourself is a **bare pod**. You just watched what happens when one goes away.

Pods also go away without anyone typing `delete`: a node runs out of memory and evicts them, a node is drained for maintenance, or the machine underneath simply dies.

**At 3am a node is drained and the API server deletes the `storefront` pod you made in lesson 1. What happens next?**
