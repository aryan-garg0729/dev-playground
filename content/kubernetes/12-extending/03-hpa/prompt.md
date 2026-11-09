# Scaling on Demand

Time to replace Dex's Monday-morning routine.

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: rush
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: rush
  minReplicas: 1
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 50    # percent of the pods' CPU REQUEST
```

The HPA re-checks about every 15 seconds. It scales **up** quickly, but waits about **5 minutes** of low load before scaling **down** (`behavior.scaleDown.stabilizationWindowSeconds`), so a short dip doesn't make it flap.

Once an HPA manages a Deployment, take `replicas` out of the Deployment's manifest. Otherwise every `kubectl apply` resets the count, and the two fight.

## Assignment

1. Create a Deployment **`rush`** (label `app: rush`, `nginx:1.27-alpine`, 1 replica) with **`requests.cpu: 100m`** and **`limits.cpu: 200m`**, and a ClusterIP Service **`rush`** on port 80.
2. Create the HPA above.
3. Watch it in one terminal: `kubectl get hpa rush -w`. In another, generate load:

   ```sh
   kubectl run load --rm -it --image=busybox:1.36 --restart=Never -- \
     sh -c "while true; do wget -q -O /dev/null http://rush; done"
   ```

   Within a minute or two, `TARGETS` climbs past 50% and `REPLICAS` goes up.
4. Stop the load (Ctrl-C) and read `kubectl describe hpa rush`. Its events narrate every decision it made. Submit any time after it has scaled up; it doesn't have to scale back down first.
