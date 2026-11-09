# Give It Less

By default a container can use all of the host's memory, so one leak can starve everything else on the machine. Set a limit:

```sh
docker run -m 64m --memory-swap 64m ...    # 64 MB of RAM, and no swap on top
docker stats --no-stream                   # live usage against each limit
```

`--memory-swap` is memory plus swap, so setting it equal to `-m` means "no swap": the container fails fast instead of crawling. The container can read its own limit:

```sh
docker run --rm -m 32m alpine:3.20 cat /sys/fs/cgroup/memory.max
# 33554432
```

Over the limit, the kernel's **out-of-memory (OOM) killer** ends the process with SIGKILL: you know that exit code. Docker also records `State.OOMKilled`, but on some setups (certain cgroup v2 and VM combinations) it reads `false` even after a real OOM kill. Exit 137, plus a limit, plus logs showing memory climbing, is the evidence that holds up.

## Assignment

Reproduce a memory leak on purpose, safely capped:

1. Run **`eb-hog`** from `alpine:3.20`, detached, with **32 MB** of memory and **no extra swap**, running this leak. It doubles a string forever and logs its size each time:

   ```sh
   awk 'BEGIN { s = "potion"; while (1) { s = s s; print "cauldron holds " length(s) " bytes"; fflush() } }'
   ```

2. It dies within a second or two. Don't remove it. Look at `docker logs eb-hog`: where did it stop?
3. Submit from `~/elixirbytes`. The checks read the limits, the exit code and the logs.
