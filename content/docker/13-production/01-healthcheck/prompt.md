# A Built-In Health Check

"Running" only means the program hasn't exited. A server can be running and completely stuck. In Compose you added a healthcheck to one service. Put it in the **image** instead, and every container made from it checks itself, wherever it runs:

```dockerfile
HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
  CMD wget -qO- http://localhost:8080/healthz || exit 1
```

Every 5 seconds Docker runs that command inside the container. After 3 failures in a row the container is marked **unhealthy**. You'll see it in `docker ps`:

```
STATUS
Up 20 seconds (healthy)
```

Tools like Compose, Docker Swarm and Kubernetes use this to decide when to send traffic to a container and when to replace it.

# Assignment

In `~/docker-course`:

1. Replace your `Dockerfile` (the same as in the Compose chapter, plus a `HEALTHCHECK`). Copy and paste:

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
   ENTRYPOINT ["webapp"]
   EOF
   ```

2. Build it as **`cq/webapp:health`**.
3. Run **`cq-health`** from it, detached. Run `docker ps` a few times: `(health: starting)`, then `(healthy)`.
4. Run, then submit `dk.production.healthcheck`.
