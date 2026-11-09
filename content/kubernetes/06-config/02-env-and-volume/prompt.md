# Two Ways In

A pod can read a ConfigMap in two ways, and the difference will matter in the very next lesson.

**As environment variables**, one key at a time:

```yaml
containers:
  - name: brewd
    env:
      - name: DB_PORT                 # the variable's name in the container
        valueFrom:
          configMapKeyRef:
            name: brewd-config        # which ConfigMap
            key: db-port              # which key
```

(`envFrom` with a `configMapRef` imports every key at once, but keys like `db-port` aren't valid variable names, so we'll name them one by one.)

**As files**, by mounting the ConfigMap as a volume. Each key becomes a file named after the key:

```yaml
spec:
  volumes:
    - name: config
      configMap:
        name: brewd-config
  containers:
    - name: brewd
      volumeMounts:
        - name: config
          mountPath: /etc/brewd       # -> /etc/brewd/db-port, /etc/brewd/greeting
          readOnly: true
```

Both of these go in the Deployment's **pod template**, so changing them is a rollout: new ReplicaSet, new pods.

## Assignment

Edit **`manifests/brewd.yaml`** so the `brewd` container gets:

1. An environment variable **`DB_PORT`** from key `db-port` of `brewd-config`.
2. The whole `brewd-config` ConfigMap mounted at **`/etc/brewd`**.

Apply it, wait for `kubectl rollout status deploy/brewd`, then check from inside a pod:

```sh
kubectl exec deploy/brewd -- sh -c 'echo $DB_PORT; ls /etc/brewd; cat /etc/brewd/greeting'
```

`exec deploy/brewd` picks one of the Deployment's pods for you, so you don't have to look up a name.
