# Build Here, Ship There

A **multi-stage build** uses two (or more) `FROM` lines in one Dockerfile:

1. A **build stage** with all the heavy tools, where the program is made.
2. A small **final stage** that only gets the finished program, copied out of the build stage.

![Two stages: the build stage on golang:1.23-alpine has the Go toolchain, the source code and the build cache, and produces /out/webapp. Only that file is copied with COPY --from=build into the final stage on alpine:3.20, which becomes the image. Everything else is discarded.](assets/stages.svg)

```dockerfile
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY app/go.mod ./
RUN go mod download
COPY app/ .
RUN CGO_ENABLED=0 go build -o /out/webapp .

FROM alpine:3.20
COPY --from=build /out/webapp /usr/local/bin/webapp
ENTRYPOINT ["webapp"]
```

- `AS build` gives the first stage a name.
- `COPY --from=build` copies from that stage instead of from your folder.
- Only the **last** stage becomes the image. The build stage is thrown away.
- `CGO_ENABLED=0` makes Go build a fully self-contained program, so it runs on any small base. (We'll need that in two lessons.)

**One gotcha:** `ARG` belongs to a single stage. Declare it in the stage that uses it, here the final one.

# Assignment

In `~/docker-course`:

1. Replace your `Dockerfile` with the multi-stage version. Copy and paste:

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
   ENV GREETING=Welcome
   EXPOSE 8080
   ENTRYPOINT ["webapp"]
   EOF
   ```

   (We've dropped the start script to keep things short. The `exec` lesson still applies whenever you use one.)

2. Build it as **`cq/webapp:slim`** with `VERSION=1.5`.
3. Compare the sizes: `docker images cq/webapp`.
4. Run, then submit `dk.multistage.two-stages`.
