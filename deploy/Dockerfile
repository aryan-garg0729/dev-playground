# One-service production image for free hosting (Render; build context: repo root).
# The API serves the built web app on the same origin (WEB_DIR), and every
# start applies migrations and compiles content/ into Postgres before serving.
# It also serves the learner CLI (/cli/, /install.sh), built for APP_URL.
# Lessons that run learner code (Python, Go, JS, SQL) need the grading workers
# and the sandbox, which this image leaves out; the current courses are quiz
# and CLI lessons, graded inside the API.

FROM node:26-alpine AS web
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY web/ ./
RUN npm run build

FROM golang:1.27-alpine AS go
WORKDIR /src
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
RUN CGO_ENABLED=0 go build -trimpath -o /out/api ./cmd/api && \
    CGO_ENABLED=0 go build -trimpath -o /out/cq ./cmd/cq

# The learner CLI, served by this site at /cli/ with /install.sh. It is built
# for this site's address: APP_URL (Render passes service env vars to Docker
# builds as build args; locally: docker build --build-arg APP_URL=...).
ARG APP_URL
ARG RENDER_GIT_COMMIT=dev
COPY deploy/install.sh /tmp/install.sh
RUN site="${APP_URL%/}"; \
    case "$site" in https://*|http://localhost*|http://127.0.0.1*) ;; \
      *) echo "APP_URL build arg must be the site's https URL, got '$APP_URL'" >&2; exit 1 ;; esac; \
    mkdir -p /out/cli && \
    for t in darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64; do \
      os="${t%/*}"; arch="${t#*/}"; ext=""; [ "$os" = windows ] && ext=".exe"; \
      CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath \
        -ldflags "-s -w -X main.version=${RENDER_GIT_COMMIT} -X main.defaultAPIURL=${site} -X main.defaultAppURL=${site}" \
        -o "/out/cli/codequest-$os-$arch$ext" ./cmd/codequest || exit 1; \
    done && \
    sed "s#__SITE__#${site}#g" /tmp/install.sh > /out/install.sh

FROM alpine:3.22
RUN adduser -D -u 10001 app
COPY --from=go /out/api /out/cq /usr/local/bin/
COPY --from=web /web/dist /srv/web
COPY --from=go /out/cli /srv/web/cli
COPY --from=go /out/install.sh /srv/web/install.sh
COPY content /srv/content
ENV WEB_DIR=/srv/web
USER app
EXPOSE 10000
ENTRYPOINT ["/bin/sh", "-c", "cq migrate up && cq build /srv/content && exec api"]
