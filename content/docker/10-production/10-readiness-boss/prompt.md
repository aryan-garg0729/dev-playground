# Production Readiness Review

Every ElixirBytes service passes a readiness review before its first real customer. Priya, Dex and Mo each brought their checklist:

> **Priya:** "It runs the image from **our registry**, tagged **`release`**, not something built on someone's laptop."
>
> **Dex:** "It tells us when it's sick, comes back after a crash or a reboot, but stays down if I stop it. It can't eat more than **64 MB**, swap included."
>
> **Security (via Mo):** "Not root. Nothing writable but its data. No capabilities. No privilege escalation."

## Assignment

In `~/elixirbytes`, with everything from this chapter:

1. Tag `elixirbytes/brewd:prod` as **`localhost:15000/elixirbytes/brewd:release`** and push it to `eb-registry`.
2. Delete your local `localhost:15000/elixirbytes/brewd:release` tag, so the next step really pulls from the registry.
3. Run **`eb-prod`** from **`localhost:15000/elixirbytes/brewd:release`**, detached, published on **`localhost:18094`**, satisfying every line of the checklist. Its orders go in a named volume **`eb-prod-data`** at `/data`.
4. Self-check: `docker ps` shows `(healthy)` after about 10 s, and an order over HTTP gets `201`.

Submit from `~/elixirbytes` to attack. The checks inspect `eb-prod`'s image, health, user, security settings, limits and restart policy, then place an order.

Every check is a hit on the final boss. Defeat it and brewd is ready to ship.
