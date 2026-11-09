# Your Own Registry

A registry is itself just a program, and of course there's an image for it: **`registry:2`**. Run one and you have a private image shelf on your own computer, with no account needed. Companies run the same thing on a server.

![On your laptop, docker tag and docker push send an image to cq-registry at localhost:15000 over its HTTP API /v2/. A server then uses docker pull and docker run. The image name localhost:15000/cq/webapp:1.0 includes the registry: the registry is part of the image name.](assets/registry.svg)

The registry listens on port 5000 inside its container. We'll use **15000** on your computer, because macOS often uses 5000 for AirPlay.

# Assignment

1. Run **`cq-registry`** from **`registry:2`**, detached, with a restart policy of `unless-stopped`, published on **127.0.0.1** port **15000** (the registry uses 5000 inside).
2. Check it's alive: `curl http://localhost:15000/v2/` should answer `{}`.
3. Run, then submit `dk.publish.local-registry`.
