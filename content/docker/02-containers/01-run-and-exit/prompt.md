# Run a Container

The command you'll type most in Docker is `docker run`. Its shape:

```sh
docker run [options] IMAGE [command]
```

Random names like `brave_curie` are cute but impossible to remember. Give your container a name with `--name`:

```sh
docker run --name cq-demo alpine:3.20 echo "hi there"
# hi there
```

Now list running containers:

```sh
docker ps
# CONTAINER ID   IMAGE   COMMAND   CREATED   STATUS   PORTS   NAMES
```

Empty! Did it break? No. **A container only lives as long as the command inside it.** `echo` printed one line and finished, so the container finished too. It's not gone, just stopped. Add `-a` to see stopped containers as well:

```sh
docker ps -a
# ... alpine:3.20   "echo 'hi there'"   Exited (0) 5 seconds ago   cq-demo
```

Throughout this course, start every container name with **`cq-`**. That's how the checks find your work, and it keeps it apart from anything else on your computer.

# Assignment

1. In `~/docker-course`, run a container named **`cq-first`** from **`alpine:3.20`** that runs `echo "Hello from a container"`.
2. Find it with `docker ps -a`. Its status should say `Exited (0)`.
3. Run, then submit `dk.containers.run-and-exit`.

The checks look at `cq-first`'s image and status and read what it printed.
