# Measuring the Cluster

Every Monday at 9am the shop gets hammered by people who just got paid, and by 10 it's quiet. Dex has been scaling brewd up and down by hand. Kubernetes can do this for him, but only if it can **measure** load first, and out of the box it can't:

```sh
kubectl top nodes
# error: Metrics API not available
```

**metrics-server** collects current CPU and memory use from every kubelet, every 15 seconds or so, and serves it through the Metrics API. It keeps no history, only "right now". (History is what Prometheus is for, in a later course.)

## Assignment

1. Install metrics-server. kind's kubelets use self-signed certificates, so it needs one extra flag to trust them. (Don't do this on a real cluster.)

   ```sh
   kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.9.0/components.yaml
   kubectl -n kube-system patch deployment metrics-server --type=json \
     -p '[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'
   kubectl -n kube-system rollout status deploy/metrics-server
   ```

2. Give it a minute to collect, then:

   ```sh
   kubectl top nodes
   kubectl top pods -n shop --sort-by=memory
   ```

3. Add one line to `answers.txt`:

   ```
   hungriest_pod=<the name of the shop pod using the most memory right now>
   ```

   This changes from minute to minute, so run `top` just before you submit.
