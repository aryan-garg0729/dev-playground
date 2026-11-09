# The Cache Is Too Good

Mo is debugging a "security fix that won't install":

```dockerfile
1  FROM alpine:3.20
2  RUN apk update
3  RUN apk add curl
4  COPY potion-check.sh /usr/local/bin/
```

The security team says a fixed `curl` was published last week. Mo rebuilds on the build server, which has built this Dockerfile many times since spring. The new image still has the old, vulnerable `curl`.

Remember: a `RUN` step's cache key is its parent's key plus the **instruction text**.

**Why?**
