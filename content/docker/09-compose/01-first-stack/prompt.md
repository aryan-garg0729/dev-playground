# The Stack in One File

By now, starting the shop means remembering a page of `docker run` flags. **Compose** writes them down instead: one YAML file describes every container (a **service**), and one command starts them all.

```yaml
# compose.yaml
services:
  web:
    build: .                       # build the Dockerfile in this folder...
    image: demo/web:dev            # ...and tag the result with this name
    ports:
      - "18000:8080"               # same host:container as -p
```

```sh
docker compose up -d --build   # build if needed, create, start (detached)
docker compose ps              # this project's containers
docker compose logs -f web     # follow one service's logs
docker compose down            # stop and remove the containers (and the project's network)
```

Compose groups everything under a **project name**, which defaults to the folder name, so containers come out as `elixirbytes-web-1`. It also creates a network for the project, `elixirbytes_default`, where **service names resolve via DNS**, just like chapter 8.

`docker compose ps -q SERVICE` prints a service's container ID, which is handy inside other commands.

## Assignment

1. In `~/elixirbytes`, create **`compose.yaml`** with one service named **`api`** that:
   - builds from your `Dockerfile` (the chapter 6 multi-stage one)
   - tags the image **`elixirbytes/brewd:compose`**
   - publishes container port 8080 on host port **18091**
2. Start it with `docker compose up -d --build` and check `curl localhost:18091/healthz`.
3. Submit from `~/elixirbytes`, where `compose.yaml` lives: the checks use `docker compose` and need to find it.
