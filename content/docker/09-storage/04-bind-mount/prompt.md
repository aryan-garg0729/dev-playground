# Share a Folder

A volume is a box that Docker manages. Sometimes you want something else: **a folder on your own computer**, shared live with a container. That's a **bind mount**.

```sh
docker run -d -p 8091:80 -v "$(pwd)/site":/usr/share/nginx/html nginx:1.27-alpine
#                           └── a path on your computer      └── where it appears inside
```

- If the left side is a **path** (it starts with `/`), it's a bind mount. If it's a plain **name**, it's a volume.
- `$(pwd)` is replaced with the folder you're in, because Docker needs a full path.

Now the container reads files straight from your folder. Edit a file in your editor, refresh the browser, and the change is there. No rebuild, no restart. This is how people use Docker while they're developing.

# Assignment

In `~/docker-course`:

1. Start **`cq-live`** from **`nginx:1.27-alpine`**, detached, on port **8091**, with your **`site`** folder bind-mounted at **`/usr/share/nginx/html`**.
2. Open http://localhost:8091. It shows your `site/index.html` from the Dockerfiles chapter.
3. Edit the file on your computer. Copy and paste:

   ```sh
   echo '<h1>Edited on my computer</h1>' > site/index.html
   ```

4. Refresh. The change is live.
5. Run, then submit `dk.storage.bind-mount`.

Seeing `403 Forbidden` or an empty folder inside the container? On a Mac, Docker can only share folders it's allowed to see. Your home folder is shared by default; others are added in Docker Desktop under **Settings → Resources → File sharing**.
