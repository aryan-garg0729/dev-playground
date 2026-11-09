# Ship It Small

Priya's ticket:

> **REL-160: brewd 1.6.0, small edition.** Security wants the toolchain and the source out of production. Ops wants deploys to pull less. Keep the healthcheck option open (we'll need a shell tool there), keep the version stamping, keep the fast builds.
>
> Budget: **under 20 MB, at most 2 filesystem layers.**

## Assignment

In `~/elixirbytes`:

1. Your `Dockerfile` must be a multi-stage build with:
   - a build stage that downloads dependencies first, stamps `BUILD_VERSION` into the binary and uses the compiler cache mount
   - a final stage on the base from this chapter's quiz, holding only the binary at `/usr/local/bin/brewd` and the `org.opencontainers.image.version` label, starting with `CMD ["brewd"]`
2. Build **`elixirbytes/brewd:1.6.0`** with `BUILD_VERSION=1.6.0`.
3. Run it detached as **`eb-web6`**, reachable at **`http://localhost:18086`**.

> Layers: count the `RootFS.Layers`, not the history rows. A `LABEL` adds no filesystem layer.

Submit from `~/elixirbytes` to attack. The checks measure size and layers, look for a Go toolchain or `/src`, read the label and call `/version` over HTTP.

Every check is a hit on the boss. Defeat it to unlock **Storage**.
