# The Whole Shop

Priya's ticket, the last one before launch:

> **LAUNCH-1: the whole shop in one `compose.yaml`.**
> - **api**: our image, orders in a named volume, healthchecked, **never published to the host**: customers only reach it through front.
> - **front**: nginx on port 18092, opens only once api is healthy.
> - Both come back on their own after a crash or a reboot, unless we stop them on purpose.
> - Prove the orders survive tearing the stack down.

## Assignment

In `~/elixirbytes`, update `compose.yaml` using only what this chapter covered:

1. **`api`**: keep the build, image, `DATA_DIR`, `ledger` volume and healthcheck, **remove the published port**, and set **`restart: unless-stopped`**.
2. **`front`**: as in the healthcheck lesson, plus **`restart: unless-stopped`**.
3. `docker compose up -d`, then order a **moonlight-mead** through the front door:
   `curl -X POST 'localhost:18092/orders?potion=moonlight-mead'`
4. Wait a couple of seconds, then **`docker compose down`** and **`docker compose up -d`**. Pick the `down` that keeps the ledger.
5. Self-check: `curl localhost:18092/orders`

Submit from `~/elixirbytes` to attack. The checks read the stack's state, the api container's ports and restart policy and the merged config, read the ledger, compare the container's creation time with the order's (the container must be newer: proof of the down and up), and ask the front door.

Every check is a hit on the boss. Defeat it to unlock **Production ready**.
