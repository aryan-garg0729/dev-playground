# The Ten-Second Stop

Dex's complaint: *"Every deploy takes ten seconds per container, and every old container dies with exit 137. Something's eating our SIGTERMs."*

Someone added a wrapper script to print a startup banner. Before we fix anything, we reproduce the bug. Debugging rule one: make it happen on purpose.

## Assignment

1. Create **`~/elixirbytes/app/entrypoint.sh`**, exactly as it is in production:

   ```sh
   #!/bin/sh
   echo "brewd starting on $(hostname)"
   brewd "$@"
   ```

   and make it executable: `chmod +x app/entrypoint.sh`. (`"$@"` passes along any arguments the script was given.)

2. Create **`Dockerfile.ch5`**: your chapter 4 Dockerfile (dependencies first, version stamp), plus:
   - copy the script to **`/usr/local/bin/entrypoint.sh`**
   - **`ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]`** in place of the `CMD`

3. Build it as **`elixirbytes/brewd:banner`**, run it detached as **`eb-slow`**, then time the stop:

   ```sh
   time docker stop eb-slow
   ```

   About ten seconds, then exit 137. Bug reproduced.

4. Submit from `~/elixirbytes`. The checks read the image's entrypoint, `eb-slow`'s exit code and its logs, which have the banner but no `brewd: shutting down` line: brewd never got the message.
