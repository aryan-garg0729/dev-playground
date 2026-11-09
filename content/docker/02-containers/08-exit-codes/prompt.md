# Exit Codes

When a program finishes, it hands back a number called its **exit code**:

- **0** means "all good".
- **Anything else** means "something went wrong". The program picks the number.

That's the `(0)` in `Exited (0)`. A container whose status says `Exited (1)` or `Exited (3)` didn't finish happily, and the logs usually say why.

```sh
docker ps -a --format '{{.Names}} {{.Status}}'
# cq-first   Exited (0) 2 minutes ago    <- fine
# cq-crash   Exited (3) 5 seconds ago    <- trouble
```

To get just the number, ask Docker to **inspect** the container. `docker inspect` shows everything Docker knows about it, and `-f` picks one field:

```sh
docker inspect -f '{{.State.ExitCode}}' cq-crash
# 3
```

# Assignment

1. Run this container. It pretends to be an app that can't find its config file. Copy and paste it:

   ```sh
   docker run --name cq-crash alpine:3.20 sh -c 'echo "error: config file not found"; exit 3'
   ```

2. Find its exit code with `docker ps -a` and with `docker inspect`.
3. Read its logs to see why it failed.
4. Run, then submit `dk.containers.exit-codes`. **Don't remove `cq-crash`**: the checks need it.
