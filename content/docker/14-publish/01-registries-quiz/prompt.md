# Where Images Live

So far every image you built lives only on your computer. To run it on a server, or give it to a teammate, it has to travel. Images travel through a **registry**: an online shelf of images that you **push** to and **pull** from.

You've been pulling from one since chapter 1: **Docker Hub**. There are many others:

| Registry | Address starts with | Often used by |
| --- | --- | --- |
| Docker Hub | (nothing) or `docker.io/` | Everyone; public images |
| GitHub Container Registry | `ghcr.io/` | Projects on GitHub |
| Amazon, Google, Azure registries | `*.amazonaws.com/`, `*.pkg.dev/`, `*.azurecr.io/` | Companies on those clouds |
| Your own | `localhost:15000/`, `registry.mycompany.com/` | Private, in-house images |

Remember image names? **The registry is part of the image's name.** `ghcr.io/acme/shop:2.1` says where it lives. That's how `docker push` knows where to send it.

The usual way to deploy:

```
build on your laptop (or a CI server) --push--> registry --pull--> production server
```

**You built an image on your laptop. What's the normal way to run that exact image on a server?**
