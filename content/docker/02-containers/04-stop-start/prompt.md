# Stop and Start

A running container keeps using memory and CPU. When you don't need it, stop it:

```sh
docker stop cq-sleeper
# (wait about 10 seconds...)
# cq-sleeper
```

Yes, that wait is weird. Hold that thought, it's explained at the end of this chapter.

A stopped container isn't deleted. Bring it back with `docker start`:

```sh
docker start cq-sleeper
# cq-sleeper
docker ps
# ... Up 2 seconds   cq-sleeper
```

`docker start` restarts the **same** container, with the same name and settings. `docker run` always makes a **new** one. (`docker restart` does stop + start in one go.)

# Assignment

1. Stop **`cq-sleeper`** and confirm `docker ps -a` shows it as `Exited`.
2. Start it again and confirm `docker ps` shows `Up`.
3. Run, then submit `dk.containers.stop-start`.

The check confirms `cq-sleeper` is running *and* was stopped at some point.
