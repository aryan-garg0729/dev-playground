# Optional Services

Some services you only want now and then: a debugging toolbox, a database admin page, a test runner. Put them in a **profile** and they're skipped unless you ask for them:

```yaml
  toolbox:
    image: alpine:3.20
    command: ["sleep", "infinity"]
    profiles: ["debug"]
```

`command:` replaces the image's `CMD`, just like typing a command after the image name in `docker run`.

```sh
docker compose up -d                    # api and front only
docker compose --profile debug up -d    # api, front and toolbox
```

# Assignment

In `~/docker-course`:

1. Add the `toolbox` service above to `compose.yaml`, under `services:` (the same indentation as `api` and `front`).
2. Start it with the `debug` profile.
3. Use it: `docker compose --profile debug exec toolbox wget -qO- http://api:8080/healthz`.
4. Run, then submit `dk.compose.profiles`.
