# Remove Volumes

Volumes outlive containers on purpose, so removing a container never removes its volume. When you really are done with the data, remove the volume yourself:

```sh
docker volume rm cq-demo-vol
```

Like images, Docker won't remove a volume that a container (even a stopped one) still uses:

```sh
docker volume rm cq-scrap
# Error response from daemon: remove cq-scrap: volume is in use - [a1b2c3...]
```

Remove the container first, then the volume.

A warning about one command you'll see online: `docker volume prune` deletes volumes that no container uses. That's often exactly what holds your data between runs. Read its list carefully before saying yes.

# Assignment

1. Set it up. Copy and paste:

   ```sh
   docker volume create cq-scrap
   docker run --name cq-scrapper -v cq-scrap:/data alpine:3.20 sh -c 'echo junk > /data/junk.txt'
   docker volume rm cq-scrap      # fails: in use
   ```

2. Remove **`cq-scrapper`**, then the **`cq-scrap`** volume.
3. Keep **`cq-notes`**!
4. Run, then submit `dk.storage.remove-volumes`.
