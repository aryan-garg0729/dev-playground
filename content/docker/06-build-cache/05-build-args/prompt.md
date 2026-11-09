# Pick the Version at Build Time

Let's make the version a build argument. We'll also add a **label**: a note stuck onto the image that tools can read. `org.opencontainers.image.version` is the standard name for "which version is this image".

```dockerfile
ARG VERSION=dev
ENV APP_VERSION=$VERSION
LABEL org.opencontainers.image.version=$VERSION
```

Read a label back with:

```sh
docker image inspect -f '{{index .Config.Labels "org.opencontainers.image.version"}}' IMAGE
```

# Assignment

In `~/docker-course`:

1. Replace your `Dockerfile`. Copy and paste:

   ```sh
   cat > Dockerfile <<'EOF'
   FROM golang:1.23-alpine
   WORKDIR /src
   COPY app/go.mod ./
   RUN go mod download
   COPY app/ .
   RUN go build -o /usr/local/bin/webapp .
   ARG VERSION=dev
   ENV APP_VERSION=$VERSION
   LABEL org.opencontainers.image.version=$VERSION
   ENV GREETING=Welcome
   EXPOSE 8080
   CMD ["webapp"]
   EOF
   ```

   `ARG` comes **after** the build steps on purpose: a changed build argument breaks the cache from that line down, so placing it late keeps the slow steps cached.

2. Build it as **`cq/webapp:1.2`** with the build argument **`VERSION=1.2`**.
3. Build it again as **`cq/webapp:nover`**, this time with **no** build argument.
4. Check: `docker run --rm cq/webapp:1.2 webapp --version`.
5. Run, then submit `dk.cache.build-args`.
