# Things You Shouldn't Commit

brewd needs a database password. Mo: *"Can we just put it in the ConfigMap?"* No. ConfigMaps are for things anyone on the team may read, and they end up in Git.

A **Secret** looks almost like a ConfigMap, but Kubernetes treats it differently: it's a separate kind, so access to it can be granted separately (chapter 11), it can be encrypted in etcd, and the kubelet keeps it in memory on the node instead of writing it to disk.

```sh
kubectl create secret generic brewd-db \
  --from-literal=username=brewd \
  --from-literal=password="$(openssl rand -hex 16)"
```

In a manifest, values under `data:` must be base64-encoded. `stringData:` takes plain text and encodes it for you:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: brewd-db
type: Opaque          # "arbitrary key-value data"; other types exist for TLS certs and registry logins
stringData:
  username: brewd
  password: change-me
```

The catch: a Secret manifest with a real password in it is exactly the file you must **not** commit. Teams create Secrets with `kubectl create secret` from a password manager, or use tools like Sealed Secrets or External Secrets that keep only an encrypted or referenced form in Git. Those tools are out of scope here; the habit isn't.

## Assignment

1. Create a Secret named **`brewd-db`** with:
   - `username`: **`brewd`**
   - `password`: a random value of **at least 16 characters** (the `openssl` command above makes a 32-character one)
2. Use `kubectl create secret`. **Don't** save it as a manifest in `manifests/`.

The checks never print your password. They read the username and the password's length only.
