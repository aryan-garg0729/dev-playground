# The Disk Is Full

A Docker host collects junk: stopped containers, old untagged images from rebuilds, build cache, forgotten volumes. Look before you delete:

```sh
docker system df
# TYPE            TOTAL     ACTIVE    SIZE      RECLAIMABLE
# Images          31        6         4.2GB     3.1GB (73%)
# Containers      18        5         12MB      9MB (75%)
# Local Volumes   7         3         310MB     120MB (38%)
# Build Cache     212       0         2.9GB     2.9GB
```

| Command | Removes |
| --- | --- |
| `docker container prune` | stopped containers |
| `docker image prune` | dangling images (untagged leftovers); `-a`: every image no container uses |
| `docker builder prune` | build cache |
| `docker volume prune` | anonymous volumes no container uses; `--all`: **named ones too** |
| `docker system prune` | stopped containers, unused networks, dangling images and build cache; `--volumes` adds volumes |

The staging shop's disk is at 97%. Its stack is stopped for maintenance, and the `ledger` volume holds the only copy of today's orders.

**What's the safe way to free space?**
