# Tags Are Pointers

An image name like `alpine:3.20` has two parts:

```
alpine : 3.20
└─ repository   └─ tag
```

It's tempting to think of a tag as a copy of the image. It isn't. A tag is a **label stuck on an image**, the way a Git branch is a label on a commit. One image can have many tags, and adding a tag copies nothing:

```sh
docker tag alpine:3.20 elixirbytes/demo:latest
docker images
# REPOSITORY         TAG      IMAGE ID       SIZE
# alpine             3.20     d9e853e87e55   13.7MB
# elixirbytes/demo   latest   d9e853e87e55   13.7MB    <- same IMAGE ID
```

Same IMAGE ID, so it's the same image. `docker images` lists it twice, but it's stored once.

![Three tags, alpine:3.20, elixirbytes/base:mine and elixirbytes/base:copy, all point to one image, d9e853e87e55, which is stored once.](assets/tags.svg)

And `docker rmi elixirbytes/demo:latest` only peels the label off. The image stays on disk until its last tag (and every container using it) is gone.

## Assignment

ElixirBytes wants its own name for the base image, so we stop typing `alpine:3.20` everywhere.

1. Give `alpine:3.20` two more tags: **`elixirbytes/base:mine`** and **`elixirbytes/base:copy`**.
2. Run `docker images` and confirm all three rows share one IMAGE ID.
3. From `~/elixirbytes`, submit `dk.images.tags`.

The checks confirm both tags exist and count how many distinct image ids the three names point at.
