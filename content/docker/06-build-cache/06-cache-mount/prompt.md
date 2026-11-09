# A Cache That Survives

Layer caching is all or nothing. Once `COPY app/ .` changes, `RUN go build` starts from zero and compiles **everything** again, even though tools like Go keep their own cache of work they've already done.

A **cache mount** gives a `RUN` step a folder that **survives between builds** and is **not** saved in the image:

```dockerfile
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /usr/local/bin/webapp .
```

- `target` is the folder inside the build where the tool keeps its cache (for Go, `/root/.cache/go-build`; for pip, `/root/.cache/pip`; for npm, `/root/.npm`).
- The first build fills it. Later builds find their old work there and finish faster.
- The folder never ends up in your image, so the image stays small.

(The `\` at the end of a line just continues the step on the next line.)

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
   RUN --mount=type=cache,target=/root/.cache/go-build \
       go build -o /usr/local/bin/webapp .
   ARG VERSION=dev
   ENV APP_VERSION=$VERSION
   LABEL org.opencontainers.image.version=$VERSION
   ENV GREETING=Welcome
   EXPOSE 8080
   CMD ["webapp"]
   EOF
   ```

2. Build it as **`cq/webapp:fast`**.
3. Run, then submit `dk.cache.cache-mount`.

The last check looks inside the image for the compiler cache and expects to find nothing.
