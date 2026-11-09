# Keep Secrets Out of Images

Passwords, API keys and tokens are **secrets**. An image is something you share: with your team, a registry, servers, sometimes the whole internet. So **a secret must never be built into an image.**

Every one of these leaks:

```dockerfile
# visible with docker image inspect, and inside every container
ENV PAYMENT_KEY=sk-live-123

# stays in a layer forever, even if you delete it later
COPY .env /app/.env

# recorded in docker history
ARG PAYMENT_KEY
```

Do this instead: build the image **without** the secret, and hand the secret to the **container** when it starts:

```sh
docker run -e PAYMENT_KEY="$PAYMENT_KEY" myapp                 # from your shell
docker run --env-file ./prod.env myapp                         # from a file you never commit
docker run -v "$(pwd)/key.txt":/run/secrets/key:ro myapp       # as a read-only file
```

(And add `.env` files to `.dockerignore` so they can't sneak into the context.)

A teammate built an image with `ENV PAYMENT_KEY=sk-live-123` and pushed it to the company registry.

**Who can read that key?**
