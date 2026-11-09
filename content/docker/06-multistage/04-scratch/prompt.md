# Nothing at All

`FROM scratch` is the empty image: no files, no shell, no `/etc`. Copy in a **static binary** (one that needs no shared libraries) and that's the whole image. `golang:1.23-alpine` builds static binaries by default (`go env CGO_ENABLED` prints `0` there), so brewd qualifies.

Mo tried it and got this:

```dockerfile
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY app/ .
RUN go build -o /out/brewd .

FROM scratch
COPY --from=build /out/brewd /brewd
CMD /brewd
```

```sh
docker run --rm elixirbytes/brewd:scratch
# docker: Error response from daemon: ... exec: "/bin/sh": stat /bin/sh: no such file or directory
```

The binary is there. What's missing is the **shell** that shell-form `CMD /brewd` quietly asked for. Any shell form (`CMD`, `ENTRYPOINT` or `RUN`) needs `/bin/sh` in the image, and so does any wrapper script.

## Assignment

1. Save Mo's Dockerfile as **`Dockerfile.scratch`** and fix it. Run **`/brewd`** as the image's **`ENTRYPOINT`**, in exec form, so that `docker run IMAGE --version` passes `--version` to brewd.
2. Build it as **`elixirbytes/brewd:scratch`**, then try `docker run --rm elixirbytes/brewd:scratch --version`.
3. Submit from `~/elixirbytes`. The checks read the entrypoint and layer count, run brewd with `--version`, and try (and should fail) to start a shell.
