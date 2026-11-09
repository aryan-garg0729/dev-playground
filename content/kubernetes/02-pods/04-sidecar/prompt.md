# A Sidecar and a Shared Folder

The menu board in the shop window shows today's potions, and it has to say when it was last updated. The web server shouldn't have to know how to do that. A tiny helper container can write the file, and nginx can just serve it.

Two containers can't see each other's files by default. A **volume** declared on the pod can be mounted into several of its containers. The simplest one is an **`emptyDir`**: an empty folder created when the pod starts and deleted when the pod goes away.

```yaml
spec:
  volumes:
    - name: shared
      emptyDir: {}
  containers:
    - name: web
      image: nginx:1.27-alpine
      volumeMounts:
        - name: shared                  # which pod volume
          mountPath: /usr/share/nginx/html  # where it appears in THIS container
```

Each container chooses its own `mountPath`. The same folder can be `/usr/share/nginx/html` in one container and `/data` in another.

With more than one container, `logs` and `exec` need to know which one you mean:

```sh
kubectl logs menu-board -c clock
kubectl exec -it menu-board -c web -- sh
```

## Assignment

Write **`manifests/menu-board.yaml`**: a pod named **`menu-board`** with label `app: menu-board` and:

1. A volume named **`shared`**, of type `emptyDir`.
2. A container **`web`**: `nginx:1.27-alpine`, port 80, mounting `shared` at **`/usr/share/nginx/html`**.
3. A container **`clock`**: `busybox:1.36`, mounting `shared` at **`/data`**, running this command so it rewrites the page every 5 seconds:

   ```yaml
   command: ["sh", "-c", "while true; do echo \"<h1>Today's potions</h1><p>menu updated $(date)</p>\" > /data/index.html; sleep 5; done"]
   ```

Containers must be listed in that order, `web` first. Apply it, wait for `READY 2/2`, then port-forward to the pod and refresh a few times: the time changes, and nginx never knew.
