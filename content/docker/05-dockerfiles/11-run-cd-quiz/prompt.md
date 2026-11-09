# Every RUN Starts Fresh

One gotcha catches nearly everyone. Look at this Dockerfile:

```dockerfile
FROM alpine:3.20
RUN mkdir -p /tmp/app
RUN cd /tmp/app
RUN pwd
```

`pwd` prints the folder you're in.

Every `RUN` is its own little shell session. Files you create **are** kept (they're saved in a layer), but things like "which folder am I in" are forgotten as soon as the step ends.

The fix is `WORKDIR`, which changes the folder for **every** step after it, and for the container too:

```dockerfile
WORKDIR /tmp/app
# prints /tmp/app
RUN pwd
```

(Or do it in one step: `RUN cd /tmp/app && pwd`.)

**What does the original `RUN pwd` print?**
