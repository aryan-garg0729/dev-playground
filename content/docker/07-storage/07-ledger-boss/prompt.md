# The Ledger Must Survive

Priya's ticket:

> **DATA-9: no more lost orders.** Run the shop so that orders outlive any container, config is read-only, and scratch files never touch disk. Then prove it: take an order, **destroy the container**, start a fresh one, and show the order is still there.

## Assignment

In `~/elixirbytes`:

1. Run **`eb-ledger-shop`** from `elixirbytes/brewd:1.6.0`, detached, published on **`localhost:18097`**, with:
   - a named volume **`eb-ledger`** at **`/data`**
   - your `config/` folder bind-mounted **read-only** at **`/config`**
   - a **tmpfs** at **`/tmp`**, written with `--tmpfs`
2. Order an **elixir-of-life**: `curl -X POST 'localhost:18097/orders?potion=elixir-of-life'`
3. Wait a couple of seconds, then **remove `eb-ledger-shop`** and start it again with the same options.
4. Check `curl localhost:18097/orders` still lists it.

Submit from `~/elixirbytes` to attack. The checks inspect every mount, read the order straight out of the volume, compare the container's creation time with the order file's last change (the container must be newer: proof you recreated it), and ask the shop over HTTP.

Every check is a hit on the boss. Defeat it to unlock **Networking**.
