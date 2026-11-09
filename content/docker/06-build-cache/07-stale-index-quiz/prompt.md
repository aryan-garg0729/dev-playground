# The Stale Package List

The cache only looks at the **text** of a step (and the files it copies). It has no idea whether the outside world has changed.

That causes a famous bug in Debian- and Ubuntu-based images. Their package manager, `apt-get`, needs an up-to-date **list** of available packages before it can install anything:

```dockerfile
FROM debian:12-slim
# download the package list
RUN apt-get update
# install from that list
RUN apt-get install -y git
```

This works the first time. Months later, you add `curl`:

```dockerfile
RUN apt-get update
RUN apt-get install -y git curl
```

The build fails with `404 Not Found` errors while downloading packages.

The fix is always to update and install **in the same step**:

```dockerfile
RUN apt-get update && apt-get install -y git curl
```

(Alpine's `apk add --no-cache` fetches a fresh list every time, which is one reason you've been using it.)

**Why did the second build fail?**
