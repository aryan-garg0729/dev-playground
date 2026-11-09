# Logs, Exec and Port-Forward

The pod is `Running`. Is it actually serving the shop? You can't SSH into a pod, but you don't need to. Three commands cover almost everything:

**Read what the container printed.** Anything it writes to stdout or stderr is kept by the node:

```sh
kubectl logs storefront
kubectl logs storefront -f            # follow, like tail -f
kubectl logs storefront --previous    # the run BEFORE the last restart (lesson 5)
```

**Run a command inside the container**, the same idea as `docker exec`:

```sh
kubectl exec storefront -- nginx -v
kubectl exec -it storefront -- sh     # an interactive shell; type exit to leave
```

Everything after `--` is the command to run in the container, not flags for kubectl.

**Reach the pod from your laptop.** The pod's IP only works inside the cluster. `port-forward` tunnels a local port to it through the API server, for as long as the command runs:

```sh
kubectl port-forward pod/storefront 8080:80     # local 8080 -> pod port 80
# in a second terminal:
curl localhost:8080
```

port-forward is a debugging tool, not how customers will reach the shop. That's chapter 4.

## Assignment

Mo wants the default nginx page gone *right now*, before anyone sees it.

1. Using `kubectl exec`, overwrite `/usr/share/nginx/html/index.html` inside `storefront` so it contains the text **`Welcome to ElixirBytes`**. (Inside the pod's shell: `echo '<h1>Welcome to ElixirBytes</h1>' > /usr/share/nginx/html/index.html`.)
2. Port-forward to the pod and `curl` the page from your laptop. You should see your heading.
3. Check `kubectl logs storefront`: your request is there, logged by nginx.

Notice what you just did: you changed a running container by hand. If this pod is ever replaced, your edit is gone. Keep that thought; it comes back in lesson 6.
