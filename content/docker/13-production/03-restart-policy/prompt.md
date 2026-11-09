# Try, Try Again

Let's watch a restart policy at work on a container that always crashes.

`on-failure:3` means "restart it when it fails, but give up after 3 retries". You don't want something that's truly broken restarting forever.

# Assignment

1. Run **`cq-flaky`** from `alpine:3.20`, detached, with the restart policy **on-failure, at most 3 retries**, running this command. Copy and paste the command part:

   ```sh
   sh -c 'echo "crashed at $(date +%T)"; exit 1'
   ```

2. Wait about 10 seconds. Docker starts it four times in total (once, plus three retries), then leaves it exited.
3. Look at `docker logs cq-flaky`: one line per attempt. And `docker inspect -f '{{.RestartCount}}' cq-flaky`.
4. Run, then submit `dk.production.restart-policy`.
