# Don't Run as Root

Remember `whoami` inside nginx? It said **root**. Containers run as the all-powerful root user unless the image says otherwise.

The walls around a container are good, but not perfect. If someone breaks into your app, being root inside makes their next step much easier. So production images run the app as an ordinary user:

```dockerfile
RUN adduser -D -u 10001 app && mkdir -p /data && chown app /data
USER app
```

1. `adduser -D -u 10001 app`: create a user called `app` with the ID number 10001 (`-D`: no password).
2. `mkdir -p /data && chown app /data`: the app needs one folder it may write to, so make `app` its owner.
3. `USER app`: every later step, and the container itself, runs as `app`.

Now the app can't change its own program or the system files, even if it's hacked.

# Assignment

In `~/docker-course`:

1. In your `Dockerfile`'s **final stage**, add those two lines just before `ENTRYPOINT`. Copy and paste to replace the whole file:

   ```sh
   cat > Dockerfile <<'EOF'
   FROM golang:1.23-alpine AS build
   WORKDIR /src
   COPY app/go.mod ./
   RUN go mod download
   COPY app/ .
   RUN --mount=type=cache,target=/root/.cache/go-build \
       CGO_ENABLED=0 go build -o /out/webapp .

   FROM alpine:3.20
   COPY --from=build /out/webapp /usr/local/bin/webapp
   ARG VERSION=dev
   ENV APP_VERSION=$VERSION
   LABEL org.opencontainers.image.version=$VERSION
   EXPOSE 8080
   HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
     CMD wget -qO- http://localhost:8080/healthz || exit 1
   RUN adduser -D -u 10001 app && mkdir -p /data && chown app /data
   USER app
   ENTRYPOINT ["webapp"]
   EOF
   ```

2. Build it as **`cq/webapp:prod`**.
3. Check: `docker run --rm --entrypoint id cq/webapp:prod`.
4. Run, then submit `dk.production.non-root`.
