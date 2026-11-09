# Volumes

The fix for disappearing data is a **volume**: a storage space that Docker manages **outside** any container. Containers come and go; the volume stays.

Think of a container as a hotel room and a volume as a safe-deposit box at the front desk. You check out of the room (remove the container), but your things in the box are still there for your next stay.

Volumes have their own commands:

```sh
docker volume create cq-demo-vol      # make one
docker volume ls                      # list them
docker volume inspect cq-demo-vol     # details
```

`inspect` shows a `Mountpoint`: where Docker keeps the files. On a Mac that path is inside Docker's hidden Linux VM, so you can't open it in Finder. That's fine, because you'll always reach a volume's files through a container.

# Assignment

1. Create a volume named **`cq-notes`**.
2. Find it with `docker volume ls` and look at it with `docker volume inspect cq-notes`.
3. Run, then submit `dk.storage.volumes`.
