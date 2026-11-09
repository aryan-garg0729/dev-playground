# Run in the Background

To keep a container alive, give it a command that doesn't end. `sleep infinity` waits forever, which is perfect for practice:

```sh
docker run --name cq-demo2 alpine:3.20 sleep infinity
```

It works, but now your terminal is stuck: the container is attached to it, and you can't type. (Open a new terminal and run `docker rm -f cq-demo2` to get out.)

Add **`-d`** (*detached*) to run the container in the background and get your terminal back:

```sh
docker run -d --name cq-demo3 alpine:3.20 sleep infinity
# 7c1e0f3a9b2d...   <- the new container's id
docker ps
# CONTAINER ID   IMAGE         COMMAND            STATUS         NAMES
# 7c1e0f3a9b2d   alpine:3.20   "sleep infinity"   Up 3 seconds   cq-demo3
```

`Up 3 seconds` means it's running. Servers, databases and websites are almost always run with `-d`.

# Assignment

1. Start a container named **`cq-sleeper`** from **`alpine:3.20`**, **detached**, running `sleep infinity`.
2. Check `docker ps` shows it as `Up`.
3. Run, then submit `dk.containers.background`.

If the name is already taken from an earlier try, remove the old one with `docker rm -f cq-sleeper` and start again. Leave `cq-sleeper` running: the next lesson uses it.
