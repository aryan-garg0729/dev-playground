# Don't Run as Root

Containers run as **root** unless the image says otherwise:

```sh
docker run --rm elixirbytes/brewd:health id
# uid=0(root) gid=0(root) groups=0(root),...
```

That root is the same user id 0 as the host's root. Namespaces and dropped capabilities fence it in, but every container escape bug ever found is far worse when the process inside was root. brewd needs no special powers, so it shouldn't have any.

```dockerfile
RUN adduser -D -H -u 10001 brewer        # alpine: -D no password, -H no home folder
USER brewer                              # everything after, and the container, runs as brewer
```

Two catches:

- Anything the program writes to needs to be owned by that user. Create and `chown` those folders **before** `USER`, while the build is still root.
- Use a numeric UID that's unlikely to clash with a real user on the host, like `10001`: bind-mounted files show up on the host under that number.

## Assignment

1. In the **final stage** of your `Dockerfile`, before `CMD`:
   - create a user **`brewer`** with UID **`10001`**
   - create **`/data`** and make `brewer` its owner
   - switch to `brewer`
2. Keep the healthcheck. Build **`elixirbytes/brewd:prod`**.
3. Submit from `~/elixirbytes`. The checks run `id`, try to write to `/data` (allowed) and to `/usr/local/bin` (not allowed), and read the healthcheck.
