# The Default Bridge

Each container gets its own network namespace: its own interfaces, IP address and ports. To talk to each other, containers have to share a Docker **network**.

```sh
docker network ls
# NETWORK ID     NAME      DRIVER    SCOPE
# 3c5e...        bridge    bridge    local    <- the default
# 9a1b...        host      host      local
# 7f2d...        none      null      local
```

Mo starts a database stand-in and tries to reach it from another container:

```sh
docker run -d --name eb-db alpine:3.20 sleep infinity
docker run --rm alpine:3.20 ping -c1 eb-db
# ping: bad address 'eb-db'
```

Both are running.

**Why doesn't the name resolve?**
