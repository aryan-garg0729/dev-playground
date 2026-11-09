# Lock It Down

A few more `docker run` options take away things the app never needs:

| Option | What it does |
| --- | --- |
| `--read-only` | The container's own files can't be changed at all |
| `--tmpfs /tmp` | ...except a small in-memory `/tmp`, because many programs need one |
| `-v cq-locked-data:/data` | ...and a volume for the app's real data |
| `--cap-drop ALL` | Removes special powers (like changing file owners or network settings) the app doesn't use |
| `--security-opt no-new-privileges` | Stops the app from gaining more powers later |

Each one is a door an attacker can no longer use. Together with a non-root user, they're the standard checklist for running containers safely.

# Assignment

1. Run **`cq-locked`** from **`cq/webapp:prod`**, detached, published on **127.0.0.1** port **8100**, with:
   - a read-only filesystem and a tmpfs at `/tmp`
   - a new volume **`cq-locked-data`** at **`/data`**
   - all capabilities dropped
   - no new privileges
2. Save a note to prove it still works: `curl -X POST "http://127.0.0.1:8100/notes?text=hello"`
3. Run, then submit `dk.production.lockdown`.

A new, empty volume starts with a copy of what the image has at that path, owner included. That's why the `app` user can write to it. (A bind mount doesn't get that treatment.)
