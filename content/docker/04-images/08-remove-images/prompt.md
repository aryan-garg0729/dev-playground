# Remove Images

Images take up disk space, and they add up. Remove one with `docker rmi` (*remove image*):

```sh
docker rmi cq/demo:test
# Untagged: cq/demo:test
```

If other names still point at the same image, `rmi` only removes that one name ("Untagged"). The files are deleted once the last name is gone.

When it's an image's **last** name, Docker won't remove it while a container still uses it, **even a stopped one**:

```sh
docker rmi busybox:1.36
# Error response from daemon: conflict: unable to delete busybox:1.36 (must be forced)
# - container 8730be270b56 is using its referenced image 73aaf090f3d8
```

Remove that container first (`docker rm`), then the image.

Two cleanup commands you'll see online, with a warning: `docker container prune` deletes **every** stopped container on your computer, and `docker image prune` deletes leftover images with no name (shown as `<none>`). They're handy, but they don't care whether those things belong to this course, so read what they list before you say yes.

# Assignment

1. Remove the **`cq/base:copy`** tag. Keep `cq/base:mine`.
2. Set up the conflict, then fix it. Copy and paste:

   ```sh
   docker run --name cq-user busybox:1.36 echo "I use busybox"
   docker rmi busybox:1.36      # fails: cq-user is using it
   ```

3. Remove **`cq-user`**, then remove the **`busybox:1.36`** image.
4. Run, then submit `dk.images.remove-images`.
