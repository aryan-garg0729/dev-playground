# What depends_on Waits For

Compose starts services in parallel. When one service needs another, say so with `depends_on`:

```yaml
services:
  web:
    image: my-web
    depends_on:
      - db
  db:
    image: postgres:16-alpine
```

Now Compose starts `db` before `web`.

But "started" and "ready" aren't the same thing. A database container starts in a fraction of a second, then spends a few more seconds getting ready to accept connections. A web app that connects straight away gets `connection refused`.

**With the short form above, how long does Compose wait before starting `web`?**
