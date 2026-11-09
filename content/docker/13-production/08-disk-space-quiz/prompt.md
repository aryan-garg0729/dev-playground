# Where Did My Disk Space Go?

After a few weeks of Docker, your disk fills up. Every build, every pulled image and every stopped container takes space. See where it went:

```sh
docker system df
# TYPE            TOTAL   ACTIVE   SIZE      RECLAIMABLE
# Images          31      6        4.1GB     3.2GB (78%)
# Containers      24      9        12MB      3MB (25%)
# Local Volumes   7       3        300MB     120MB (40%)
# Build Cache     188     0        2.6GB     2.6GB
```

The cleanup commands, from gentle to dangerous:

| Command | Removes |
| --- | --- |
| `docker system prune` | Stopped containers, unused networks, nameless (`<none>`) images, leftover build cache |
| `docker builder prune` | The build cache (next builds are slower, nothing else changes) |
| `docker image prune -a` | Every image no container uses (you'll have to pull or build them again) |
| `docker system prune --volumes` | All of the above plus **volumes no container is using right now**: your data! |

All of them list what they'll remove and ask first. Read that list.

Your disk is nearly full. You have a Compose stack that's stopped for the weekend, and its notes volume must survive.

**What's the safe way to free space?**
