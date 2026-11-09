# Look, Don't Touch

A bind mount works both ways: the container can **change** your files too. For something like a config file, the container should only read it.

Add **`:ro`** (read-only) to the end of the mount:

```sh
-v "$(pwd)/config":/config:ro
```

Now any attempt to write there fails with `Read-only file system`. A bug (or an attacker) inside the container can't touch your files.

# Assignment

In `~/docker-course`:

1. Create a config folder with a settings file. Copy and paste:

   ```sh
   mkdir -p config
   echo "theme=dark" > config/settings.txt
   ```

2. Start **`cq-ro`** from **`alpine:3.20`**, detached, running `sleep infinity`, with your **`config`** folder mounted **read-only** at **`/config`**.
3. Try it: `docker exec cq-ro cat /config/settings.txt` works, and `docker exec cq-ro sh -c 'echo hacked > /config/settings.txt'` fails.
4. Run, then submit `dk.storage.read-only`.
