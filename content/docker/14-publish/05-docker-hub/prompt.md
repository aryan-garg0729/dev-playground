# Publish to Docker Hub

Time to put an image on the internet, where anyone can `docker pull` it.

On Docker Hub, your images live under **your username**: `USERNAME/name:tag`. Docker Hub is the default registry, so there's no address in front.

```sh
docker login                                    # once: sign in from the terminal
docker tag cq/webapp:prod USERNAME/cq-webapp:1.0
docker push USERNAME/cq-webapp:1.0
```

The first push creates the repository, and it's **public**: anyone can pull it. That's fine for this tiny demo app. Never push an image with secrets in it (remember the production chapter).

# Assignment

1. Create a free account at [hub.docker.com](https://hub.docker.com) if you don't have one.
2. Sign in from your terminal with `docker login`. (Docker Desktop users may already be signed in.)
3. Tag **`cq/webapp:prod`** as **`USERNAME/cq-webapp:1.0`**, using your own Docker Hub username, and push it.
4. Open `https://hub.docker.com/r/USERNAME/cq-webapp` in your browser. There it is!
5. Add your username to `answers.txt` as `hub_user=USERNAME`.
6. Run, then submit `dk.publish.docker-hub`.

The check asks Docker Hub for your image, the same way anyone in the world could. You can delete the repository on Docker Hub afterwards if you like.
