# Exit Code 42

When a container's main process ends, Docker records its **exit code**: the number every program hands back to the operating system when it finishes. `0` means success. Anything from `1` to `125` means the program itself decided something went wrong.

```sh
docker run --name eb-demo alpine:3.20 sh -c 'echo brewing; exit 3'
# brewing

docker ps -a --filter name=eb-demo
# ... Exited (3) 5 seconds ago ...

docker inspect -f '{{.State.ExitCode}}' eb-demo
# 3
```

`docker inspect` dumps everything Docker knows about a container as JSON. The `-f` flag takes a Go template, `{{.Path.To.Field}}`, and pulls out one value. Get used to it; you'll use it all course.

A stopped container isn't deleted. It stays in `docker ps -a`, along with its exit code, logs and writable layer, until you `docker rm` it. That's what lets you investigate a crash after it happens. (Add `--rm` to `docker run` when you *want* it cleaned up automatically.)

## Assignment

Dex wants a reproducible crash to test alerting with.

1. Run a container named **`eb-crash`** from **`alpine:3.20`** that exits on its own with exit code **42**. `sh -c 'exit 42'` is one way to do it.
2. **Leave it in the `exited` state.** Don't use `--rm` and don't remove it afterwards: it's evidence.
3. Confirm it with `docker inspect -f '{{.State.ExitCode}}' eb-crash`, then submit.
