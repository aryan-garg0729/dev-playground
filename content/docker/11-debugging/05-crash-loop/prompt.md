# It Exits Straight Away

The most common Docker support question: "I ran my container and it's not in `docker ps`!"

The checklist is always the same:

1. **`docker ps -a`**: is it there, and what's the exit code?
2. **`docker logs NAME`**: what did it say before it died? Nine times out of ten, the answer is right there.
3. Fix the cause, **remove** the broken container and run it again.

You can't change a container's settings (like `-e` or `-p`) after it's created. A fix almost always means `docker rm` and a new `docker run`.

# Assignment

1. Start this container. Copy and paste:

   ```sh
   docker run -d --name cq-broken alpine:3.20 sh -c '[ -n "$API_KEY" ] || { echo "error: API_KEY is not set"; exit 2; }; echo "started with key"; sleep infinity'
   ```

2. It's not in `docker ps`. Work through the checklist and find out why.
3. Fix it: remove `cq-broken` and run it again (same name, same command) with the setting it needs. Any value works.
4. Run, then submit `dk.debug.crash-loop`.
