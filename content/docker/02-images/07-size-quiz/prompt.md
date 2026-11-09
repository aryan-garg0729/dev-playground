# Where Did the Disk Go?

Dex runs `docker images` after your tagging lesson:

```
REPOSITORY         TAG    IMAGE ID       SIZE
alpine             3.20   d9e853e87e55   13.7MB
elixirbytes/base   mine   d9e853e87e55   13.7MB
elixirbytes/base   copy   d9e853e87e55   13.7MB
```

Mo adds up the SIZE column, gets 41.1 MB and asks why a potion shop needs three copies of Linux.

**Roughly how much disk do these three rows actually use?**

> Worth trying yourself: `docker system df`, then `docker system df -v`.
