# depends_on Doesn't Wait

Mo adds the nginx front door to the stack:

```yaml
services:
  api:
    build: .
  front:
    image: nginx:1.27-alpine
    depends_on:
      - api
```

On a slow laptop, `front` sometimes logs `connect() failed (111: Connection refused) while connecting to upstream` for the first few requests.

`depends_on` also has a long form:

```yaml
    depends_on:
      api:
        condition: service_healthy     # or service_started, service_completed_successfully
```

**With the short form above, what does Compose wait for before starting `front`?**
