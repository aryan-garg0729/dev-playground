# Everything Docker Knows

You've used `docker inspect -f` a few times. Without `-f`, it prints **everything** Docker knows about a container as one big JSON document: its settings, mounts, networks, state, exit code, start time...

```sh
docker inspect cq-api | less        # scroll with the arrow keys, q to quit
docker inspect cq-api | grep -i ipaddress
```

When you know what you want, `-f` picks one value out. The `{{ }}` path follows the JSON's structure:

```sh
docker inspect -f '{{.State.Status}}' cq-api
# running
docker inspect -f '{{.Config.Image}}' cq-api
# cq/webapp:1.6
```

Network names with a dash need `index` to look them up:

```sh
docker inspect -f '{{(index .NetworkSettings.Networks "cq-net").IPAddress}}' cq-api
# 172.19.0.2
```

That's the container's address on `cq-net`. It changes whenever the container is recreated, which is exactly why you call containers by name instead.

# Assignment

1. Run `docker inspect cq-api` and skim the output. Find the `State`, `Mounts` and `NetworkSettings` sections.
2. Add two lines to `answers.txt`:
   - `api_ip=` followed by `cq-api`'s address on `cq-net`
   - `api_image=` followed by the image `cq-api` runs
3. Run, then submit `dk.debug.inspect`.
