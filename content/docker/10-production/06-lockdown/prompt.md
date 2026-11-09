# Lock It Down

Running as `brewer` is step one. Three more `docker run` flags take away what's left:

| Flag | Takes away |
| --- | --- |
| `--read-only` | writing anywhere in the container's filesystem. Only mounts stay writable, so give it a volume for data and a `--tmpfs /tmp` for scratch files. |
| `--cap-drop ALL` | Linux **capabilities**: the pieces root's power is split into (binding ports below 1024, changing file owners, raw network access...). brewd listens on 8080 and needs none of them. |
| `--security-opt no-new-privileges` | gaining privileges later through setuid programs like `su` or `sudo`, even if one is in the image |

An attacker who gets code running in this container finds no writable binaries, no root, no capabilities, and no way to get them back.

## Assignment

1. Run **`eb-locked`** from `elixirbytes/brewd:prod`, detached, published on **`localhost:18093`**, with:
   - a read-only root filesystem and a tmpfs at `/tmp`
   - a named volume **`eb-locked-data`** at **`/data`**
   - all capabilities dropped
   - no new privileges
2. Order a potion to confirm brewd still works: `curl -X POST 'localhost:18093/orders?potion=test'`
3. Submit from `~/elixirbytes`. The checks read the security settings, try to write to `/etc`, and order a potion over HTTP.

> The new volume starts as a copy of the image's `/data`, ownership included, which is why `brewer` can write to it. A bind mount wouldn't get that.
