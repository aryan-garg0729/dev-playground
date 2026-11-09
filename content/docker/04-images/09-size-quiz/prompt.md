# Why Is This Image So Big?

Look at two versions of nginx side by side:

```sh
docker images nginx
# REPOSITORY   TAG           SIZE
# nginx        1.27          192MB
# nginx        1.27-alpine   48MB
```

Same web server, same version, four times the size.

Many images on Docker Hub come in **variants**, shown in the tag:

| Tag ends in | Built on | Size |
| --- | --- | --- |
| (nothing) | Debian, a full Linux | Big |
| `-slim` | Debian with the extras removed | Medium |
| `-alpine` | Alpine, a Linux built to be tiny | Small |

Smaller images download faster, start faster and carry fewer programs that could have security holes. You'll make your own images tiny in the multi-stage chapter.

**Why is `nginx:1.27` so much bigger than `nginx:1.27-alpine`?**
