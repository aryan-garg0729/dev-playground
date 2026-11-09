# Memory-Only Storage

The third kind of mount is **tmpfs**: a folder that lives in the computer's **memory**, not on disk.

```sh
docker run -d --tmpfs /scratch alpine:3.20 sleep infinity
```

Why would you want that?

- **Speed**: memory is much faster than disk.
- **Secrets**: things like temporary keys are never written to disk where they could be found later.
- **Read-only containers**: in the production chapter you'll lock a container's files completely, and tmpfs gives the app one small place it may still write to.

A container writes some files into `/scratch` (a tmpfs mount). Then you run `docker stop` and `docker start` on it.

**What happens to those files?**
