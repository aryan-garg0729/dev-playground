# Settings and Storage

Everything from earlier chapters has a Compose spelling:

```yaml
services:
  web:
    image: demo/web:dev
    environment:              # docker run -e
      LOG_LEVEL: debug
    volumes:                  # docker run -v
      - cache:/var/cache/web          # named volume (declared below)
      - ./config:/etc/web:ro          # bind mount, relative to compose.yaml
volumes:
  cache:                      # Compose creates it as <project>_cache
```

A named volume must also be listed under the top-level `volumes:` key. Compose prefixes it with the project name, so `cache` becomes `elixirbytes_cache`, and keeps it across `docker compose down` (unless you add `-v`).

`docker compose exec SERVICE COMMAND` runs a command in a service's container. Add **`-T`** when no terminal is attached, as in scripts or the CodeQuest checks.

## Assignment

brewd reads `DATA_DIR` to decide where `orders.log` goes.

1. Update the `api` service in `compose.yaml` so that:
   - **`DATA_DIR`** is set to **`/orders`**
   - a named volume **`ledger`** is mounted at **`/orders`**
2. `docker compose up -d` (Compose recreates the container because its config changed).
3. Order a **starlight**: `curl -X POST 'localhost:18091/orders?potion=starlight'`
4. Submit from `~/elixirbytes`. The checks read the environment, the mount and the order file through `docker compose exec`.
