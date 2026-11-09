# Pin It by Digest

Anywhere you can write `name:tag`, you can write `name@sha256:<digest>` instead:

```sh
docker run --rm alpine@sha256:d9e853e87e55526f6b2917df91a2115c36dd7c696a35be12163d44e6e2a4b6bc echo pinned
# pinned
```

The digest form is long and ugly, but it's a promise: these exact bytes, forever. Production deploys and security-sensitive base images are often pinned this way, and a tool updates the digests on purpose when the team decides to upgrade.

## Assignment

Priya wants proof you can run something by digest instead of by tag.

1. Find the digest of your local `alpine:3.20`. The previous quiz showed a template that prints `RepoDigests`.
2. Run a container named **`eb-pinned`** from **`alpine@sha256:<that digest>`** that prints the Alpine version with `cat /etc/alpine-release` and exits. Don't use `--rm`.
3. Submit from `~/elixirbytes`.

The checks read the digest themselves, then confirm `eb-pinned` was started from that digest (and not from the `alpine:3.20` tag), exited 0 and logged a `3.20.x` version.
