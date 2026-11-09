# How Small Can You Go?

Your final stage can start from many bases. From biggest to smallest:

| Base | What's inside | Typical size |
| --- | --- | --- |
| `debian:12` / `ubuntu` | A full Linux with lots of tools | ~80-120 MB |
| `debian:12-slim` | Debian with extras removed | ~30-75 MB |
| `alpine:3.20` | A tiny Linux: shell, `apk`, basic tools | ~8 MB |
| `gcr.io/distroless/static` | Almost nothing: no shell, no package manager | ~2 MB |
| `scratch` | **Literally nothing.** An empty image | 0 MB |

The smaller the base, the fewer programs there are to attack, and the less you can do inside it when something goes wrong.

`scratch` only works for programs that need nothing else at all, which is why we built with `CGO_ENABLED=0`.

**Your team is containerising a brand-new app and expects to debug it inside the container now and then. Which final base is the best start?**
