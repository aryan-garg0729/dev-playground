# One-service production image for free hosting (Render; build context: repo root).
# The API serves the built web app on the same origin (WEB_DIR), and every
# start applies migrations and compiles content/ into Postgres before serving.
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

FROM alpine:3.22
RUN adduser -D -u 10001 app
COPY --from=go /out/api /out/cq /usr/local/bin/
COPY --from=web /web/dist /srv/web
COPY content /srv/content
ENV WEB_DIR=/srv/web
USER app
EXPOSE 10000
ENTRYPOINT ["/bin/sh", "-c", "cq migrate up && cq build /srv/content && exec api"]
