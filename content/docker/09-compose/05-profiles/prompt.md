# Tools on Demand

Dex wants a debugging container in the stack: something with a shell, sitting on the project network, for poking at services by name. But not in production, and not started by every `docker compose up`.

**Profiles** mark services as optional:

```yaml
services:
  web:
    image: demo/web            # no profiles: always part of the stack
  debugger:
    image: alpine:3.20
    profiles: [debug]          # only with --profile debug
```

```sh
docker compose up -d                        # web only
docker compose --profile debug up -d        # web and debugger
docker compose --profile debug exec debugger sh
```

## Assignment

1. Add a **`toolbox`** service to `compose.yaml`: image **`alpine:3.20`**, command **`sleep infinity`**, in the profile **`debug`**.
2. Start it with the profile, then try `docker compose --profile debug exec toolbox wget -qO- http://api:8080/healthz`.
3. Submit from `~/elixirbytes`. The checks list services with and without the profile and make that request from `toolbox`.
