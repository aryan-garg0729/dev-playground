FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/api ./cmd/api && \
    CGO_ENABLED=0 go build -trimpath -o /out/cq ./cmd/cq && \
    CGO_ENABLED=0 go build -trimpath -o /out/worker ./cmd/worker

# Grading worker: needs the docker CLI and a Docker daemon (socket), no database.
FROM alpine:3.22 AS worker
RUN apk add --no-cache docker-cli
COPY --from=build /out/worker /usr/local/bin/
ENTRYPOINT ["worker"]

FROM alpine:3.22
RUN adduser -D -u 10001 app
COPY --from=build /out/api /out/cq /usr/local/bin/
USER app
EXPOSE 8080
ENTRYPOINT ["api"]
