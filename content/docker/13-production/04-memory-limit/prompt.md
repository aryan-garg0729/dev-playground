# Put a Lid on Memory

By default a container may use **all** of your computer's memory. One buggy app with a memory leak can slow down or crash everything else on the machine.

Set a limit:

```sh
docker run -m 32m --memory-swap 32m ...
#          │      └── memory + swap in total: equal to -m means "no swap"
#          └── at most 32 MB of memory
```

(Swap is disk space used as overflow memory. It's much slower, so a leaking app would just crawl instead of stopping.)

When the container goes over its limit, the system kills it immediately. You'll recognise the exit code: **137**, forced to stop. With a restart policy, it would come right back.

# Assignment

Make a memory leak on purpose, safely capped:

1. Run **`cq-hog`** from `alpine:3.20`, detached, with **32 MB** of memory and **no extra swap**, running this command. It keeps doubling a piece of text and prints its size. Copy and paste the command part:

   ```sh
   awk 'BEGIN { s = "leak"; while (1) { s = s s; print "holding " length(s) " bytes"; fflush() } }'
   ```

2. It dies within a second or two. Don't remove it. Look at `docker logs cq-hog`: where did it stop?
3. Run, then submit `dk.production.memory-limit`.

(`--cpus 0.5` limits CPU the same way: at most half of one processor.)
