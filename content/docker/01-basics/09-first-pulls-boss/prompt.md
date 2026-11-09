# Image Explorer

Boss time! No new ideas here, just everything from this chapter together: pull images, run commands inside them and read what they tell you.

You can run **any command** inside an image by putting it after the image name:

```sh
docker run --rm alpine:3.20 ls /
# bin  dev  etc  home  lib  ...
```

That `ls` ran *inside* a tiny Alpine Linux, not on your computer. `--rm` cleans up the container afterwards.

# Assignment

In `~/docker-course`:

1. Pull a third image, **`busybox:1.36`** (a toolbox of small Linux commands). Make sure `alpine:3.20` and `nginx:1.27-alpine` are still there too.
2. Find the exact Alpine version inside `alpine:3.20`. It's stored in the file `/etc/alpine-release`:

   ```sh
   docker run --rm alpine:3.20 cat /etc/alpine-release
   ```

3. Find the BusyBox version. Running `busybox` with no arguments prints it on the first line (`BusyBox v1.36.1 ...`):

   ```sh
   docker run --rm busybox:1.36 busybox | head -n 1
   ```

4. Create **`answers.txt`** with exactly these two lines, filling in the versions you found (numbers only, like `1.36.1`, no `v`):

   ```
   alpine=VERSION
   busybox=VERSION
   ```

   For example: `echo "alpine=3.20.3" >> answers.txt`.

5. Submit `dk.basics.first-pulls-boss`.

The checks run the same two commands and compare their answers with yours.
