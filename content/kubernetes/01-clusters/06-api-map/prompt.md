# Map of the API

Kubernetes has dozens of kinds of object, each with hundreds of fields. Nobody memorises them, and blog posts go stale with every release. The cluster carries its own documentation, and it's always right for *your* version.

**What kinds exist?**

```sh
kubectl api-resources
# NAME          SHORTNAMES   APIVERSION   NAMESPACED   KIND
# pods          po           v1           true         Pod
# services      svc          v1           true         Service
# deployments   deploy       apps/v1      true         Deployment
# ...
```

- **SHORTNAMES** save typing: `kubectl get po` is `kubectl get pods`.
- **APIVERSION** is what goes in a manifest's `apiVersion:` line.
- **NAMESPACED** says whether the kind lives inside a namespace or belongs to the whole cluster.

**What fields does a kind have?**

```sh
kubectl explain pod.spec.containers          # one field, with its docs
kubectl explain pod.spec --recursive | less  # the whole tree
```

`kubectl explain` answers "what was that field called again?" faster than any search engine.

## Assignment

Add three lines to `answers.txt`, using only those two commands:

```
cronjob_short=<the short name for cronjobs>
node_namespaced=<true or false: is Node a namespaced kind?>
restart_default=<the restartPolicy a pod gets if you don't set one>
```

For the last one, read the description `kubectl explain pod.spec.restartPolicy` prints. Copy the value with the capitalisation it uses.
