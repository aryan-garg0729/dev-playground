# Describe First

Dex's on-call notes, the part he wants everyone to memorise:

> Every broken workload is one of a handful of shapes. Find the shape before you touch anything.
>
> | you see | it means | look at |
> | --- | --- | --- |
> | `Pending` | the scheduler can't place it (resources, taints, selectors, or a claim that isn't bound) | `describe pod` → FailedScheduling |
> | `ImagePullBackOff` / `ErrImagePull` | the image can't be pulled | `describe pod` → Events |
> | `CreateContainerConfigError` | a ConfigMap, Secret or key it needs is missing | `describe pod` → Events |
> | `ContainerCreating` for minutes | usually a volume that's bound but can't attach or mount | `describe pod`, then `describe pvc` |
> | `CrashLoopBackOff` | it starts and dies | `logs --previous`, then `describe` for the exit code |
> | `Running` but `0/1` | readiness probe failing | `describe pod` → Unhealthy events |
> | `1/1 Running`, Service serves nothing | the Service doesn't reach the pods | endpoints, selector, `targetPort` |
>
> Fix the **cluster**, then fix the YAML in Git so the next deploy doesn't undo you.

Events are kept for about an hour. `kubectl get events --sort-by=.lastTimestamp` shows a namespace's recent history in order.

From now on, every lesson starts with something already broken. Dex's first rule: **don't delete anything until you know why it's broken.**

**A pod has been `Pending` for ten minutes. What's your first command?**
