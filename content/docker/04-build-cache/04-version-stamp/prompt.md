# Stamp the Version

Dex's question during last week's outage was "which version is running?", and nobody knew. Every image we ship should say what it is, in two places:

1. **In the binary**, so `brewd --version` and `GET /version` answer.
2. **On the image**, as a **label**: a `key=value` note stored in the image's metadata, readable without running anything. `org.opencontainers.image.version` is the standard key that registries and scanners understand.

Go can overwrite a variable at compile time. brewd's `var version = "dev"` becomes `1.4.0` with:

```sh
go build -ldflags "-X main.version=1.4.0" -o /usr/local/bin/brewd .
```

To pass a value into a build, declare an `ARG` and use it like a shell variable:

```dockerfile
ARG BUILD_VERSION=dev
RUN echo "building $BUILD_VERSION"
LABEL org.opencontainers.image.version=$BUILD_VERSION
```

```sh
docker build --build-arg BUILD_VERSION=1.4.0 -t demo .
docker image inspect -f '{{json .Config.Labels}}' demo
# {"org.opencontainers.image.version":"1.4.0"}
```

**Cache tip:** declare the `ARG` just before the first step that uses it. A changed `ARG` value invalidates the cache from its first use, and you don't want a version bump to re-download dependencies.

## Assignment

1. Add an **`ARG BUILD_VERSION`** with default **`dev`** to your Dockerfile, after the dependency download.
2. Use it twice: stamp it into the binary with `-ldflags "-X main.version=..."`, and set the label **`org.opencontainers.image.version`**.
3. Build **`elixirbytes/brewd:v4`** with `BUILD_VERSION` set to **`1.4.0`**.
4. Submit from `~/elixirbytes`. The checks run `brewd --version`, read the label, and confirm `BUILD_VERSION` isn't in the container's environment, because an `ARG` should never be.
