# Percent of What?

A **HorizontalPodAutoscaler** (HPA) watches a metric and edits a Deployment's `replicas` for you. "Horizontal" means more pods, as opposed to bigger ones. Its formula is short:

```
desired replicas = ceil( current replicas × current metric / target metric )
```

With 2 pods at 90% CPU and a 60% target: `ceil(2 × 90 / 60) = 3`.

Mo creates an HPA for brewd, which has no `resources` set:

```sh
kubectl autoscale deploy brewd --min=2 --max=8 --cpu=50%
kubectl get hpa
# NAME    REFERENCE          TARGETS              MINPODS   MAXPODS   REPLICAS
# brewd   Deployment/brewd   cpu: <unknown>/50%   2         8         4
```

It stays `<unknown>` forever, though `kubectl top pods` shows brewd using CPU.

**Why?**
