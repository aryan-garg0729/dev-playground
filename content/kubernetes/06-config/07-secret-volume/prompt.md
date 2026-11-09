# Secrets as Read-Only Files

A Secret can reach a container the same two ways as a ConfigMap: environment variables (`secretKeyRef`) or a volume. Prefer the **volume**. Environment variables leak:

- into crash reports and error pages that dump the environment
- into every child process the app starts
- through `/proc/<pid>/environ`, to anything on the node that can read it
- through `kubectl describe pod`, which shows env sources, and logging libraries that print "config at startup"

A file can be locked down, and like a ConfigMap volume, it updates when the Secret changes.

```yaml
spec:
  volumes:
    - name: creds
      secret:
        secretName: brewd-db
        defaultMode: 0400       # owner can read; nobody else, and nobody can write
  containers:
    - name: brewd
      volumeMounts:
        - name: creds
          mountPath: /etc/creds  # -> /etc/creds/username, /etc/creds/password
          readOnly: true
```

**One YAML trap:** `0400` with a leading zero is an *octal* number, which YAML reads as 256. That's correct, and the API stores and shows it as `256`. Without the leading zero, `400` would be decimal 400, which is a different and odd permission.

## Assignment

1. Edit `manifests/brewd.yaml` to mount the **`brewd-db`** Secret at **`/etc/creds`**, `readOnly: true`, with `defaultMode: 0400`. Don't add any env vars from it.
2. Apply it and wait for the rollout.
3. Look from inside a pod: `kubectl exec deploy/brewd -- ls -lL /etc/creds` (`-L` follows the symlinks Kubernetes uses to swap files atomically on update).
