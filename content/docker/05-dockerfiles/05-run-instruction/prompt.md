# Install Things with RUN

Alpine is tiny because it has almost nothing in it. Want a tool like `curl` (which downloads web pages)? Install it with Alpine's package manager, `apk`, in a `RUN` step:

```dockerfile
FROM alpine:3.20
RUN apk add --no-cache curl
```

During `docker build`, Docker starts a temporary container, runs `apk add` inside it, and saves the result as a new layer of your image. (`--no-cache` stops apk from keeping its download list, which keeps the image small.)

Different base images have different package managers: Debian and Ubuntu images use `apt-get`, Alpine uses `apk`.

# Assignment

1. Create `tools/Dockerfile` in `~/docker-course`. Copy and paste:

   ```sh
   mkdir -p tools
   cat > tools/Dockerfile <<'EOF'
   FROM alpine:3.20
   RUN apk add --no-cache curl
   EOF
   ```

2. Build it as **`cq/tools:v1`**. Watch the `RUN` step download curl.
3. Try it: `docker run --rm cq/tools:v1 curl --version`.
4. Run, then submit `dk.dockerfile.run-instruction`.

The checks also confirm that `alpine:3.20` itself was never changed. Your image is a new layer on top.
