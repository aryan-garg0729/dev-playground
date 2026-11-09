# Step Inside

`docker run` starts a **new** container. To work inside one that's already running, use `docker exec`. It starts an extra process next to the main one, in the same isolated world:

```sh
docker exec eb-hello cat /etc/os-release
# NAME="Alpine Linux"
# ...

docker exec -it eb-hello sh
# / # ps
# PID   USER     COMMAND
#     1 root     sleep infinity
#     7 root     sh
#    13 root     ps
```

`-i` keeps input open and `-t` gives you a terminal, so together they make an interactive shell. Type `exit` to leave. Only your shell ends; `sleep infinity` keeps going as PID 1.

Notice the PIDs. Inside the container, `sleep` believes it's process 1, the first process on the machine. On the host it's just another process with a big PID. That's a namespace at work.

## Assignment

Every ElixirBytes container gets a badge file so Dex knows who set it up.

1. **Without restarting or recreating** `eb-hello`, create the file **`/etc/eb-badge`** inside it, containing the word **`apprentice`**. Use an interactive shell or a one-liner: `docker exec eb-hello sh -c '...'`.
2. Check it with `docker exec eb-hello cat /etc/eb-badge`.
3. Run `codequest submit dk.containers.exec`.

The checks confirm `eb-hello` is still running and was never restarted, then read the badge.
