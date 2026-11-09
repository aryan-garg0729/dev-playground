# Read the Logs

A container running in the background with `-d` can't print to your screen. So where does its output go?

Docker saves everything a container prints. These are its **logs**:

```sh
docker logs cq-demo
```

Two handy options:

```sh
docker logs -f cq-demo        # follow: keep showing new lines live (Ctrl+C to quit)
docker logs --tail 5 cq-demo  # only the last 5 lines
```

When something breaks, the logs are the first place to look.

# Assignment

1. Start this container. It counts to five, one number per second, then says goodbye. Copy and paste it:

   ```sh
   docker run -d --name cq-counter alpine:3.20 sh -c 'for i in 1 2 3 4 5; do echo "count $i"; sleep 1; done; echo "all done"'
   ```

2. Straight away, watch it live with `docker logs -f cq-counter`. Press Ctrl+C once you see `all done`.
3. Run `docker logs cq-counter` again. The logs are still there even though the container has finished.
4. Run, then submit `dk.containers.logs`.
