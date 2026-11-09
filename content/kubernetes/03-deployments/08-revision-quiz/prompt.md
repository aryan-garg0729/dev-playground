# Why the Revision Went Up

You rolled out 1.28 and then undid it. Your history now looks like this:

```sh
kubectl rollout history deploy/brewd
# REVISION  CHANGE-CAUSE
# 2         <none>
# 3         <none>
```

Revision 1 is gone, and the current state is revision **3**, even though it runs exactly what revision 1 ran.

**Why?**
