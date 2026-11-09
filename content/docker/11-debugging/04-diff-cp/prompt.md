# What Changed, and Copy It Out

Two more tools for poking at a container from the outside.

**`docker diff`** lists every file a container has changed compared with its image:

```sh
docker diff cq-x
# C /etc            <- Changed
# A /etc/my.conf    <- Added
# D /tmp/old.txt    <- Deleted
```

Great for answering "what did this app actually write?".

**`docker cp`** copies files between a container and your computer, in either direction:

```sh
docker cp cq-x:/path/in/container ./here     # out of the container
docker cp ./file.txt cq-x:/tmp/               # into the container
```

It works even on stopped containers, so you can rescue a crash report from a container that won't start.

# Assignment

1. Start a container that writes a report. Copy and paste:

   ```sh
   docker run -d --name cq-writer alpine:3.20 sh -c 'mkdir -p /reports && echo "token=$(date +%s)" > /reports/report.txt && sleep infinity'
   ```

2. Find what it wrote with `docker diff cq-writer`.
3. Copy `/reports/report.txt` out of `cq-writer` into `~/docker-course` as **`report.txt`**.
4. Run, then submit `dk.debug.diff-cp`.
