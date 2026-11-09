# The LoadBalancer That Never Comes

Mo: *"On the cloud, don't you just say `type: LoadBalancer` and get an address?"*

On a cloud, yes. A **LoadBalancer** Service is a NodePort Service plus a request: "please create a real load balancer outside the cluster and point it at my nodes." The component that answers is the **cloud controller manager**, which each cloud (AWS, GCP, Azure...) writes for its own load balancers. On EKS or GKE, the same YAML gets a public address in about a minute.

```yaml
spec:
  type: LoadBalancer
  selector:
    app: brewd
  ports:
    - port: 80
      targetPort: 80
```

Your laptop has no cloud. Try it anyway, and watch:

```sh
kubectl get svc brewd-lb -w
```

## Assignment

1. Create a LoadBalancer Service named **`brewd-lb`** for the `app: brewd` pods, port 80 to target port 80. Don't set a `nodePort`.
2. Watch it for a minute and work out why it looks the way it does before reading the explanation at the end of the lesson.
3. Add two lines to `answers.txt`:

   ```
   lb_external=<exactly what the EXTERNAL-IP column says>
   lb_nodeport=<the nodePort Kubernetes gave it anyway>
   ```

The second one is the point: a LoadBalancer is a superset of a NodePort, which is a superset of a ClusterIP. Even with no cloud, `brewd-lb` still works as a NodePort and a ClusterIP.

(Tools like `cloud-provider-kind` or MetalLB can play the cloud's part on a laptop. We won't install one: seeing the empty column is the lesson.)
