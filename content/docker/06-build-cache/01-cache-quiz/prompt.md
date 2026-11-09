# Why Was That So Fast?

Build the same image twice and the second build finishes almost instantly. Every step says **`CACHED`**:

```
=> CACHED [2/4] WORKDIR /src
=> CACHED [3/4] COPY app/ .
=> CACHED [4/4] RUN go build -o /usr/local/bin/webapp .
```

Docker remembers the layer each step made. Before running a step, it asks: **"Have I done exactly this before, with exactly these inputs?"**

- Same instruction, and (for `COPY`) the same files: reuse the saved layer. Instant.
- Anything different: run the step for real.

There's one catch. Layers stack on top of each other, so **once one step has to run again, every step after it runs again too.** The cache is a chain, and a broken link breaks everything below it.

Here's a Dockerfile:

```dockerfile
# step 1
FROM nginx:1.27-alpine
# step 2
RUN apk add --no-cache curl
# step 3
COPY index.html /usr/share/nginx/html/
# step 4
RUN echo "built" > /usr/share/nginx/html/built.txt
```

**You edit `index.html` and build again. Which steps actually run?**
