# Ship Without a Registry

Usually images travel through a **registry** like Docker Hub: you `push` from one machine and `pull` on another. But ElixirBytes has a brewing server on a locked-down network with no internet access. You can still move an image there, as a file:

```sh
docker save -o thing.tar some/image:tag   # image -> tarball (every layer + metadata)
# ...copy thing.tar over with scp or a USB stick...
docker load -i thing.tar                  # tarball -> image, tags included
```

Peek inside one and you'll see the layers from two lessons ago, stored by digest:

```sh
tar -tf thing.tar
# blobs/sha256/1db32a79...
# blobs/sha256/3f26bc2d...
# index.json
# manifest.json
# oci-layout
```

(`docker export` sounds similar but is different: it flattens one *container's* filesystem into a single tar with no layers or image settings. You'll rarely want it.)

## Assignment

1. From `~/elixirbytes`, save **`elixirbytes/base:mine`** to **`base.tar`**.
2. Delete the tag with `docker rmi elixirbytes/base:mine`. `docker images` no longer lists it.
3. Load it back from `base.tar`. It should reappear with its tag.
4. Submit from `~/elixirbytes`. The checks list the archive, read its manifest and confirm the tag is back.
