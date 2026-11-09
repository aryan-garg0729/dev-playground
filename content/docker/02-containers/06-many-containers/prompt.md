# Many Containers

Remember the recipe and the cakes? One image can become as many containers as you like, all running at once and each in its own box. They don't share files and they don't get in each other's way.

This is a big part of why Docker is popular: you can run three different databases on one laptop and they never fight.

With lots of containers, `docker ps` gets crowded. Two tricks:

```sh
docker ps --filter name=cq-            # only names containing "cq-"
docker ps --format '{{.Names}} {{.Status}}'
# cq-one Up 4 seconds
# cq-two Up 3 seconds
```

`--format` picks the columns you want. Those `{{ }}` blocks are placeholders, and you'll see them a lot in this course.

# Assignment

1. Start three detached containers from **`alpine:3.20`**, named **`cq-one`**, **`cq-two`** and **`cq-three`**, each running `sleep infinity`.
2. List just those three with `docker ps --filter name=cq-`.
3. Run, then submit `dk.containers.many-containers`.
