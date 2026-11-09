# Choosing the Final Base

Popular choices for the final stage of a Go service:

| Base | Size | Has a shell? | Has CA certificates? | Good for |
| --- | --- | --- | --- | --- |
| `scratch` | 0 MB, literally empty | no | no | static binaries that need nothing at all |
| `gcr.io/distroless/static` | about 2 MB | no | yes | static binaries that make HTTPS calls |
| `alpine:3.20` | about 8 MB | yes (busybox) | yes | when you need a shell or small tools |
| `debian:12-slim` | about 75 MB | yes | yes | programs that need glibc or apt packages |

Fewer tools means less for an attacker to use, and less for you when debugging.

Dex wants a **healthcheck** later: a command Docker runs *inside* the container to ask brewd if it's alive. brewd has no health flag, so it'll be an HTTP request to `localhost:8080/healthz`, made by a tool in the image.

**Which final base fits brewd today?**
