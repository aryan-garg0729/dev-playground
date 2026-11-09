# The ReplicaSet Counts Labels

`brewd-rs` is running 3 pods. Mo, trying to help, creates one more pod by hand:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: extra-brewd
  labels:
    app: brewd-rs
spec:
  containers:
    - name: brewd
      image: nginx:1.27-alpine
```

Remember the loop: *count the pods that match my selector, then fix the difference.*

**What does the ReplicaSet do about `extra-brewd`?**
