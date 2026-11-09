# Who Is PID 1?

`docker top` lists a container's processes as the host sees them. Here's `eb-slow` before you stopped it:

```
$ docker top eb-slow -o pid,args
PID       COMMAND
1828964   /bin/sh /usr/local/bin/entrypoint.sh
1828980   brewd
```

Inside the container, the first row is PID 1:

```sh
docker exec eb-slow cat /proc/1/comm
# entrypoint.sh
```

Two facts you need:

1. `docker stop` sends SIGTERM to **PID 1 only**.
2. Linux treats PID 1 specially: a signal it has no handler for is **ignored**, not fatal.

**Why did `docker stop eb-slow` take ten seconds and end in exit 137?**
