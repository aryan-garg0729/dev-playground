# The exec Fix

The fix is one word. `exec` tells the shell: **"replace yourself with this program."** No shell is left in the middle, so the webapp itself becomes PID 1 and hears `docker stop` directly.

```sh
#!/bin/sh
echo "starting webapp..."
exec webapp "$@"
```

Any start script that ends by launching the main app should launch it with `exec`.

(Another option you'll see: `docker run --init` adds a tiny helper as PID 1 that passes signals along. The `exec` fix is better because it needs no extra flag.)

# Assignment

In `~/docker-course`:

1. Change the last line of `app/start.sh` to `exec webapp "$@"`. Copy and paste:

   ```sh
   cat > app/start.sh <<'EOF'
   #!/bin/sh
   echo "starting webapp..."
   exec webapp "$@"
   EOF
   ```

2. Build as **`cq/webapp:graceful`** (same Dockerfile as last lesson).
3. Run **`cq-fast`** from it, detached. Then `time docker stop cq-fast`: well under a second.
4. Check its logs show both the banner and `webapp: shutting down`, and its status is `Exited (0)`.
5. Run, then submit `dk.entrypoint.exec-fix`.
