# Config Outside the Image

brewd's database port and the banner on the shop page are baked into the image. Marketing wants the banner changed every week, and staging uses a different database port from production. Rebuilding an image to change one string is silly, and keeping one image per environment is worse.

A **ConfigMap** holds configuration as plain key-value pairs, separate from the image, so the same image runs everywhere with different settings.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: brewd-config
data:
  db-port: "5432"            # values are always strings: quote numbers
  greeting: "Fresh potions daily v1"
```

The imperative version is handy for experiments:

```sh
kubectl create configmap demo --from-literal=colour=green --from-file=nginx.conf
kubectl get configmap brewd-config -o yaml
kubectl describe configmap brewd-config
```

A ConfigMap is for configuration that's safe to read. Passwords and tokens go in Secrets (lesson 5).

## Assignment

Create **`manifests/brewd-config.yaml`** with a ConfigMap named **`brewd-config`** holding exactly:

- `db-port`: `5432`
- `greeting`: `Fresh potions daily v1`

Apply it. Nothing uses it yet; that's the next lesson.
