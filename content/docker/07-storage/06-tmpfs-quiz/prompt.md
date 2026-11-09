# Scratch Space

The third kind of mount isn't storage at all:

```sh
docker run --tmpfs /tmp alpine:3.20 ...
docker run --mount type=tmpfs,dst=/tmp,tmpfs-size=16m alpine:3.20 ...
```

| Mount | Data lives | Survives the container? |
| --- | --- | --- |
| volume | Docker's storage area, on disk | yes |
| bind | a host folder, on disk | yes |
| tmpfs | **memory** | no: gone when the container stops |

A new brewd feature decrypts each order's payment token into a scratch file, uses it for a second and deletes it. Dex doesn't want those files touching a disk, ever.

**Where should brewd's scratch directory live?**
