# Docker course: authoring notes (no story)

Beginner course, boot.dev-style: one idea per lesson, ~1 minute read, then do it.
No story, no characters. Docker only: the learner never reads or edits app code;
any file they need comes as one copy-paste block.

## Conventions (the checks find work by these names)
- Workspace: `~/docker-course` (run `codequest` from there).
- Containers, volumes, networks start with `cq-`. Images start with `cq/`.
- Base images: alpine:3.20, nginx:1.27-alpine, hello-world, golang:1.23-alpine.
- Host ports are unique per lesson, in the 8081-8099 range.

## The sample app (chapter 5 onwards): `app/main.go` + `app/go.mod`
Tiny Go web server, pasted with a heredoc. Docker compiles it; no Go on the host.
- `GET /` -> `<GREETING> from webapp <APP_VERSION>` (defaults Hello / dev)
- `GET /healthz` -> `ok`
- `GET /notes`, `POST /notes?text=...` -> `$DATA_DIR/notes.txt` (default /data)
- `webapp --version` -> `webapp <APP_VERSION>`
- Listens on `$PORT` (default 8080). On SIGTERM prints `webapp: shutting down`, exits 0.

## Checks
- Read-only `docker inspect -f` / `--format` with `key=value` output.
- Never check ids, digests, sizes or timestamps exactly; compare in shell.
- Learner runs builds/pulls; steps only inspect (30 s limit per command).
