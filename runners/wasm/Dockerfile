# Wasm grading worker: the WASI modules plus cmd/worker in WORKER_RUNNER=wasm
# mode. Build context is the repo root (see Dockerfile.dockerignore).
FROM golang:1.27-alpine AS modules
WORKDIR /src
COPY runners/wasm/go.mod runners/wasm/go.sum ./
RUN go mod download
COPY runners/wasm/ ./
RUN GOOS=wasip1 GOARCH=wasm go build -trimpath -ldflags=-s -o /out/go.wasm ./cmd/go && \
    GOOS=wasip1 GOARCH=wasm go build -trimpath -ldflags=-s -o /out/js.wasm ./cmd/js

FROM golang:1.27-alpine AS worker
WORKDIR /src
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
RUN CGO_ENABLED=0 go build -trimpath -o /out/worker ./cmd/worker

# No Docker CLI, no shell tools needed at run time; runs as a non-root user in
# a read-only container (see docker-compose.yml).
FROM alpine:3.22
RUN adduser -D -u 10001 app
COPY --from=modules /out/ /opt/codequest/wasm/
COPY --from=worker /out/worker /usr/local/bin/
ENV WORKER_RUNNER=wasm WASM_DIR=/opt/codequest/wasm WASM_CACHE_DIR=/tmp/wasm-cache
USER 10001
ENTRYPOINT ["worker"]
