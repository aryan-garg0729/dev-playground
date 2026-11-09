# When a Job Gives Up

A Job's `restartPolicy` decides *how* it retries:

- **`OnFailure`**: the same pod's container is restarted in place.
- **`Never`**: the failed pod is left as evidence, and a **new pod** is created for the retry.

`backoffLimit` is how many retries the Job gets before it gives up (default 6).

Mo's import Job has a bug: its container always exits with code 1.

```yaml
spec:
  completions: 5
  backoffLimit: 3
  template:
    spec:
      restartPolicy: Never
```

**How many pods does this Job create, and how does it end?**
