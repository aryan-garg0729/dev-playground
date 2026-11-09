# I Changed It, Nothing Changed

Here's a classic beginner moment. You edit a file, refresh the browser... and nothing changes.

Remember the chain:

```
Dockerfile + files  --build-->  image  --run-->  container
```

`cq-site` is a container made from the image `cq/site:v1`, which was built **before** your edit. The image is frozen, so the container serves the old page. Editing files in your folder changes nothing until you:

1. **build** a new image, and
2. **replace** the container with one made from the new image.

Giving the new build a new tag (`v2`) keeps the old one around in case you need to go back.

# Assignment

1. Change the page. Copy and paste:

   ```sh
   echo '<h1>Version 2</h1>' > site/index.html
   ```

2. Refresh http://localhost:8084. Still the old page, right?
3. Build the `site` folder again as **`cq/site:v2`**.
4. Replace the container: remove `cq-site`, then start a new **`cq-site`** from **`cq/site:v2`**, detached, on port **8084**.
5. Refresh: **Version 2**.
6. Run, then submit `dk.dockerfile.rebuild`.
