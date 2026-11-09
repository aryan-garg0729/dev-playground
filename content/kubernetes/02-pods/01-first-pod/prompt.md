# Your First Pod

Kubernetes doesn't run containers directly. It runs **pods**: a pod is one or more containers that are always scheduled together, on the same node, sharing one network address. Most pods hold a single container, and you can think of "pod" as "a running container plus the box Kubernetes keeps it in".

You could start one with a single command:

```sh
kubectl run storefront --image=nginx:1.27-alpine
```

That's imperative and gone from memory as soon as you close the terminal. We'll write it down instead. A manifest has four top-level parts, and every kind of object uses them:

```yaml
apiVersion: v1          # which API group/version this kind lives in (kubectl api-resources)
kind: Pod               # what kind of object
metadata:
  name: storefront      # its name, unique in the namespace
  labels:
    app: storefront     # tags for finding it later (lesson 7)
spec:                   # what you want: the "desired state"
  containers:
    - name: web
      image: nginx:1.27-alpine
      ports:
        - containerPort: 80
```

`containerPort` is documentation for humans and tools; it doesn't open anything. The Kubernetes object also grows a `status:` section that you never write: it's the cluster reporting the *actual* state.

```sh
kubectl apply -f manifests/storefront.yaml
kubectl get pods
# NAME         READY   STATUS    RESTARTS   AGE
# storefront   1/1     Running   0          8s
```

In production this image would be `elixirbytes/brewd`. We'll use public stand-in images all course, so the cluster works without the registry from the Docker course.

## Assignment

The shop's front page needs somewhere to live.

1. Save the manifest above as **`~/elixirbytes/k8s/manifests/storefront.yaml`**.
2. Apply it to the `shop` namespace.
3. Wait until `kubectl get pods` shows `storefront` as `Running`. The first run pulls the image, so it can sit at `ContainerCreating` for a few seconds.

If it's stuck, run `kubectl describe pod storefront` and read the **Events** at the bottom. That's the first move for every problem in this course.
