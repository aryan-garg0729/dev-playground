# Production Ready

Boss fight! Run the app the way you would on a real server, with every safety setting from this chapter.

# Assignment

Run **`cq-prod`** from **`cq/webapp:prod`**, detached, so that it:

1. is reachable only on **127.0.0.1**, port **8101**
2. stores its data in a volume named **`cq-prod-data`** at **`/data`**
3. has a read-only filesystem, with a tmpfs at `/tmp`
4. has all capabilities dropped and no new privileges
5. may use **128 MB** of memory and no extra swap
6. restarts after crashes and reboots, but not after you stop it on purpose

Wait until `docker ps` shows it as `(healthy)`, then submit `dk.production.production-boss`.
