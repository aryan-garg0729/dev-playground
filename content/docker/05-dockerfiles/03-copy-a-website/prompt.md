# Put a Website in an Image

In the last chapter you changed nginx's page with `docker exec`, and the change vanished with the container. The proper way is to **bake the page into an image**, so every container made from it has the page from the start.

The `COPY` step copies files from your folder into the image:

```dockerfile
FROM nginx:1.27-alpine
COPY index.html /usr/share/nginx/html/index.html
#    └── from your folder   └── to this path inside the image
```

The image starts as plain nginx, then your `index.html` replaces the default page.

# Assignment

1. In `~/docker-course`, create the page and the Dockerfile. Copy and paste:

   ```sh
   mkdir -p site
   echo '<h1>Built into the image</h1>' > site/index.html
   cat > site/Dockerfile <<'EOF'
   FROM nginx:1.27-alpine
   COPY index.html /usr/share/nginx/html/index.html
   EOF
   ```

2. Build it as **`cq/site:v1`** from the `site` folder.
3. Run a container named **`cq-site`** from `cq/site:v1`, detached, on your port **8084** (nginx listens on 80 inside).
4. Open http://localhost:8084.
5. Run, then submit `dk.dockerfile.copy-a-website`.
