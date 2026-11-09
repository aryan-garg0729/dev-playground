# Tags Are Nicknames

A tag is just a **nickname** pointing at an image. One image can have many nicknames, and adding one copies nothing.

Add a nickname with `docker tag`:

```sh
docker tag alpine:3.20 cq/demo:test
#          └── existing   └── new name
```

Now list them:

```sh
docker images
# REPOSITORY   TAG    IMAGE ID       SIZE
# alpine       3.20   0a9a5dfd008f   7.8MB
# cq/demo      test   0a9a5dfd008f   7.8MB
```

Same **IMAGE ID**, so it's the same image under two names. Only the label is new, and that 7.8 MB is on your disk once.

Teams use this to give one image several useful names: `myapp:1.4.2`, `myapp:1.4` and `myapp:stable` can all point at the same build.

# Assignment

1. Tag `alpine:3.20` as **`cq/base:mine`**.
2. Tag `cq/base:mine` as **`cq/base:copy`**. (You can tag a tag.)
3. Run `docker images` and check all three share one IMAGE ID.
4. Run, then submit `dk.images.tags`.
