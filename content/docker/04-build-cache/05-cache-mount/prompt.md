# A Cache That Survives

Layer caching is all or nothing: if `COPY app/ .` changes, the `go build` step runs from scratch, and the Go compiler recompiles every package it already compiled last time. On a real codebase that's minutes.

BuildKit (Docker's builder) can **mount a persistent cache directory** into a single `RUN` step:

```dockerfile
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /usr/local/bin/brewd .
```

- The directory at `target` survives between builds, kept by the builder on your machine.
- It is **never** part of the image. When the step ends, it's unmounted.
- So the step still re-runs when the source changes, but the compiler only rebuilds what changed.

`/root/.cache/go-build` is where Go keeps its compiler cache (`go env GOCACHE`). Each tool has its own: `/root/.npm`, `/root/.cache/pip`, `/var/cache/apk`.

To prove a cache hit, save the build output. `--progress=plain` prints every step as plain text instead of the animated view:

```sh
docker build --progress=plain -t demo . 2>&1 | tee build.log
```

## Assignment

1. Add a cache mount for **`/root/.cache/go-build`** to the `go build` step in your Dockerfile. Keep the version stamping from the last lesson.
2. Build `elixirbytes/brewd:v4` once (version `1.4.0` again).
3. Edit a comment in `app/main.go`, then build again with `--progress=plain`, saving the output to **`~/elixirbytes/build.log`**. Use the `tee` line above.
4. Submit from `~/elixirbytes`. The checks read your Dockerfile and `build.log` (the dependency step should be `CACHED`, and the build step should show its mount), then look inside the image to confirm the cache didn't end up in it.
