# Docker Hub

Where did `hello-world` come from? From **[Docker Hub](https://hub.docker.com)**, a free online library of ready-made images. Want a database, a web server or a plain Linux to play with? Somebody has already packed it, and you can download it with one command.

Downloading an image is called **pulling**:

```sh
docker pull alpine:3.20
# 3.20: Pulling from library/alpine
# Status: Downloaded newer image for alpine:3.20
```

The part after the colon (`3.20`) is the **tag**, which usually means the version. Always pick a version on purpose; you'll see why in the images chapter.

To see every image on your computer:

```sh
docker images
# REPOSITORY    TAG       IMAGE ID       CREATED       SIZE
# alpine        3.20      0a9a5dfd008f   2 weeks ago   7.8MB
# hello-world   latest    d2c94e258dcb   1 year ago    13kB
```

Meet the two images you'll use most in this course:

- **alpine**: a tiny Linux, about 8 MB. Perfect for trying things.
- **nginx**: a popular web server, the program that sends web pages to browsers.

Browse Docker Hub sometime. Look for the **Docker Official Image** badge: those images are maintained with Docker's help, so they're a safe first choice.

# Assignment

1. In `~/docker-course`, pull both images:

   ```sh
   docker pull alpine:3.20
   docker pull nginx:1.27-alpine
   ```

2. Run `docker images` and find them in the list.
3. Run, then submit `dk.basics.docker-hub`.
