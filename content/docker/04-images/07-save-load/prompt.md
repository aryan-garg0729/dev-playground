# Images as Files

Usually images travel through a registry like Docker Hub. But you can also pack an image into a single file and carry it on a USB stick, or to a computer with no internet.

```sh
docker save -o demo.tar cq/demo:test     # image -> file
docker load -i demo.tar                  # file -> image
# Loaded image: cq/demo:test
```

The `.tar` file holds every layer plus the image's name, so `docker load` brings it back exactly as it was.

# Assignment

1. In `~/docker-course`, save **`cq/base:mine`** to a file named **`base.tar`**.
2. Delete the image's tag from Docker: `docker rmi cq/base:mine`. (`docker images` no longer lists `cq/base:mine`.)
3. Load it back from `base.tar`.
4. Run, then submit `dk.images.save-load`.
