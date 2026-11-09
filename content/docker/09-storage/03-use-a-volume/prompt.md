# Data That Survives

Attach a volume to a container with `-v VOLUME:PATH`:

```sh
docker run -d -v cq-notes:/data cq/webapp:1.6
#             └── volume  └── where it appears inside the container
```

Everything the app writes to `/data` now goes into the volume, not into the container's scratch layer. Delete the container and the volume (and your data) stays.

# Assignment

In `~/docker-course`:

1. If `cq-temp` from the quiz is still around, remove it: `docker rm -f cq-temp`.
2. Start **`cq-notes-1`** from **`cq/webapp:1.6`**, detached, on port **8090**, with the **`cq-notes`** volume at **`/data`**.
3. Save a note:

   ```sh
   curl -X POST "http://localhost:8090/notes?text=buy+milk"
   ```

4. Remove `cq-notes-1` completely.
5. Start **`cq-notes-2`** the same way: same image, same port, same volume at `/data`.
6. Check `curl http://localhost:8090/notes`. The note survived!
7. Run, then submit `dk.storage.use-a-volume`.
