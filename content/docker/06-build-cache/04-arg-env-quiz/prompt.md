# ARG vs ENV

Right now the version number is written straight into the Dockerfile: `ENV APP_VERSION=1.1`. Releasing 1.2 means editing the file. It would be nicer to choose the version when you **build**.

That's what `ARG` is for. It declares a **build argument**: a value you pass to `docker build`.

```dockerfile
# "dev" is the default if you don't pass one
ARG VERSION=dev
RUN echo "building $VERSION"
```

```sh
docker build --build-arg VERSION=1.2 -t myapp:1.2 .
```

| | `ARG` | `ENV` |
| --- | --- | --- |
| Set with | `--build-arg` when you build | the Dockerfile (or `-e` when you run) |
| Exists during the build | Yes | Yes |
| Exists inside running containers | **No** | Yes |

Want a build argument to reach the app? Copy it into an `ENV`:

```dockerfile
ARG VERSION=dev
ENV APP_VERSION=$VERSION
```

**You want to choose a version number each time you build, without it hanging around as an extra setting in every container. Which do you use?**
