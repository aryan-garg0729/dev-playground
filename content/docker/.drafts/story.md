# Docker course: story bible (ElixirBytes)

## World
You've joined **ElixirBytes**, a small online potion shop. Its only backend is
**brewd**, a Go HTTP API (standard library only) that takes potion orders. Today it runs
because someone SSHed into a box and started it in `tmux`. Your job: containerise it,
shrink it, persist it, network it, compose it and make it production-ready.

## Characters
- **Priya**: platform lead, writes the tickets. Precise, dry.
- **Dex**: on-call SRE, has been paged too many times. Grumpy, right.
- **Mo**: product manager, asks "can we just..." questions.

## Naming conventions (the grader finds work by these names)
- Containers, volumes, networks: `eb-...` (e.g. `eb-hello`, `eb-api`, `eb-net`, `eb-orders`).
- Images: `elixirbytes/...` (e.g. `elixirbytes/brewd:v1`).
- Workspace folder: `~/elixirbytes` (one folder for the whole course). App code in `~/elixirbytes/app`.
- brewd listens on **8080** inside the container. Host ports: `18080` + lesson-specific.

## brewd endpoints (given in chapter 3)
- `GET /healthz` -> `ok`
- `GET /version` -> `brewd <version>`
- `GET /orders` / `POST /orders` (stored in `$DATA_DIR/orders.log`, default `/data`)
- `--version` flag prints `brewd <version>`
- Handles SIGTERM: prints `brewd: shutting down` and exits 0.

## Check style
- Read-only `docker ... --format` / `inspect -f` with `key=value` labels, e.g.
  `docker inspect -f 'status={{.State.Status}} exit={{.State.ExitCode}}' eb-crash`.
- Never check ids, digests, timestamps or sizes exactly; use shell comparisons that print
  `size=ok` / `size=too-big`.
- Learner runs the slow stuff (builds, pulls) themselves; steps only inspect (30 s limit).

## Bosses
Bosses are `project_cli` lessons (decision 2026-11-08): each one is a realistic ticket that
combines the chapter's skills on the learner's machine, 5-8 tests (one HP each). Where a
value differs per learner (random exit codes, digests, layer counts), steps capture it from
Docker and compare it with what the learner wrote down.

## Workspace state per chapter
- ch1: no files. Containers eb-hello (running, /etc/eb-badge), eb-crash (exited 42),
  eb-sleepy (137), eb-polite (143), eb-ticker (exited 0).
- ch2: ~/elixirbytes/answers.txt (key=value lines), alpine.tar; tags elixirbytes/base:mine|copy.
- ch3: ~/elixirbytes/app/{main.go,go.mod,testdata/fixtures.bin}, Dockerfile, .dockerignore;
  images elixirbytes/brewd:v1, :lean; container eb-web on 18080.
- ch4: Dockerfile (cache-friendly, ARG BUILD_VERSION, cache mount); build.log; images :v4, :cached.
- ch5: app/entrypoint.sh (exec), Dockerfile.ch5; images :banner, :ch5, :graceful.
- ch6: Dockerfile (multi-stage, alpine final); images :slim, :scratch, :1.6.0.
- ch7: volumes eb-orders, eb-ledger; orders/ (bind), config/ (read-only).
- ch8: networks eb-net, eb-private; front/default.conf; containers eb-api, eb-front.
- ch9: compose.yaml (api, front, toolbox profile); project name elixirbytes.
- ch10: registry eb-registry on 127.0.0.1:15000; image localhost:15000/elixirbytes/brewd.
