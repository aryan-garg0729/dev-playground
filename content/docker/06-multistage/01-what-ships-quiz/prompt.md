# What Are We Shipping?

Look at our release image next to the base it was built on:

```sh
docker images --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}' | grep -E 'golang|brewd'
# golang:1.23-alpine        250MB
# elixirbytes/brewd:1.5.1   260MB
```

brewd itself is about 8 MB. The security team asks why a potion API ships with a compiler, and what an attacker who gets a shell in the container would find.

Our Dockerfile has a single `FROM golang:1.23-alpine`, then copies the source and runs `go build`.

**Which of these are inside `elixirbytes/brewd:1.5.1`?** Pick every one that is.
