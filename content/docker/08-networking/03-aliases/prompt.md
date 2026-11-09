# Another Name for It

Other services shouldn't care that the orders API happens to run in a container named `eb-api`. Next month it might be `eb-api-blue` during a deploy. A **network alias** is an extra DNS name for a container on one network:

```sh
docker run -d --name eb-x --network eb-demo --network-alias cache alpine:3.20 sleep infinity
# on eb-demo, both "eb-x" and "cache" resolve to this container
```

For a container that's already running, swap its connection for one with an alias:

```sh
docker network disconnect eb-demo eb-x
docker network connect --alias cache eb-demo eb-x
```

Aliases belong to a network, so the same container can have different names on different networks.

## Assignment

1. Give **`eb-api`** the alias **`orders`** on **`eb-net`**. Recreating it with `--network-alias` or reconnecting it both work.
2. Check it: `docker run --rm --network eb-net alpine:3.20 wget -qO- http://orders:8080/healthz`
3. Submit from `~/elixirbytes`. The checks read `eb-api`'s aliases, then reach it by both names.
