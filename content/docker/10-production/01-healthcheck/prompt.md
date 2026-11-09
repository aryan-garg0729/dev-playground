# Is It Actually Up?

`running` only means the process exists. A brewd stuck in a deadlock is `running`. So is one whose disk is full. Dex wants Docker itself to know whether brewd is **working**.

In chapter 9 you put the healthcheck in `compose.yaml`. It can also live in the image, so every `docker run` gets it for free:

```dockerfile
HEALTHCHECK --interval=5s --timeout=2s --retries=3 \
  CMD wget -qO- http://localhost:8080/healthz || exit 1
```

`docker ps` then shows `(health: starting)`, `(healthy)` or `(unhealthy)`, and `docker inspect` keeps the last few results:

```sh
docker inspect -f '{{json .State.Health}}' eb-x
# {"Status":"healthy","FailingStreak":0,"Log":[{"ExitCode":0,"Output":"ok\n",...}]}
```

Docker itself only *reports* health. Compose's `depends_on`, orchestrators and your monitoring are what act on it.

## Assignment

1. Add a **`HEALTHCHECK`** to the **final stage** of your `Dockerfile`: every **5 s**, fetch `http://localhost:8080/healthz` with `wget`.
2. Build it as **`elixirbytes/brewd:health`**.
3. Run **`eb-health`** from it, detached. Within about 10 s, `docker ps` should say `(healthy)`.
4. Submit from `~/elixirbytes`. The checks read the image's healthcheck and `eb-health`'s current health.
