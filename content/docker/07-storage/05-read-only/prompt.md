# Look, Don't Touch

Config files are for a container to **read**, not edit. If the container is compromised, or just buggy, a writable config mount lets it change files on your host. Add `readonly` (or `:ro` with `-v`):

```sh
docker run --mount type=bind,src="$PWD/config",dst=/config,readonly ...
docker run -v "$PWD/config":/config:ro ...
```

Writes then fail inside the container with `Read-only file system`, while you can still edit the files on your side.

## Assignment

1. In `~/elixirbytes`, create **`config/motd.txt`** containing `welcome, brewer`.
2. Run **`eb-ro`** from **`alpine:3.20`**, detached, running `sleep infinity`, with `~/elixirbytes/config` bind-mounted **read-only** at **`/config`**.
3. Try `docker exec eb-ro sh -c 'echo hi > /config/motd.txt'` and watch it fail.
4. Submit from `~/elixirbytes`. The checks inspect the mount, read the file through the container, and try to overwrite it.
