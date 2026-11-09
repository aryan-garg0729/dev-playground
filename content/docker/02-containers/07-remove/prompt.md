# Clean Up

Stopped containers stick around until you remove them. They're small, but they pile up fast. Delete one with `docker rm`:

```sh
docker rm cq-demo
# cq-demo
```

Docker refuses to remove a **running** container, just to be safe:

```sh
docker rm cq-one
# Error response from daemon: cannot remove container "cq-one": container is running
```

Either stop it first, or force it with `-f` (stop + remove in one go):

```sh
docker rm -f cq-one
```

You can remove several at once: `docker rm -f cq-two cq-three`.

And for throwaway containers, skip the cleanup completely: `docker run --rm ...` deletes the container the moment it finishes. You've used it since chapter 1.

# Assignment

1. Remove **`cq-first`** (it's already stopped).
2. Remove **`cq-one`**, **`cq-two`** and **`cq-three`** (they're running).
3. Keep **`cq-sleeper`**.
4. Run, then submit `dk.containers.remove`.

Removing a container never removes its image. `alpine:3.20` is still on your computer, ready for the next run.
