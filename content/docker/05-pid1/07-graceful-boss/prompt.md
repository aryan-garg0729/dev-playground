# Zero-Downtime Deploys

Dex's ticket, written at the end of a long week:

> **OPS-77: deploys must stop old containers in under 3 seconds.** No more 137s in the deploy log. Keep the startup banner (support relies on it). Ship it as brewd **1.5.1** so I can tell the fixed build from the broken one.

## Assignment

In `~/elixirbytes`, using what you've learned so far:

1. Build **`elixirbytes/brewd:graceful`** with `BUILD_VERSION=1.5.1`. The image must:
   - start through `/usr/local/bin/entrypoint.sh`, set as the `ENTRYPOINT` in exec form
   - print the `brewd starting` banner, then make **brewd PID 1**
   - pass any `docker run` arguments through to brewd
2. Self-check before you attack:

   ```sh
   docker run -d --name eb-try elixirbytes/brewd:graceful
   docker exec eb-try cat /proc/1/comm      # brewd
   time docker stop eb-try                  # well under 3 s
   docker logs eb-try                       # banner ... brewd: shutting down
   docker rm eb-try
   ```

Submit from `~/elixirbytes` to attack. The checks read your image's entrypoint, then start their own throwaway container, **`eb-probe-stop`**, check its PID 1, time its stop and read its logs, and finally run the image with `--version`. The probe container is removed afterwards.

Every check is a hit on the boss. Defeat it to unlock **Multi-stage builds**.
