# Named Volumes

A **volume** is storage that Docker manages outside any container's layers. You mount it at a path inside the container, and whatever is written there lives in the volume, which outlives every container that uses it:

```sh
docker volume create eb-demo
docker run --rm -v eb-demo:/stuff alpine:3.20 sh -c 'echo hello > /stuff/note'
docker run --rm -v eb-demo:/stuff alpine:3.20 cat /stuff/note
# hello                      <- a different container, same data
docker volume ls
docker volume rm eb-demo
```

`-v NAME:/path` mounts the volume called `NAME` at `/path`. Docker creates the volume if it doesn't exist. Where the bytes actually live on disk is Docker's business (on macOS it's inside the Docker VM), which is the point: you manage volumes by name.

## Assignment

1. Create a volume named **`eb-orders`**.
2. Run **`eb-store`** from `elixirbytes/brewd:1.6.0`, detached, with `eb-orders` mounted at **`/data`**, published on **`localhost:18087`**.
3. Order a **phoenix** potion: `curl -X POST 'localhost:18087/orders?potion=phoenix'`
4. Remove the container (`docker rm -f eb-store`) and start a **new** `eb-store` the same way.
5. Order a **basilisk** potion, then check `curl localhost:18087/orders` lists both.
6. Submit from `~/elixirbytes`. The checks inspect the volume and the mount, read the orders file straight out of the volume with a throwaway container, and ask the shop over HTTP.
