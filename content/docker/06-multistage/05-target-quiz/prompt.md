# Which Stages Get Built?

Priya sketches a Dockerfile with four stages:

```dockerfile
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY app/ .
RUN go build -o /out/brewd .

FROM build AS test
RUN go test ./...                 # slow

FROM alpine:3.20 AS docs
RUN echo "generating docs..." && sleep 60

FROM alpine:3.20 AS final
COPY --from=build /out/brewd /usr/local/bin/brewd
CMD ["brewd"]
```

`docker build --target <stage>` stops at a named stage. Without `--target`, the target is the **last** stage.

**Running `docker build -t elixirbytes/brewd:x .` with BuildKit, which stages actually run?**
