# Build, Settings and Volumes

Everything you've passed to `docker run` has a place in `compose.yaml`:

| `docker run` | `compose.yaml` |
| --- | --- |
| `docker build -t cq/webapp:compose .` | `build: .` and `image: cq/webapp:compose` |
| `-e NAME=value` | `environment:` |
| `-v notes:/notes` | `volumes:` (and the volume listed at the bottom) |
| `-p 8097:8080` | `ports:` |

With `build: .`, Compose builds the image from the Dockerfile in this folder before starting it. Rebuild after a change with `docker compose up -d --build`.

# Assignment

In `~/docker-course`:

1. From here on, Compose will build the app itself, so let's agree on one Dockerfile. It's the multi-stage one with an **Alpine** final stage: later lessons need the small tools Alpine includes. Copy and paste:

   ```sh
   cat > Dockerfile <<'EOF'
   FROM golang:1.23-alpine AS build
   WORKDIR /src
   COPY app/go.mod ./
   RUN go mod download
   COPY app/ .
   RUN --mount=type=cache,target=/root/.cache/go-build \
       CGO_ENABLED=0 go build -o /out/webapp .

   FROM alpine:3.20
   COPY --from=build /out/webapp /usr/local/bin/webapp
   ARG VERSION=dev
   ENV APP_VERSION=$VERSION
   LABEL org.opencontainers.image.version=$VERSION
   EXPOSE 8080
   ENTRYPOINT ["webapp"]
   EOF
   ```

2. Replace `compose.yaml`. Copy and paste:

   ```sh
   cat > compose.yaml <<'EOF'
   services:
     api:
       build: .
       image: cq/webapp:compose
       ports:
         - "8097:8080"
       environment:
         GREETING: Hello from Compose
         DATA_DIR: /notes
       volumes:
         - notes:/notes

   volumes:
     notes:
   EOF
   ```

3. Run `docker compose up -d --build`. Compose builds the image, then replaces the old `api` container.
4. Open http://localhost:8097.
5. Run, then submit `dk.compose.build-env-volumes`.
