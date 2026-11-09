# CodeQuest Build Spec — React + Go

## How to use this doc

This is the source of truth for building CodeQuest (working name), a gamified, hands-on coding-learning platform in the style of boot.dev, with a React frontend and a Go backend. `CLAUDE.md` at the repo root points here.

Prompt Claude Code one phase at a time, for example: "Implement Phase 1 from docs/SPEC.md. Stop when its acceptance criteria pass and summarize what you built."

How to write courses, chapters and lessons is in `docs/CONTENT.md`. Design references live in `docs/design/`: `tokens.json` (all design tokens), `README.md` (the design system's rules), `components/` (component guidelines and HTML previews) and `screens/` (a PNG of every designed screen).

## Product in one paragraph

Learners work through courses made of short lessons. Each lesson is an exercise: read a concept, write code in a browser editor, run it, submit it. Submissions are graded automatically against tests. Passing earns XP, keeps a daily streak alive, moves the learner up a weekly league, and unlocks the next lesson. An AI tutor gives hints, never full solutions. The first chapters of every course are free; finishing a course requires a paid subscription.

Non-goals for v1: mobile apps, user-generated courses, video content, live classes, certificates.

## Tech stack and repo layout

Use a single monorepo with a Go API, a React SPA, a separate sandbox runner page, and content as files in Git.

| Layer | Choice | Notes |
| --- | --- | --- |
| Frontend | React 18 + TypeScript + Vite | SPA; React Router; TanStack Query for server state; Zustand for small client state |
| Styling | Tailwind CSS | Keep components in `web/src/components` |
| Code editor | Monaco (`@monaco-editor/react`) | Lazy-load; it is large |
| Browser execution | Pyodide (Python), sql.js (SQLite), native JS | Runs in a Web Worker inside a sandboxed iframe |
| Backend | Go 1.23+, `net/http` + chi router | Plain handlers, no heavy framework |
| DB | PostgreSQL 16 | Migrations with `golang-migrate`; queries with `sqlc` |
| Cache / leaderboards | Redis 7 | Sorted sets for leagues; rate limiting |
| Server-side execution | Docker containers (later gVisor) for Python; WASI modules under wazero for Go and JS (Phase 10) | Worker pools behind Redis job queues |
| Auth | Google and GitHub OAuth (decision 2026-11-09: email magic links removed, no SMTP) | Session cookie, `HttpOnly`, `SameSite=Lax` |
| Payments | Stripe (global), Razorpay (India) | Behind one `billing` interface |
| AI tutor | LLM API via a Go `tutor` package | Provider behind an interface |
| Local dev | Docker Compose | Postgres, Redis, API, web, runner |

```
codequest/
  CLAUDE.md  docker-compose.yml  Makefile
  docs/
    SPEC.md                 # this file
    design/                 # tokens.json, README.md, components/, screens/
  api/                      # Go module
    cmd/api/main.go         # HTTP server
    cmd/worker/main.go      # server-side grading worker (Redis queue -> Docker; no DB access)
    cmd/cq/main.go          # CLI: validate, test, build content
    cmd/codequest/          # learner CLI: login, run and submit project_cli lessons (Phase 8)
    internal/
      content/              # load + validate exercise packages
      exercise/             # type registry + plugins (code, mcq, sql, projectcli)
      checks/               # reusable check primitives
      runner/               # Runner interface + docker impl
      grading/              # submission pipeline, result storage
      game/                 # XP, levels, streaks, leagues, medals
      economy/              # Bytes, shop, inventory, power-ups, quests
      tutor/                # AI hints
      billing/              # Stripe, Razorpay, entitlements
      auth/  store/  httpapi/
    db/migrations/  db/queries/
  web/                      # React app (app.<domain>)
    tailwind.config.ts      # generated from docs/design/tokens.json
    src/styles/tokens.css   # generated from docs/design/tokens.json
    src/exercise-types/     # one folder per type: renderer components
    src/runner-client/      # postMessage bridge to sandbox
  runners/python/           # codequest/python:3.12 grading image + harness (Phase 2)
  sandbox/                  # static runner page (served from separate origin)
    index.html  src/main.ts src/worker.ts   # Pyodide / sql.js execution
    src/protocol.ts         # zod message schemas (identical copy in web/src/runner-client)
  content/                  # courses as files (see Content format); paths.yaml lists career paths
```

## Architecture overview

The browser runs code for instant feedback inside an isolated sandbox; every submission is re-graded on the server before any XP is awarded.

```
 React app (lessons, editor, XP, streaks)  <--postMessage-->  Sandbox iframe (separate origin, Web Worker:
        |                                     Run only,          Pyodide, sql.js; no cookies)
        | Submit                              ungraded
        v
 Go API (auth, content, grading, XP, streaks, leagues, economy, billing, tutor)
   |            |                 |                     |
 Postgres     Redis            Grading worker         External APIs
 (users,      (job queue,      (Docker runners,       (Stripe, Razorpay,
  progress,    leagues,         no network)            LLM provider)
  ledgers,     rate limits)
  content)
   ^
 Content repo (YAML, prompts, tests) loaded by `cq build`
```

The worker pulls grading jobs from Redis and runs them in locked-down Docker containers. Course content is compiled from Git into Postgres by `cq build`.

## Content format

Courses are plain files in `content/`, versioned in Git. Platform code never contains course content, and authors never touch platform code.

```
content/python-basics/
  course.yaml
  03-loops/
    chapter.yaml
    01-sum-list/
      exercise.yaml
      prompt.md
      starter.py
      solution.py          # reference solution, server-only
      tests_visible.py     # shown to learner
      tests_hidden.py      # server-only
      wrong/               # known-bad solutions; CI asserts each fails
        off_by_one.py
```

Every `exercise.yaml` has a common header plus a type-specific `spec`:

```yaml
id: py.loops.sum-list      # globally unique, stable forever
type: code                  # selects the plugin
version: 3                  # bump when grading behaviour changes
title: Sum a list
xp: 20
difficulty: 2               # 1-5
tags: [loops, lists]
free: true                  # in the free tier
boss: false                 # true for chapter bosses (see Boss fights)
spec:
  language: python
  runner: pyodide           # pyodide | docker-python (python); wasm (go, js)
  starter: starter.py
  entry: main.py
  limits: { timeout_ms: 3000, memory_mb: 128 }
  examples:                 # shown as "Example n: Input / Output / Explanation"
    - { input: "nums = [1, 2, 3]", output: "6", explanation: "1 + 2 + 3 = 6" }
    - { input: "nums = []", output: "0", explanation: "an empty list sums to 0" }
  constraints:
    - "0 ≤ len(nums) ≤ 10⁴"
    - "-1000 ≤ nums[i] ≤ 1000"
  checks:
    - kind: unit_tests
      file: tests_visible.py
      visible: true
    - kind: unit_tests
      file: tests_hidden.py
      weight: 2
    - kind: forbidden_calls
      names: [sum]
      message: Try it with a loop instead of sum()
```

```yaml
id: py.loops.range-quiz
type: mcq
version: 1
xp: 5
spec:
  question: prompt.md
  multiple: false
  options:
    - { id: a, text: "0 to 4" }
    - { id: b, text: "1 to 5" }
  answer: [a]
  explain: range(5) starts at 0 and stops before 5.
```

Rules:

- `id` never changes, even if the folder moves. Progress and submissions key on `id`.
- Lesson images (Phase 9): an exercise folder may hold a flat `assets/` folder of images (`png`, `jpg`/`jpeg`, `webp`, `gif`, `svg`; the bytes must match the extension), referenced from its markdown as `![alt](assets/loop.svg)` and from spec strings such as an mcq option's `image: assets/a.png`. Everything in `assets/` is public, so `cq validate` fails on: a reference to a file that is not there, any external image (`http:`, `//`, `data:`), a file no public markdown or spec string uses (`solution.md` does not count), files over 1 MB or over 5 MB per exercise, subfolders and symlinks, server-only names (`solution.*`, `tests_hidden.*`, `wrong`), and SVGs with script, event handlers, `foreignObject`, entities, or `href`/`url()` other than `#fragment`. Alt text is required in practice: it is what shows if an image cannot load.
- Order comes from `course.yaml` and `chapter.yaml` lists, not from folder names.
- The build step (`cq build`) validates everything and writes a compiled bundle to Postgres. The API never parses YAML at request time.
- Files marked server-only (`solution.*`, `tests_hidden.*`, `wrong/`) are never included in any API response.

Manifests (decided in Phase 1; unknown fields are a validation error everywhere):

```yaml
# course.yaml
id: langchain-basics        # lowercase slug, unique
title: LangChain Basics     # required
description: ...
order: 1                    # position among courses (then id)
published: true             # default true
chapters: [01-chat-models, 02-prompt-templates]   # folder names, in order
icon: lc                    # 1-4 lowercase letters or digits for the catalog tile; default: first 2 of the id
level: intermediate         # beginner (default) | intermediate | advanced
tags: [Python, AI]          # catalog filters, at most 4, each 1-20 characters
topics: Chat models, prompts  # one line for the course card, at most 60 characters

# chapter.yaml
id: lc.chat-models          # lowercase words joined by . or -, globally unique
title: Chat models & messages
exercises: [01-message-roles, 02-runnable-methods] # folder names, in order
```

- Career paths (decided 2026-11-06, `screens/courses.png`): an optional `paths.yaml` at the root of a folder of courses lists paths in order: `paths: [{id (slug), title, description, goal, courses: [{course, note}]}]`. `goal` is an onboarding goal (`zero`, `backend`, `interviews`, `ai`), at most one path per goal; it makes the path "Your path" for learners who picked that goal. Every course must exist in the tree and be published; `note` (≤ 60 characters) is the line under the course on the path card. `cq build` upserts `career_paths` and unpublishes paths no longer listed. The catalog holds only real courses; "coming soon" stubs, course prerequisites and estimated hours from the design are not built.
- Every folder must be listed, and every listed folder must exist with its manifest; both are errors otherwise. Hidden folders (`.x`) are ignored; folders inside an exercise (e.g. `wrong/`) belong to that exercise.
- Exercise header: `id`, `type`, `version` (≥ 1), `title` and `xp` (≥ 0) are required; `difficulty` defaults to 1 (1–2 shows as Easy), `tags` to `[]`, `free` and `boss` to false. Only `code` and `project_cli` exercises can be bosses.
- `cq validate <dir>` and `cq build <dir>` print every problem as `path: reason` and exit 1. `<dir>` is a folder of courses or a single course folder.
- `cq build` writes everything in one transaction, upserting by id. Courses, chapters and exercises no longer in the content are **unpublished, never deleted**, so progress and submissions keyed on old ids survive.
- mcq: `question` names a markdown file; options may carry `image: assets/<file>` (Phase 9; the option's `text` is still required and is its alt text); `answer` and `explain` live only in `private_spec`. `explain` is returned after every submit (pass or fail) as the `Feedback` of the implicit `exact_choice` check.
- Run location (decided in Phase 3): `code` and `sql` specs take `run: browser | server | off`. `browser` runs in the sandbox iframe (Pyodide for Python, sql.js for SQL); `server` sends the answer to `POST /exercises/{id}/run`, which runs it on the grading worker with visible checks only, ungraded and not stored; `off` hides Run. Defaults: `pyodide` → browser, `docker-python` → server, `sql` → browser. Changing it is a manifest edit plus `cq build`; content CI checks the solution's visible checks wherever `run:` points, so a lesson moved to `browser` whose packages don't load there fails CI.
- code (decided in Phase 2): `language: python`; `runner: docker-python` (submit-only) or `pyodide` (server grades it on `docker-python`); `prompt` (markdown), `entry` (a module file name such as `main.py`), `starter` and `solution` are required files; `limits` default to 3000 ms / 128 MB, capped at 10000 ms / 512 MB; `examples`, `constraints` and `hints` are optional and public. At least one check. `solution.*` and `wrong/*.py` are read by `cq build`/`cq test` but never stored in `private_spec`; hidden test sources are stored there (server-only). Visible test files are sent to the browser as `visible_tests`.
- code in other languages (decided in Phase 10): `language: go` or `js` with `runner: wasm`, `run: server` (default) or `off`; `entry` is `main.go` / `main.js` style, and `wrong/` holds files of that extension. Checks are `stdout_equals`, `stdout_regex`, `io_cases` (at most one per exercise) and `exit_code`; `unit_tests` and `forbidden_calls` are Python-only. Go is interpreted (yaegi) with a fixed standard-library allowlist, so `go build` rules such as "declared and not used" are not enforced. JS is ES2015+ (goja) with Node-like `console`, `process.exit`, `process.stdout.write` and `require("fs").readFileSync(0)` for stdin; no modules or npm. See `docs/CONTENT.md`.
- Test files are Python modules with top-level `test_*` functions that import from the entry module. A test's docstring is its row name. `from cq import expect` gives `expect(got, expected)`, which compares exact types (so a learner object with a rigged `__eq__` cannot pass) and reports got/expected. Label hidden tests in words: a hidden label that looks like code (contains `()[]{}=`) is replaced by "Hidden test N". A hidden test may declare its input on a docstring line `input: <call>` (e.g. `input: fizzbuzz(15)[14]`); the harness parent reads it from the AST, the API keeps it server-side (`submissions.hidden_inputs`, never serialized in a `Result`) and only Peek reveals it, for one failed check. The API cuts it at `==`, `->` or `→`, so never write the expected value there. `cq.fakes` (fake chat models) is available to tests and learner code in the image.

## Exercise type plugin system

Each exercise type is a plugin registered by name in Go and mirrored by a React renderer. Adding a type means adding one Go package, one React component, and one spec schema; nothing else changes.

Go interface (`api/internal/exercise/types.go`):

```go
type Type interface {
    Name() string                                     // "code", "mcq", "sql"
    ParseSpec(raw yaml.Node, dir fs.FS) (Spec, error) // validate + load files (cq build only)
    LoadSpec(private json.RawMessage) (Spec, error)   // rebuild Spec from exercises.private_spec (API)
    PublicView(ex Exercise) (json.RawMessage, error)  // what the browser may see
    Grade(ctx context.Context, ex Exercise, sub Submission) (Result, error)
}

// Optional: types with reference answers implement Testable for `cq test`.
type Testable interface {
    ReferenceAnswer(s Spec) (json.RawMessage, error)               // must pass
    WrongAnswers(s Spec) (map[string]json.RawMessage, error)       // each must not pass
}
```

Runner-based types get the server runner from the context (`runner.FromContext`), set by the grading service (queue runner) or `cq test` (Docker runner), so plugins hold no global state.

```go

var registry = map[string]Type{}

func Register(t Type) { registry[t.Name()] = t }   // called from each plugin's init()
func Get(name string) (Type, bool) { t, ok := registry[name]; return t, ok }
```

React registry (`web/src/exercise-types/index.ts`):

```ts
export interface ExerciseRendererProps<V = unknown, A = unknown> {
  header: React.ReactNode;                   // pills + title from the page; renderer places it in its problem pane
  publicView: V;                             // from GET /exercises/:id
  draft: A | undefined;                      // autosaved answer
  onChange: (answer: A) => void;
  onRun?: (answer: A) => Promise<RunOutput>; // client-side run, optional
  onSubmit: (answer: A) => Promise<GradeResult>;
  result: GradeResult | undefined;
}

export const renderers: Record<string, React.ComponentType<ExerciseRendererProps<any, any>>> = {
  code: CodeExercise,
  mcq: McqExercise,
  sql: SqlExercise,
};
```

v1 types:

| Type | Answer shape | Default runner |
| --- | --- | --- |
| `code` | `{ files: { [path]: string } }` | `pyodide` or `js` |
| `mcq` | `{ choices: string[] }` | none |
| `sql` | `{ query: string }` | `sqljs` (graded on the worker) |
| `project_cli` | `{ base_url?, steps: [{ cli: { exit_code, stdout } } \| { http: { status, headers, body, error? } } \| { skipped: true }] }` | the learner's machine, via the `codequest` CLI (Phase 8); graded in-process |

Later types (do not build in v1): `fill_blank`, `predict_output`, `order_steps`, `open_answer` (AI-graded).

`project_cli` (Phase 8, modelled on boot.dev's CLI lessons). Manifest `spec`: `prompt` (markdown file), `solution` (`solution.json`, a recorded passing answer; `wrong/*.json` are recorded failing ones, both for `cq test` only), optional `base_url` (the default for `${baseURL}`; the learner may override it with `codequest config base_url`), optional `allowed_os` (`linux`, `darwin`, `windows`), and 1–20 `steps`. Each step has a `description` and exactly one of:

- `cli_command`: `command` (run with `sh -c`, or `cmd /C` on Windows, in the learner's current folder; 30 s limit), optional `captures`, and `tests`, each one of `exit_code`, `stdout_contains_all`, `stdout_contains_none`, `stdout_lines_gt`. "stdout" is stdout and stderr combined.
- `http_request`: `request` (`method` GET/POST/PUT/PATCH/DELETE/HEAD, `url` starting with `${baseURL}` or http(s)://, `headers`, `body_json`, `follow_redirects`; 10 s limit), optional `captures`, and `tests`, each one of `status_code`, `body_contains`, `json_body: {path, op: eq|ne|gt|lt|contains, value}`, `header: {name, contains}`.

A capture (`name` plus `regex`, first submatch, or `json_path`, HTTP only) defines `${name}` for the step's own tests and later steps; `${var}` is interpolated in commands, URLs, header values, body strings and expected values. Paths are a small jq subset (`.a.b[0]`, `.`), not jq. Every test is visible: the public view is the whole spec plus `lesson_id`, minus `solution`/`wrong/`, and its JSON key for tests is `visible_tests`. Output is capped at 16 KB per step and headers at 4 KB. `Grade` makes one visible check row per test ("Step 2: status code is 200"). It **recomputes captures from the reported output** and ignores anything else the client sends; a step whose input variable was not captured fails with "Needs ${x}…". `Result.stdout` is a transcript of the steps.

Prefer configuration over new types. "Fix the bug" is a `code` exercise with buggy starter code. "Make it faster" is `code` plus a `max_runtime` check. Only add a type when the UI is genuinely different.

## Grading

The server is the only authority on whether a submission passed. The browser runs code for instant feedback; the server re-grades every submission before awarding XP.

### Two paths

1. **Run (client, ungraded).** The learner clicks Run. The sandbox executes their code plus visible tests and shows output. Nothing is stored.
2. **Submit (server, authoritative).** The browser POSTs the answer. The API calls the type's `Grade`, which executes all checks (visible and hidden) on a server runner, stores the result, and returns it. Only this path awards XP.

Client runners map to server runners: `pyodide` → `docker-python`, `sqljs` → SQLite in the `docker-python` grading container (decided in Phase 3, instead of in-process `modernc.org/sqlite`, so untrusted SQL never runs inside the API process). `mcq` grades in-process with no runner. The same test files must produce the same result on both. `wasm` lessons (Phase 10) have no browser runner yet: Run goes to the server, which runs the same module Submit uses.

### Contracts

```go
type Submission struct {
    ExerciseID string          `json:"exercise_id"`
    Version    int             `json:"version"`
    Answer     json.RawMessage `json:"answer"`   // type-specific shape
}

type Result struct {
    Status    string        `json:"status"`    // pass | fail | error | timeout
    Score     float64       `json:"score"`     // 0..1, weighted by checks
    Checks    []CheckResult `json:"checks"`
    Stdout    string        `json:"stdout"`    // truncated to 16 KB
    RuntimeMs int           `json:"runtime_ms"`
}

type CheckResult struct {
    Name     string `json:"name"`
    Passed   bool   `json:"passed"`
    Visible  bool   `json:"visible"`
    Feedback string `json:"feedback,omitempty"` // hidden checks: message only, never test source
    Got      string `json:"got,omitempty"`      // failed visible checks only: short repr
    Expected string `json:"expected,omitempty"` // failed visible checks only
    Line     int    `json:"line,omitempty"`     // learner source line the failure points at
    Weight   float64 `json:"-"`                 // score weight, 0 = 1; never sent
}
```

`status` is `pass` only when every check passes; `score` is the weighted passing fraction. Hidden test rows never carry got/expected, exception messages or assertion internals: only an authored assertion message, "Wrong result", or "Your code raised <ExceptionType>". Game, tutor and analytics code must depend only on `Result`, never on the exercise type.

### Check primitives (`api/internal/checks`)

A grader runs the `checks` list from the manifest. Each kind is a small function registered by name, reusable by any type.

| Kind | What it does | Used by |
| --- | --- | --- |
| `unit_tests` | Runs a test file against learner code; one CheckResult per test | code |
| `stdout_equals` | Exact match on program output (trim trailing whitespace) | code |
| `stdout_regex` | Regex match on output | code |
| `io_cases` | Runs the program once per stdin; output must equal each case's stdout; one row per case | code (wasm) |
| `exit_code` | The program's exit code equals `expected` | code |
| `forbidden_calls` | AST/static scan for banned names | code |
| `max_runtime` | Fails if runtime exceeds N ms | code |
| `exact_choice` | Selected option ids equal the answer set | mcq |
| `sql_result_equals` | Compares the result set to `expected` (visible) or to the solution's on the check's dataset (hidden); `ordered` flag | sql |

```go
type Check interface {
    Kind() string
    Run(ctx context.Context, in CheckInput) ([]CheckResult, error)
}
// Optional, called by cq build: reject bad configs (unknown fields are errors).
type Validator interface { Validate(config json.RawMessage) error }
```

Manifest fields (every code check also takes `name`, `visible`, `weight ≥ 0`):

| Kind | Fields | Defaults |
| --- | --- | --- |
| `unit_tests` | `file` (a `.py` in the exercise folder) | `visible: false`; one row per test, weight applies per row; `name` not allowed |
| `stdout_equals` | `expected` | `visible: true`, name "Program output"; runs the entry as `__main__`; trailing whitespace per line and trailing blank lines ignored |
| `stdout_regex` | `pattern` (Go RE2, unanchored), `message` | `visible: true` |
| `io_cases` | `cases: [{stdin, stdout, name?, visible?}]`, 1-50 | `visible: true` (per case overridable), rows "<name> N" (name default "Case"); compared like `stdout_equals`; hidden rows show no got/expected and their stdin is never sent on Run; `runner: wasm` only, one per exercise |
| `exit_code` | `expected` (0-125) | `visible: true`, name "Exit code"; runs the entry as a program with empty stdin |
| `sql_result_equals` | `ordered`, `expected: {columns, rows}` (required when visible), `setup` (a `.sql` file; only on hidden checks) | `visible: true`; values only (column count must match, names may differ); numbers compare with a small tolerance; one statement; `ATTACH`/`DETACH`/`PRAGMA` denied; 1000-row cap |
| `forbidden_calls` | `names` (identifiers), `message` | static AST scan of learner files: `f(...)` and `x.f(...)` count. A teaching aid, not a sandbox |

### Runner interface (`api/internal/runner`)

```go
type Runner interface {
    Run(ctx context.Context, job Job) (Output, error)
}
type Job struct {
    Image  string // e.g. codequest/python:3.12
    Stdin  []byte // the harness job (files, tests, limits) as JSON
    Limits Limits // timeout, memory, pids, output bytes
}
type Output struct { Stdout, Stderr string; ExitCode, RuntimeMs int; TimedOut bool }
```

Files reach the container on stdin (decided in Phase 2), so the root filesystem stays read-only and no host path is ever mounted; this also works when the worker itself runs in a container.

Wasm runner (`api/internal/runner/wasm`, Phase 10): job image `wasm:<language>`, same job JSON and report as the Python harness (plus `cases` in and `cases` out). Each phase is a fresh instance of the language's WASI module under wazero with the entry as the only file (read-only, `/work`), empty environment, no sockets or subprocesses (WASI preview 1 has none), real clocks, a memory cap (`limits.memory_mb` as wasm pages) and a deadline; the whole job shares 2 × `timeout_ms`, as the Python harness's two phases do. stdout is capped at the job's output limit and stderr at 16 KB; out-of-memory becomes "Memory limit: your program used more than N MB". The interpreters exit 3 with `<line>:<col>: <message>` on a compile error, which becomes `syntax_error`. Modules are built from pinned source (`runners/wasm`, go.sum) with `make wasm-modules` and never fetched at run time; the worker compiles them once at start (about 7 s, cached on disk).

Docker runner (docker CLI): `--network none`, `--read-only` with tmpfs `/work` (16 MB) and `/tmp` (8 MB, noexec), `--user 10001`, `--cap-drop ALL`, `no-new-privileges`, `--memory` = `--memory-swap`, `--pids-limit 64`, `--cpus 1`, nofile/fsize ulimits, no logging, stdout capped at 256 KB and stderr at 16 KB, and a hard kill (`docker rm -f`) at 2 × timeout + 5 s.

Python harness (`runners/python/harness`), PID 1 in the container and non-dumpable (no ptrace or `/proc/1/mem` from learner code):

1. Reads the job from stdin, writes only the learner files to `/work`. Test sources are never written to disk.
2. Statically scans learner files for syntax errors and forbidden calls (AST).
3. If a stdout check exists, runs the entry as `__main__` in a child process (program phase).
4. Runs the tests in a second child that gets the test sources and a random nonce on stdin, and reports one JSON line per test on a dedicated pipe fd. Records without the nonce, for unknown tests or repeated are dropped.
5. Each phase is killed (whole process group) at `timeout_ms`. Child stdout/stderr are pipes the harness reads, so learner `print` output can only ever appear inside the report's `stdout` field.
6. Writes exactly one JSON report line to its own stdout; the Go side parses it and runs the checks.

Known limit: tests execute in the learner's process, so a determined learner could introspect or monkey-patch test objects in memory. Hidden tests deter casual peeking; they are not a vault. The grade itself still comes only from the server.

Server grading runs in `cmd/worker`, which pulls jobs from Redis (`cq:jobs` for Docker jobs and `cq:jobs:wasm` for wasm jobs, each served by its own worker kind, `WORKER_RUNNER=docker|wasm`; a job on the wrong queue gets an error reply and never runs; reply on `cq:result:<id>`, 2 min TTL; jobs older than 60 s in the queue are skipped). The worker has no database access. The API grades in a goroutine and waits up to 15 s: a result in time is stored with progress in one transaction (200); otherwise a `pending` submission is stored, the API returns `202 {submission_id, status: "pending"}`, and the background grade finishes it. `GET /submissions/{id}` polls (owner only). A runner outage before the wait ends is `503 grader_unavailable` and nothing is stored; after it, the submission becomes `error` and does not count as an attempt. At startup, submissions pending for over 5 minutes are marked `error`.

### Content CI (`cq test`)

`cq test <dir>` (needs Docker and the runner image, and `make wasm-modules` for Go/JS lessons, which it runs in-process; part of `make test`) grades every `Testable` exercise's reference answer and wrong answers on the Docker runner and prints one line per case. For mcq, the wrong answers are each single non-answer option. For `project_cli`, the solution and wrong answers are the recorded `solution.json` and `wrong/*.json`. For every exercise, CI must verify:

- [ ] Manifest parses and matches its type's spec schema.
- [ ] `solution.*` passes every check on the server runner.
- [ ] Every file in `wrong/` fails at least one check.
- [ ] Visible tests pass in the browser runner for the solution (headless run via Playwright, Phase 3).

## Browser sandbox and security

Learner code never runs in the main app's page. It runs in a Web Worker inside a sandboxed iframe served from a separate origin, and talks to the app only through `postMessage`.

- **Separate origin.** Serve `sandbox/` from its own domain (e.g. `run.<other-domain>`; in dev, `localhost:5174`). Embed it with `<iframe sandbox="allow-scripts">` and never add `allow-same-origin`.
- **Web Worker.** Pyodide, sql.js and JS execution happen in a worker. The iframe kills it with `worker.terminate()` after `timeout_ms` and reports `timeout`. Recreate the worker for the next run; keep a warm spare to hide Pyodide startup.
- **Message protocol.** App → sandbox: `{ type: "run", runId, runner, entry, files, tests, sql?, timeoutMs }`. Sandbox → app: `{ type: "result", runId, stdout, stderr, checks, runtimeMs, timedOut, table?, error? }` and `{ type: "ready" }`. Both sides validate the shape with zod (strict objects; unknown messages are dropped). The sandbox accepts messages only from its parent window with an allowed app origin (`VITE_APP_ORIGINS`). Because the frame has no `allow-same-origin`, its origin is opaque, so the app accepts replies only when `event.source` is its own frame and `event.origin` is `"null"`. For the same reason the app posts run requests with target origin `"*"` (a named origin never matches an opaque frame, so the browser would drop the message); it posts only to its own frame's `contentWindow`, and a run request carries nothing secret (learner code and visible tests).
- **Workers from an opaque origin.** The worker is bundled as one classic script and started from a blob URL (an opaque-origin page cannot start a same-origin worker). Runtime assets are fetched by absolute URL, so the sandbox serves them with `Access-Control-Allow-Origin: *`; they are public. The sandbox keeps one active worker and one warm spare; a timeout terminates the active worker and promotes the spare.
- **CSP on the sandbox page:** `default-src 'none'; script-src 'self' <sandbox origin> 'wasm-unsafe-eval'; connect-src 'self' <sandbox origin>; worker-src blob:; base-uri 'none'; form-action 'none'; frame-ancestors <app origins>`. The sandbox origin is listed explicitly because `'self'` does not match from an opaque origin. Pyodide and sql.js are self-hosted (`sandbox/scripts/copy-assets.mjs`), so learner code can reach only the sandbox's own static files. Browser Python is Pyodide 314 (Python 3.14), stdlib only; the server runs Python 3.12, so keep browser lessons to features both share.
- **Output.** Render stdout and stderr with `textContent` only. No `innerHTML` or `dangerouslySetInnerHTML` anywhere for learner-derived text.
- **Never ship secrets to the browser.** `PublicView` strips `solution`, hidden tests and `wrong/`. Add an integration test that asserts no server-only file content appears in any API response.

Server-side rules:

- Never trust a client-reported pass. XP, streaks, leagues and course completion are written only by the server grading path.
  - **Exception (decision 2026-11-04, at the author's request): `project_cli` lessons.** Their steps must run on the learner's machine, so the server judges output the CLI *reports* (the boot.dev model), and a determined learner can forge it. They still go through the normal submit path and count for everything: XP, streaks, leagues, Bytes, quests and medals. What stays server-side: the verdict, captures (recomputed from raw output) and the tests. Every other type is unchanged.
- Rate-limit submissions per user (30/min) and per IP (120/min) in Redis, fixed one-minute windows; over the limit is `429 rate_limited` with `Retry-After`. If Redis is down, submits are refused (503), never let through unlimited.
- Docker runner hardening as listed under Grading. The worker host runs no other services and has no credentials to the database beyond writing results.
- **Wasm worker (decision 2026-11-06, Phase 10).** Wasm jobs run inside the worker process, so the wasm worker is never the Docker worker (which holds the root-equivalent Docker socket). It runs in its own container: non-root, read-only root, `cap_drop: ALL`, `no-new-privileges`, memory and pids limits, tmpfs `/tmp` only, and on an internal network that reaches only Redis. Learner code is contained by WASI first (only its own file, no network, no environment) and by that container second. This replaces "one container per submission" for Go and JS lessons only; Python lessons keep the Docker runner.
- Store passwords nowhere (OAuth only). CSRF protection on all state-changing endpoints. Validate all request bodies.

## Data model and API

Postgres holds users, compiled content, submissions, progress and game state. Redis holds leaderboards, rate limits and the job queue.

### Tables

| Table | Key columns |
| --- | --- |
| `users` | id (uuid), email, name, avatar, country, timezone, goal (onboarding: zero, backend, interviews, ai), onboarded_at, created_at |
| `sessions` | id (SHA-256 of the cookie value, never the value), user_id, expires_at, created_at |
| `oauth_identities` | provider, provider_user_id, user_id |
| `courses` | id (slug), title, description, order_index, published, icon, level, tags (text[]), topics |
| `career_paths` | id, title, description, goal, order_index, steps (jsonb `[{course, note}]`), published |
| `chapters` | id, course_id, title, order_index, published |
| `exercises` | id (manifest id), chapter_id, type, version, title, xp, difficulty, free, boss, tags (text[]), order_index, public_view (jsonb), private_spec (jsonb, server-only), assets (jsonb: image name → `<sha256>.<ext>`), published |
| `content_assets` | sha256 (key), mime, size (≤ 1 MB), bytes: lesson images, content-addressed, upserted by `cq build` and never deleted |
| `submissions` | id, user_id, exercise_id, exercise_version, answer (jsonb), status, score, result (jsonb), outcome (jsonb, the game Outcome), created_at |
| `progress` | user_id, exercise_id, status (started, passed), first_passed_at, attempts, used_hint |
| `xp_events` | id, user_id, amount, reason, ref_id, created_at (append-only ledger; unique (user_id, reason, ref_id), so a lesson's pass XP can exist once) |
| `user_stats` | user_id, total_xp, level, bytes, streak_days, streak_freezes, last_active_date, league_tier, title, daily_goal |
| `day_passes` | user_id, day (the user's local date), exercise_id: lessons passed per day (streak week, daily goal) |
| `league_members` | week (Monday, UTC), user_id, tier, group_no, final_rank, result |
| `league_weeks` | week, closed_at (a row means the week's results are final) |
| `subscriptions` | user_id, provider, provider_ref, plan, status, current_period_end |
| `tutor_messages` | id, user_id, exercise_id, role (user, assistant), content, input_tokens, output_tokens, model, created_at |
| `analytics_events` | id, user_id (null when signed out), anon_id, name, props (flat jsonb), path, created_at; kept 90 days |
| `cli_login_codes` | code_hash (SHA-256), user_id, expires_at (10 minutes), used_at (single use), created_at |
| `cli_tokens` | id (SHA-256 of the token), user_id, label (hostname), created_at, last_used_at, expires_at (90 days, extended on use) |
| `client_errors` | fingerprint, day, message, stack, url, user_agent, release, count, users (uuid[], at most 100), first_seen, last_seen; kept 90 days |

Game economy tables (medals, Bytes, inventory, quests) are listed under Game economy. `total_xp` and `bytes` are always derivable from their ledgers; treat `user_stats` as a cache.

### Endpoints (all JSON under `/api/v1`)

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/auth/providers` | Enabled sign-in methods (`github`, `google`, `dev`) |
| GET | `/auth/{github,google}/start`, `/auth/{github,google}/callback` | OAuth web flow, sets the session cookie (404 for a provider that is not configured) |
| POST | `/auth/dev` | Sign in without email; only when `DEV_AUTH=1` (404 otherwise) |
| POST | `/auth/logout` | End the session |
| POST | `/cli/codes` | A one-time code for `codequest login` (session cookie only, not a CLI token) |
| POST | `/cli/token` | Trade a code for a CLI bearer token (no Origin check; 20/h per IP) |
| POST | `/cli/logout` | Revoke the calling CLI token |
| GET | `/cli/next` | The learner's next unpassed `project_cli` lesson (locks apply on `GET /exercises/{id}`) |
| GET | `/me` | Current user (with `goal`, `onboarded`), stats (XP, level, streak, league tier, daily goal) and `features` (`tutor`); entitlement in Phase 5 |
| PATCH | `/me` | Name, timezone, daily goal, title, onboarding goal, `onboarded: true`. XP, level, streak and league fields are rejected |
| GET | `/home` | Continue card, today (goal, streak week), path cards (started courses plus the next course on the learner's path, with next boss), league summary, activity |
| GET | `/profile` | The signed-in learner's own profile: joined date, path, projects shipped, activity (lessons per local day for 53 weeks), projects (every boss: defeated, next, upcoming), skills, recent passes. Never public |
| GET | `/courses` | Course list with progress, status, current chapter, catalog fields, where lessons run (`browser`, `cli`), free lesson count; and career paths |
| GET | `/courses/{id}` | Chapters and exercises with lock state |
| GET | `/exercises/{id}` | `public_view`, type, `assets` (image name → URL), saved draft, last result; 402 if locked |
| GET | `/assets/{sha256}.{ext}` | A lesson image (public; `Cache-Control: immutable`, ETag, `nosniff`, `Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'; sandbox`, `Cross-Origin-Resource-Policy: same-site`); 404 for an unknown hash or an extension that does not match the stored type |
| PUT | `/exercises/{id}/draft` | Autosave answer |
| POST | `/exercises/{id}/run` | Ungraded server Run (`run: server`): visible checks only, nothing stored; 60/min per user |
| POST | `/exercises/{id}/submit` | Grade; returns `Result` + XP awarded + level/streak changes |
| GET | `/submissions/{id}` | Poll a pending submission (owner only; `result` is null while pending) |
| GET | `/exercises/{id}/tutor` | The learner's tutor thread for the exercise and today's hints left |
| POST | `/exercises/{id}/hint` | Ask the AI tutor; the reply streams as Server-Sent Events |
| POST | `/events` | Self-hosted analytics events (batched, allow-listed names) |
| POST | `/client-errors` | Browser error report (error monitoring) |
| GET | `/leagues/current` | This week's league group standings (names only, never other users' emails) |
| GET | `/medals` | Earned and available medals with progress |
| POST | `/billing/checkout` | Start Stripe or Razorpay checkout |
| POST | `/billing/webhook/{provider}` | Provider webhooks |

Errors always use `{ "error": { "code": "locked", "message": "..." } }`, including 402. Codes in use: `not_found`, `bad_request`, `invalid_answer` (400, answer doesn't fit the type's shape; not stored), `unauthorized`, `forbidden_origin`, `unsupported_media_type`, `stale_version` (409, the submitted `version` is not the current one), `run_unavailable` (400, the exercise does not run on the server), `invalid_code` (400, a CLI login code that is unknown, used or expired), `chapter_locked` (403, a lesson behind an earlier chapter's undefeated boss: GET, Run, Submit and power-ups), `insufficient_bytes`, `limit_reached`, `not_owned`, `nothing_to_peek`, `request_reused` (409, economy actions; state unchanged), `tutor_unavailable` (404, the tutor is off), `limit_reached` (also 429 when today's tutor hints are used), `tutor_failed` (an SSE `error` event; the hint is refunded), `rate_limited` (429), `grader_unavailable` (503), `unavailable` (503), `internal`. Request bodies are capped at 128 KB (a 50 KB code file can double when JSON-escaped). Generate TypeScript types from `api/openapi.yaml` with `openapi-typescript` (`make api-types`) and call the API with `openapi-fetch`.

CSRF: every non-GET request must carry an `Origin` header in `ALLOWED_ORIGINS` (browsers always send one on these methods; a missing one is refused) and `Content-Type: application/json` (a cross-site form cannot send it without a CORS preflight). The session cookie is also `SameSite=Lax`. Two exceptions skip the Origin check, never the content type (Phase 8): requests with an `Authorization` header (the CLI; a page cannot attach one cross-site without a preflight, which the API never answers) and `POST /cli/token` (sets no cookie; needs a code shown only to the signed-in user).

### Accounts and sessions

- **Sessions.** A random 32-byte value in the `cq_session` cookie (`HttpOnly`, `SameSite=Lax`, `Secure` when `APP_URL` is https, 30 days). Postgres stores only its SHA-256. Sign-out deletes the row. Expired sessions are deleted hourly.
- **OAuth only** (decision 2026-11-09, at the author's request: no SMTP server to run). Email magic links were removed (migration `0009` drops `magic_links`). The sign-in page passes the browser's timezone as `tz` to `/start`; it rides in the state cookie and is saved on the first sign-in (streak days follow it). Callbacks are limited to 20 per IP per hour. With `APP_ENV=production` the API refuses to start unless GitHub or Google is configured.
- **GitHub.** Hand-rolled OAuth web flow (no new dependency): a random `state` in a 10-minute `HttpOnly` cookie, compared in constant time; scopes `read:user user:email`; the primary verified email is required. A GitHub identity links to the account with that email, or creates one. Enabled when `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` are set.
- **Google.** The same hand-rolled flow against Google's OpenID Connect endpoints: scopes `openid email profile`, `prompt=select_account`, profile from the userinfo endpoint; `email_verified` must be true. A Google identity (`sub`) links to the account with that email, or creates one, so one person signing in with both providers with the same email gets one account. Enabled when `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` are set; the redirect URI is `APP_URL/api/v1/auth/google/callback`.
- **Post-sign-in redirect** (`next`) is accepted only as a same-site path (`/…`, not `//`, no backslashes, not under `/api/`).
- **CLI tokens** (Phase 8, paste-code flow like boot.dev's, without its localhost listener). `codequest login` opens `APP_URL/cli/login`; the signed-in page calls `POST /cli/codes` and shows a 10-character code (`XXXXX-XXXXX`, no ambiguous letters, single use, 10 minutes, stored hashed); the learner pastes it and the CLI gets a random 32-byte token from `POST /cli/token`, kept in `<user config dir>/codequest/config.json` (mode 0600). The API accepts `Authorization: Bearer <token>` in place of the cookie on every endpoint; a request with both is `400`, a bad or expired token `401`. `codequest logout` revokes it (`POST /cli/logout`) and always deletes it locally. Expired codes and tokens are deleted hourly.
- **Dev sign-in** replaces Phase 1's `DEV_USER_EMAIL` (decision 2026-11-03): `DEV_AUTH=1` (on in Compose) enables `POST /auth/dev`, which creates a real session for `dev@codequest.local` (or a given email) without email. Browser tests use it. The default dev user (no email given) starts with 1,000 XP and 10,000 Bytes (decision 2026-11-05; raised from 5,000 on 2026-11-09), granted once through the XP and Byte ledgers (`dev_grant`/`starter`) so every power-up can be tried (the Skeleton Key still needs 5 bosses defeated); dev sign-ins with an email start from zero. The `feature-tour` course has one lesson of every type (mcq single and multiple, browser and server Python, SQL, a boss, CLI command and HTTP) for manual testing. With `APP_ENV=production` the API refuses to start if `DEV_AUTH` is set, no OAuth provider is configured or `APP_URL` is not https.

## Gamification engine

All game logic lives in `api/internal/game` and reacts to one event: a graded submission. It reads only `Result`, the exercise's `xp`, and the user's history, so it works for every exercise type.

```go
func (e *Engine) OnGraded(ctx context.Context, tx pgx.Tx, u UserID, ex Exercise, r Result, p Progress) (Outcome, error)
// Outcome: XPAwarded, NewLevel, StreakDays, MedalsEarned, QuestsCompleted, LeagueRankDelta
```

Rules for v1 (keep all numbers in `game/config.go`):

- **XP.** Award only on the first pass of an exercise. Base = exercise `xp`. Bonuses: +50% first-try pass (the first submission of that exercise passed), +25% no hint used. Each bonus is rounded half away from zero (5 XP → +3, +1). Re-solving a passed exercise gives 0 XP. The ledger's unique key makes "once" hold under concurrent submissions too.
- **Levels.** Everyone starts at level 1; leaving level n needs `round(100 * n^1.5)` cumulative XP (level 2 at 100, level 3 at 283, level 8 at 1852). Return level-ups in `Outcome` for the UI celebration.
- **Streak.** A day counts when the user passes at least one exercise (a re-solve counts too), on the user's own calendar (`users.timezone`, IANA). Each missed day consumes a Streak Freeze if one is held; if the missed days outnumber the freezes, the streak restarts at 1 with that pass and the freezes are kept. Freezes are spent when the next pass arrives; until then the shown streak is the stored one while the gap is coverable, else 0. Earn 1 freeze each time the streak reaches a multiple of 7, max 2 held. A pass dated on or before the last counted day (e.g. after moving timezone west) changes nothing.
- **Leagues.** Weekly, Monday 00:00 UTC. A user joins the current week's group on their first XP of the week: the lowest-numbered group of their tier with fewer than 30 members, else a new group (an advisory lock serialises joins). Standings are by weekly XP, ties broken by who reached it first. Redis keeps a sorted set per group, rebuilt from Postgres after each commit (so the update is idempotent and a lost key heals); Postgres is the source of truth. The API closes finished weeks at startup and every minute (advisory lock plus a `league_weeks` row, so each week closes once even with many replicas): the top 5 of each group move up a tier and the bottom 5 of the rest move down (none above Diamond or below Bronze; in groups under 10, promotion is applied first). Users who earn no XP in a week keep their tier. A pass in the first minute after Monday 00:00 may join with the previous tier; accepted.

All writes for one submission (submission row, progress, xp_events, user_stats, day_passes, league membership; medals and quests later) happen in a single DB transaction, which locks the user's `user_stats` row. The leaderboard update to Redis happens after commit and is idempotent (the group's set is rewritten from Postgres).

Frontend: XP bar and level in the header, streak flame with day count, a short celebration on pass (see `screens/pass-celebration.png`), league page, profile medals tab. Celebrations must be skippable and respect `prefers-reduced-motion`. Animations and sounds (Phase 9) are for rare, earned moments only: lesson clear, level up, boss fight and chests; routine submits, runs and fails stay quiet. The celebration shows on a first pass (when XP was awarded); a re-solve shows the normal result. It omits the design's "Compare with solution" (solutions never reach the browser) and "Share" (no share target yet). Home showed Daily quests and Review as disabled "Coming soon" placeholders (decision 2026-11-03); quests arrived with Phase 7, Review stays a placeholder until reviews exist. The tutor button arrived in Phase 6 (development only; see AI tutor).

## Game economy: Bytes, power-ups, medals, quests

The economy rewards learning but never changes grading. Rule for every feature here: **nothing earned or bought can mark a check as passed or reveal a solution.** `grading` never imports `economy`; the game engine reads power-ups only when it computes bonuses.

Designed screens: `screens/home.png` (daily quests), `screens/boss-fight.png`, `screens/power-ups-shop.png`, `screens/profile-medals.png`.

### Bytes (soft currency)

- Earn only, never sold. Pro grants a small monthly amount.
- Ledger table `byte_events(id, user_id, amount, reason, ref_id, created_at)`, append-only like `xp_events`; balance cached in `user_stats.bytes`.
- Starting values (in `economy/catalog.go`): daily chest 40–80, chapter complete 150, boss defeated 300, medal 50–500 by tier (bronze 50, silver 150, gold 300, diamond 500), league top 5: 200, weekly quest 150.
- Every reward is a ledger row unique on `(user_id, reason, ref_id)` (e.g. `chapter`/chapter id, `medal`/`id:tier`, `buy:<powerup>`/request id), so each lands once; `user_stats.bytes` changes in the same statement batch and a `CHECK (bytes >= 0)` plus the guarded update make overspending impossible.

### Power-ups

| Power-up | Effect | Limits | Cost (Bytes) |
| --- | --- | --- | --- |
| Streak Freeze | Covers one missed day, used automatically | Hold max 2; 1 free per 7-day streak | 300 |
| Streak Repair | Restores a streak lost in the last 48 h | Once a month | 2,000 |
| Holiday Mode | Pauses the streak up to 7 days | Pro only; scheduled 3+ days ahead | Free (Pro) |
| Focus Boost | 2× XP for 15 minutes | 1 per day | 200 |
| Extra Hint | One tutor hint past the daily limit | Sets `used_hint` | 100 |
| Peek | Shows the input of one failed hidden check, never the expected output | Removes the first-try bonus | 150 |
| Second Wind | Keeps the first-try bonus after one failed submit | 1 per chapter | 350 |
| League Shield | No demotion this week | Not two weeks in a row | 700 |
| Skeleton Key | Opens the next locked chapter of a course without passing the boss before it | Needs 5 bosses defeated; hold max 1; once per course; never in chests | 6,000 |
| Double or Nothing | Bet that the daily goal is hit 7 days running; win 2× | One active bet; 100–500 | Stake |

> Decision 2026-11-09, at the author's request: prices were raised to about 1–2 days of typical earning for small items and about two weeks for Streak Repair (a learner on the default goal earns roughly 120–150 Bytes a day), the Skeleton Key was added as a deliberately rare item (about 40 days of typical earning), and chests drop a power-up only half the time.

- Tables: `inventory(user_id, powerup, count)`, `powerup_uses(id, user_id, powerup, ref_id, created_at)`, `active_effects(user_id, effect, ref_id, expires_at)`, `bets(id, user_id, stake, goal, starts_on, status)`, `idempotency(user_id, request_id, action, response json)`. Held Streak Freezes stay in `user_stats.streak_freezes` (`CHECK` 0–2).
- Endpoints: `GET /shop`, `POST /shop/{powerup}/buy`, `POST /powerups/{powerup}/use`, `GET /inventory` (same body as `/shop`). Buying and using are server-side, in one transaction with the Byte ledger, idempotent by request id: the first response is stored (as `json`, byte-for-byte) and replayed; a failed request stores nothing; reusing an id for another action is `409 request_reused`. The client sends a new UUID per click.
- Rules as built (Phase 7):
  - **Streak Freeze** is spent by the daily rollover on each missed day (or lazily by the next pass if the job has not run); a streak the freezes cannot cover is lost (recorded for Streak Repair, freezes kept).
  - **Streak Repair** is used, not bought: 1,500 Bytes when a streak was lost within the last 2 days; it restores the lost length (plus any days since) and is allowed once per 30 days.
  - **Focus Boost**: one use per local day, 15 minutes; the boost equals the rest of the award (2×).
  - **Second Wind** is armed on an exercise right after exactly one failed submit (one per chapter); the next submit keeps the first-try bonus. **Peek** sets `progress.peeked`, which removes the first-try bonus even with Second Wind. Peek with no failed hidden check that declares an input is refused and costs nothing.
  - **League Shield** is used for the current week (not if the previous week was shielded); at the close a demotion becomes a stay.
  - **Double or Nothing**: stake 100–500, uses the daily goal at the time of the bet, starts today; settled at a pass, the rollover or a shop visit: lost once a finished day misses the goal, won (2× the stake back) when all 7 days meet it.
  - **Skeleton Key** is bought only with the `bosses` counter at 5 or more, and only while none is held. It is used on a chapter (`chapter_id`): the first locked chapter of its course, locked only by the undefeated boss of an earlier chapter. That boss is waived for the chapters after it (`chapter_keys(user_id, course_id, gate_chapter_id, chapter_id)`, one row per user and course) but stays unpassed: no XP, Bytes, boss chest or Boss Slayer progress, and its own chapter does not count as complete. A refused use costs nothing.
  - **Holiday Mode** (Pro, Phase 5) and **Extra Hint** (tutor, Phase 6) show in the shop as unavailable ("Pro", "Coming soon") until those phases ship.
  - Daily rollover job (API, every 10 minutes): users whose local date moved on and who were active in 14 days or hold an open bet; spends freezes, settles bets, evaluates medals; once per local date (`user_stats.rolled_over_on`).

### Medals

- About 22 medals in families: Consistency, Mastery, Craft, Memory, League, Secret. Tiers Bronze → Silver → Gold, a few Diamond. Secret medals show "???" until earned.
- `medals(id, family, name, secret, tiers jsonb)` where each tier holds a rule (`{"kind":"streak","gte":30}`), Bytes reward and optional title or cosmetic. `user_medals(user_id, medal_id, tier, earned_at)`.
- Evaluated by the game engine after each graded submission, review, league close and daily rollover.
- As built: 20 medals are defined in `economy/catalog.go` (synced to `medals` at startup) plus Reviewer (Memory), shown as "coming soon" until reviews exist. A rule is a stat and thresholds: user stats (`longest_streak`, `total_xp`) or `counters(user_id, key, value)` bumped on events: first passes (lessons, first_try, no_hint, comebacks, persistent, night 22:00–04:00 and early 05:00–08:00 local), flawless chapters, courses, bosses, goal days, league top 5 / podium / highest tier / shield saves, bets won. Each tier is awarded once (`user_medals` key) with its Bytes once (ledger key). Every tier grants a title ("Boss Slayer I"); the user wears one via `PATCH /me {title}` (only earned titles), shown on the profile and league rows. Cosmetics: not built yet.
- Medals grant titles (shown on profile and league rows) and cosmetics (editor themes, avatar frames). Cosmetics are never sold.
- Medal colors are design tokens: `medal-bronze`, `medal-silver`, `medal-gold`, `medal-diamond`.

### Quests and chests

- Three daily quests drawn from a pool (pass 2 lessons, one first-try pass, finish a review…) plus one weekly quest. `quests(id, cadence, rule jsonb, reward jsonb)`, `user_quests(user_id, quest_id, period, progress, done_at)`.
- Completing all three opens the daily chest: Bytes plus, half the time, a random power-up. Odds live in config and are shown to users.
- As built: daily pool of five (pass 2 lessons, a first-try pass, earn 50 XP, pass a hands-on lesson (`code`, `sql` or `project_cli`, so CLI-only courses can finish it), pass a quiz; "finish a review" waits for reviews), three per user per local day chosen by a hash of user and date (stable all day, nothing stored until progress); weekly pool of three (10 lessons, goal on 5 days, 300 XP), 150 Bytes on completion. Progress counts a lesson once per day (re-solving the same one does not add). The chest is opened with `POST /quests/chest/open` (once per local day, `chests` key): 40–80 Bytes and a draw from Focus Boost 20, Peek 15, Second Wind 10, League Shield 5, no power-up 50 (weights out of 100, shown on Home; the API's `odds` lists "no power-up" as `powerup: ""`, and an empty chest's `powerup` is `""`). Randomness is server-side (`crypto/rand`).

### Boss fights

- A boss is a `code` or `project_cli` exercise with `boss: true`. Its checks are its HP; each passing check removes one HP. Same grader and `Result`; only the UI and rewards differ.
- Decision 2026-11-08, at the author's request: `project_cli` exercises can be bosses, so a hands-on course (Docker, Git) ends each chapter on its own tooling instead of a Python program. HP is the number of step tests (each is one result row). Its result is forgeable like any `project_cli` lesson (see the Phase 8 exception), so its boss rewards are too; this is accepted with the same reasoning. On the page, the HP bar and defeat reward sit under the prompt, a submit from the terminal is the attack, and a result that arrives after the page loaded (via "Check attack result") plays as a fresh attack.
- Defeat reward: +100 XP, boss chest, Boss Slayer progress, unlocks the next chapter.
- Primary button reads "Attack"; results show HIT / BLOCK rows. HP bar uses the `boss-hp` token.
- As built: a code exercise's public view carries `check_counts {visible, hidden}` (rows grading will produce; `cq test` fails if they differ), so the HP bar is right before the first attack. A chapter is locked while any earlier chapter of the course has a boss the user has not passed: its lessons answer `403 chapter_locked` (GET, Run, Submit, power-ups) and the course map shows it locked, unless a Skeleton Key waived that boss (see Power-ups). The first defeat opens the boss chest (300 Bytes and, with the same odds as the daily chest, a power-up) once.

## AI tutor

The tutor gives hints, never full solutions. It lives behind a `tutor.Provider` interface so the LLM vendor can change.

- **Context sent per request:** lesson `prompt.md`, the learner's current code, the latest `Result` (failing check names and feedback), the last 6 tutor messages for this exercise. Never send `solution.*` or hidden test source.
- **System prompt rules:** answer with questions and small nudges; point to the line or concept, not the fix; never output more than 3 consecutive lines of code; if asked for the answer, refuse briefly and give the next hint.
- **Limits:** free users 5 hints/day, paid 50/day, enforced in Redis. Max 400 output tokens. Log token usage per user.
- Using a hint sets `progress.used_hint = true` (removes the no-hint XP bonus).
- Stream responses to the browser with Server-Sent Events.

As built (Phase 6, decision 2026-11-03: **development only for now**, to avoid paying for an LLM this early):

- Provider: `TUTOR_PROVIDER=ollama` with a local Ollama (`OLLAMA_URL`, `TUTOR_MODEL`, default `qwen2.5-coder:3b`), streamed from `/api/chat` over plain HTTP (no new dependency). Empty `TUTOR_PROVIDER` turns the tutor off. With `APP_ENV=production` the API refuses to start if `TUTOR_PROVIDER` is set; production gets a larger hosted model behind the same `tutor.Provider` interface when that is decided. Compose turns it on and reaches the host's Ollama at `host.docker.internal:11434`.
- Off: `GET /me` has `features.tutor = false`, the tutor endpoints answer `404 tutor_unavailable`, the "Ask the tutor" button is disabled ("Coming soon") and Extra Hint stays "Coming soon" in the shop.
- `POST /exercises/{id}/hint` takes `{message (≤ 2000 chars), code}` (the learner's current answer) and answers `text/event-stream`: `delta` `{text}` events, then `done` `{remaining, used_hint}` or `error` `{code, message}`. Errors before the stream (validation, limit, locked chapter) use the JSON envelope. It runs outside the 30 s request timeout with its own 2-minute limit.
- Context: the exercise's public view (prompt or question, SQL setup), the learner's code with line numbers, their latest stored `Result` (check names, got/expected, feedback, stdout) and the last 6 messages. The system prompt carries the rules above; there is no output filtering yet.
- Limit: `TUTOR_DAILY_LIMIT` (default 5) per user per local day, counted in Redis (`tutor:<user>:<date>`). Every user is on the free limit until Phase 5 adds plans. A reply that fails before any text refunds the hint and changes nothing; once text has streamed the hint counts, and the question, reply and token counts are saved (even if the browser left), with `progress.used_hint` set in the same transaction. Token usage is also logged per reply.
- UI: `TutorPanel` at the bottom of the problem pane (code and sql renderers receive it as a `tutor(answerText)` render prop). It shows hints left, warns once that asking removes the no-hint bonus, and renders replies with `Markdown` (no `innerHTML`).
- Deferred until the tutor goes to production: output guardrails (code-length caps, solution-leak checks) and the scripted "give me the answer" eval in the acceptance criteria.

## Payments and entitlements

One `billing` package with a `Provider` interface (`CreateCheckout`, `HandleWebhook`, `CancelSubscription`), implemented for Stripe and Razorpay.

- **Plans:** monthly and annual. Prices live in config, not code. Show purchasing-power-adjusted prices by country (country from IP at checkout, confirmed by billing address).
- **Provider choice:** Razorpay for India, Stripe for everyone else.
- **Entitlement check:** `billing.HasAccess(user, exercise)` returns true if the exercise is `free`, or the user has an active subscription. `GET /exercises/{id}` returns 402 with `{ code: "locked" }` otherwise; the UI shows the content read-only with an upgrade prompt.
- **Webhooks** are the only source of truth for subscription status. Verify signatures. Handle events idempotently by provider event id.
- No card data ever touches our servers; use hosted checkout pages.

## Design system and UI

The UI follows the CodeQuest design system in `docs/design/`: one set of color tokens with a Dark (Charcoal, default) and a Light theme, Inter for interface text, JetBrains Mono for code, and Silkscreen only for game moments (streak, level, +XP, Bytes, medals). Every designed screen is in `docs/design/screens/`; match them.

Implementation rules:

- `web/src/styles/tokens.css` and `web/tailwind.config.ts` are generated from `docs/design/tokens.json` and already in the repo. Add a `make tokens` script that regenerates them; never edit them by hand. No hex values in components.
- Theme: default dark; the header toggle sets `data-theme` on `<html>` and saves the choice in localStorage; respect `prefers-color-scheme` on first visit.
- Motion and sound (Phase 9): `motion-*` durations and `ease-*` easings are tokens; keyframes live in `web/src/styles/effects.css` (hand-written). Sound is on by default; the header speaker toggle mutes it (`cq-sound` in localStorage).
- Amber (`accent`) is the only accent. Use it as a fill with `accent-ink` text; when amber must be text, use `accent-text`.
- Load fonts from Google Fonts (Inter 400–800, JetBrains Mono 400/500, Silkscreen 400/700) or self-host them with `@fontsource`.
- Build these components first, matching `docs/design/components/<Name>/README.md`: `Button`, `Pill`, `GameStats`, `LessonRail`, `ProblemExample`, `CheckList`.

Lesson screen layout (`/learn/:courseId/:exerciseId`; `screens/lesson-*.png`):

1. Header, 56px: logo + wordmark · previous / lesson name and "Chapter · lesson n of m" / next · streak, level, XP bar, theme toggle, avatar. No progress bar in the header.
2. `LessonRail` on the left: collapsed (48px, one vertical segment per lesson) by default; clicking it opens the 236px outline. It is the only progress indicator. Remember open/closed per device.
3. Problem pane: pills, title, statement, `ProblemExample` blocks (Example n: Input / Output / Explanation), Constraints, then hints and the tutor button at the bottom.
4. Editor pane (remaining width): file tab with Run and Submit, Monaco editor themed from the code tokens with the failing line tinted `danger-soft`, then `CheckList` rendering the grader `Result`.
5. Resizable panes (decision 2026-11-03): a drag handle between the problem and editor panes (problem width 240–720px, default 320, the editor keeps at least 360px) and one between the editor and the checks once a result is shown (checks height 96–720px, default 240, the editor keeps at least 160px). The handles are WAI-ARIA window splitters: pointer drag, arrow keys (Shift for larger steps), Home/End, and Enter or double-click to reset. Sizes are remembered per device (`cq-panes` in localStorage).
6. Below 900px the panes stack in one column, the handles are hidden and the rail stays collapsed.

Other screens and their references: Home (`home.png`), Courses (`courses.png`), Course map (`course-map.png`), Pass celebration overlay (`pass-celebration.png`), Onboarding (`onboarding.png`), Landing (`landing.png`), Pricing (`pricing.png`), Boss fight (`boss-fight.png`), Power-ups (`power-ups-shop.png`), Profile overview (`profile-overview.png`), Profile medals (`profile-medals.png`). Only the lesson screen is drawn in both themes; build every screen from tokens so Light works everywhere.

Screens updated 2026-11-06 (`home.png`, `courses.png`, `profile-overview.png`), as built:

- **Activity heatmap** (`components/Heatmap`): one column per week (Monday first), one square per local day, colored by lessons passed (0, 1, 2, 3, 4+ via `surface-raised`, `accent-line` and `accent` at rising opacity), today outlined. Home shows 27 weeks with active days, lessons passed (both for the shown weeks) and the streak ("best yet" when it is the longest); the profile shows 53 weeks.
- **Home path cards**: courses the learner has started ("8 chapters · 2 done · next boss: …"), then the first unstarted course on their path ("Up next", "After <previous course> · <topics>"). With no path and nothing started, the first course is up next.
- **Courses** (`/courses`): career path cards (steps with status dots, % of lessons passed, Continue or Start path) and course cards (icon tile, status, level, where it runs, progress with the current chapter, topics, action). Status: In progress / Completed / Free to start (has free lessons) / Pro (no free lessons, action "Preview"). Search matches title, description, topics and tags; `/` focuses it. Topic filters come from course tags.
- **Profile** (`/profile`, `/profile/medals`, `/profile/projects`, `/profile/activity`): header (path, joined month, worn title colored by its medal tier), Level, Total XP, Current streak, Projects shipped (defeated bosses). Overview: year heatmap, Projects completed (latest 3 defeated plus the next boss; "View code" opens the lesson, which shows the learner's last answer), Skills, Recent medals. Skills rate each chapter of the course the learner last worked in: Bronze once a lesson is passed, Silver when all are, Gold when all passed on the first try (reviews will replace this when they exist). "Edit profile" is an inline form for name and title. The design's @handle and "Share profile" are left out: profiles are private until public profiles are decided.

## Build phases

Build in this order. Each phase ends with a working, demoable product and must pass its acceptance criteria before the next starts.

### Phase 1 — Skeleton and content pipeline

Scope: monorepo, Docker Compose, Postgres migrations, Go API skeleton, React app shell with routing and the design tokens, `cq validate` and `cq build`, exercise registry with `mcq` plugin, one sample course with 2 chapters.

- [ ] `docker compose up` starts Postgres, Redis, API and web; `make test` passes.
- [ ] `cq build content/` loads the sample course into Postgres; an invalid manifest fails with file path and reason.
- [ ] Course page lists chapters and exercises in manifest order.
- [ ] An `mcq` exercise renders, submits, and shows pass/fail with explanation.
- [ ] The app shell uses `tokens.css`; switching `data-theme` flips dark/light.

### Phase 2 — Code exercises and grading

Scope: `code` type, check primitives (`unit_tests`, `stdout_equals`, `stdout_regex`, `forbidden_calls`), Docker runner + worker + Redis queue, Monaco editor, Python harness that writes JSON results to a dedicated fd, lesson screen layout.

- [ ] Submitting the reference solution of every sample `code` exercise returns `pass`; every `wrong/` file returns `fail`.
- [ ] `while True: pass` returns `timeout` within limit + 1 s; a fork bomb and a network call fail safely.
- [ ] Learner `print` output that mimics the result JSON cannot change the grade.
- [ ] `cq test` runs the content CI checks and exits non-zero on any failure.
- [ ] The lesson screen matches `screens/lesson-dark.png` and `screens/lesson-light.png`.

### Phase 3 — Browser sandbox

Scope: `sandbox/` app on a separate origin, worker with Pyodide and sql.js, postMessage protocol, Run button, `sql` exercise type.

- [ ] Run shows output for Python in under 1 s after first load.
- [ ] Code in the sandbox cannot read the app's cookies or DOM (automated Playwright test).
- [ ] Infinite loop is killed at `timeout_ms`; the UI stays responsive.
- [ ] For every sample exercise, the solution passes visible tests in both browser and server runners.

### Phase 4 — Accounts, progress, gamification

Scope: magic link + GitHub auth, progress, XP ledger, levels, streaks (with Streak Freeze), leagues, header UI, pass celebration, Home and Course map screens.

- [ ] First pass awards XP once; re-solving awards 0.
- [ ] Streak logic unit-tested across timezones, missed days and freezes.
- [ ] Weekly league job promotes and demotes correctly (unit test with fake clock).
- [ ] Forged client requests cannot award XP (integration test).

### Phase 5 — Payments and paywall

Scope: billing interface, Stripe + Razorpay, entitlement checks, pricing page, PPP pricing.

- [ ] Locked exercises return 402 for free users and open after a test-mode subscription.
- [ ] Replaying the same webhook twice changes state once.
- [ ] Cancelling keeps access until `current_period_end`.

### Phase 6 — AI tutor and polish

Scope: tutor with SSE streaming and limits, onboarding, landing page, analytics events, error monitoring.

- [ ] Tutor never returns the reference solution for any sample exercise (scripted eval of 20 "give me the answer" prompts). **Deferred** (decision 2026-11-03): the tutor is development-only on a small local model; guardrails and this eval come with the production tutor.
- [x] Daily hint limits enforced per plan (one plan until Phase 5: `TUTOR_DAILY_LIMIT`, integration-tested).
- [x] Lighthouse score above 90 on the landing page (`make lighthouse`, part of `make test`: every category ≥ 90, mobile profile, production build).

As built:

- **Landing** (`/` when signed out; `screens/landing.png`): hero with a static "Try it" lesson card (no runtime loads on this page), features, real courses from `GET /courses`, how a lesson works, FAQ. Pricing, "For teams", Privacy/Terms links and the cancellation FAQ are left out until Phase 5 and those pages exist. Every route except the landing page is a lazy chunk.
- **Onboarding** (`/welcome`; `screens/onboarding.png`): step 1 is goal and time per day (10/20/40 min → daily goal 1/2/4), step 2 picks the first course (preselected for the goal). The design's "5-question check to skip what you know" is not built (no design or question bank). Finishing or skipping sets `users.onboarded_at`; `/` sends a signed-in learner who has not onboarded to `/welcome`. Accounts from before Phase 6 count as onboarded.
- **Analytics** (self-hosted, no vendor): `track(name, props)` batches events to `POST /events` every 5 s or when the tab is hidden (`fetch` with `keepalive`). Names are allow-listed on the server (`page_view`, `landing_cta`, `sign_in`, `onboarding_step`, `onboarding_completed`, `lesson_opened`, `run_clicked`, `submit_clicked`, `tutor_opened`, `tutor_asked`, `shop_bought`); props are a flat object ≤ 1 KB; ≤ 20 events per batch; 60/min per user, 240/min per IP. Do Not Track and Global Privacy Control turn it off. A random `anon_id` in localStorage (`cq-anon`) counts signed-out visitors.
- **Error monitoring** (self-hosted): `window` `error` and `unhandledrejection` handlers and a React `ErrorBoundary` send `POST /client-errors` (once per message per page, ≤ 20 per page; the URL fragment is dropped because it can hold a sign-in token). The API groups them by a fingerprint of the message and top stack frame (line numbers and bundle hashes removed) per day and logs each at error level. API panics are recovered, logged with stack and request id, and answered with `500 internal`. `VITE_RELEASE` tags reports. A vendor (e.g. Sentry) can be added later behind the same calls.

> Decision 2026-11-03: Phases 5 (payments) and 6 (AI tutor) are skipped for now and Phase 7 is built first, at the author's request. Phase 6 was then built after Phase 7 (same day), with the tutor development-only; Phase 5 is still skipped. Phase 7 features that need them (Holiday Mode, Extra Hint, Pro's monthly Bytes, review quests and the Memory medals) show as unavailable or "coming soon" until those phases ship.

### Phase 7 — Game economy

Scope: Byte ledger, shop and inventory, power-up effects, medals with tiers and titles, daily quests and chests, boss fight UI, Home quests card, Profile medals tab.

- [ ] Byte balance always equals the sum of `byte_events`; buying with too few Bytes is refused and nothing changes.
- [ ] Replaying a buy or use request with the same request id changes state once.
- [ ] No power-up can change a `Result`: an integration test submits wrong code with every power-up active and still gets `fail`.
- [ ] Peek returns only the input of a failed hidden check, never expected output or test source.
- [ ] Streak Freeze is consumed automatically on a missed day; the streak survives; a third freeze cannot be held.
- [ ] Medal rules are evaluated after submission, review, league close and daily rollover; each tier is awarded once.
- [ ] A boss shows HP = number of checks; the HP bar drops by one per passing check; defeat unlocks the next chapter.

### Phase 8 — Learner CLI and `project_cli` lessons

Scope: the `codequest` CLI (`api/cmd/codequest`, standard library only: `login`, `logout`, `status`, `run [id]`, `submit [id]`, `config`), the `project_cli` type (see Exercise type plugin system), CLI login codes and bearer tokens, `/cli/login` page and the `project_cli` renderer, sample course `content/command-line`.

> Decision 2026-11-04: Phase 8 is built before Phase 5 (payments), at the author's request, as Phase 7 was. Lessons that are not free return 402 to the CLI like to the browser.

- [x] `codequest login` with a code from `/cli/login` stores a token. A code works once, expires after 10 minutes and is refused on a second use.
- [x] Bearer requests skip the Origin check. A bearer request that also carries a cookie is refused. A revoked or expired token gives 401.
- [x] `codequest run` executes the steps locally and prints their output with no verdict. `submit` returns the server's `Result`. A passing submit awards XP, league points and quest progress once, exactly like a web lesson.
- [x] The server recomputes captures, so forged variables in an answer have no effect. A malformed answer is a 400 and is not stored.
- [x] `cq test` passes every `project_cli` `solution.json` and fails every `wrong/*.json`.
- [x] The CLI never sends the CodeQuest token to lesson URLs (unit test against an `httptest` server).
- [x] The config file is 0600. `logout` deletes the token locally even when the API is unreachable.

As built: `run` and `submit` list the commands and ask before running them (`-y` skips the prompt), because lesson commands run in the learner's own shell. Output is printed as it arrives: the raw output, then the tests it will be judged on; only `submit` prints ✓/✗, from the server. `submit` exits 0 on pass and 1 otherwise. API and app URLs default to localhost and are set at build time with `-ldflags "-X main.defaultAPIURL=… -X main.defaultAppURL=…"` (or `CODEQUEST_API_URL` / `CODEQUEST_APP_URL` / `codequest config`); plain http is refused except for localhost. The lesson page shows the steps, tests, copyable `codequest run/submit <id>` commands and the latest result ("Check for new results" refetches the lesson). `/cli/login` has no design screen; it reuses the sign-in page's shell. The name `codequest` changes with the product name (open question).

### Phase 9 — Game feel and lesson images

> Decision 2026-11-05: built after Phase 8, at the author's request, with Phase 5 (payments) still skipped. Taken with the author: sounds are synthesized in the browser (Web Audio, no audio files); diagrams are image files (no Mermaid, which would need `innerHTML`); sound is on by default with a mute toggle. After review, the author narrowed effects to rare, earned moments: routine submits, runs and fails get no sound and no new animation, so the game never feels childish or nagging.

Scope: motion tokens, a sound engine and mute toggle, effects for lesson clear, level up, boss fights and chests, and lesson images (`assets/`, the assets endpoint, rendering in prompts and mcq options). No new dependencies.

- [x] Normal submits and runs play no sound and show no new animation (`effects/quiet.test.ts`: only the lesson-clear dialog, the boss fight and the chest may import the sound engine).
- [x] Lesson clear, level up, boss hit/block/defeat and chest open each play their animation and sound. With `prefers-reduced-motion` there is no motion but the end state shows. Muting silences every sound and survives a reload.
- [x] Every overlay closes with Escape, a backdrop click or its button, and focus returns to where it was.
- [x] Effects never change a `Result`: they read only the `Result` and `Outcome` the server returned.
- [x] No hex colors in components (confetti reads token variables at runtime). Token files are regenerated with `make tokens`, and `tokens:check` passes.
- [x] `cq validate` rejects a missing image, an external image, an over-size file, a disallowed or mislabelled type, a scriptable SVG, a server-only name and an unused file in `assets/`.
- [x] The assets endpoint serves the stored type with `nosniff`, the sandbox CSP, an ETag and immutable caching. An unknown hash or a mismatched extension is a 404.
- [x] Prompts and mcq options show lesson images in both themes. External or unknown images show as their alt text, and nothing goes through `innerHTML`.
- [x] Lighthouse stays ≥ 90 on the landing page (effects code is not on the landing page; `make lighthouse`: 99 / 96 / 96 / 100).

As built:

- **Sounds** (`web/src/effects/sound.ts`): six short chiptune recipes (`success`, `levelUp`, `bossHit`, `bossBlock`, `bossDefeat`, `chestOpen`), oscillators plus gain envelopes through one quiet master gain (0.4 × the user's volume). The `AudioContext` is created on first use, which always follows a click. Missing Web Audio is a silent no-op. Callers use `play(name)` only, so recorded files can replace the recipes later.
- **Lesson clear**: one soft chime and an XP count-up. On a level up the fanfare replaces the chime, and the dialog shows a "LEVEL UP" banner, the badge popping from the old level to the new one, and one confetti burst (none under reduced motion).
- **Boss fight**: after Attack the checks are revealed 150 ms apart. Each HIT drains one HP segment (`bossHit`), each BLOCK plays `bossBlock`, and a full clear ends with the `bossDefeat` fanfare. The celebration opens once the reveal has finished. A result loaded with the page (not a fresh attack) shows at once, silently.
- **Chests** (`ChestOpen`): the daily chest on Home and the boss chest (from the celebration's "Open boss chest") open in a modal. The chest wobbles while the server decides (at least 600 ms), then the lid opens (`chestOpen`), the Bytes count up and the power-up rises. The loot comes from the server; the modal only shows it.
- **Images**: `GET /exercises/{id}` carries `assets`. `LessonAssetsProvider` (on the lesson page) lets `Markdown` and the mcq renderer resolve `assets/…` names, so no renderer props changed. Images are lazy, framed on `surface-raised` so transparent diagrams read in both themes, and open full size on click. `Markdown` outside a lesson (tutor replies) loads no images at all. Assets of a locked lesson are reachable by anyone who has the hash; accepted, because hashes cannot be guessed and images are not answers.
- **Effects lab** (`/dev/effects`, development builds only: the route and its chunk are left out of production builds by `import.meta.env.DEV`): buttons for every sound, each lesson-clear scenario (plain, level up, boss with chest, boss with level up), a boss attack with any number of checks and hits, and the daily chest with a fast server, a slow server or an error. It needs no sign-in and makes no requests; it drives the real components with made-up `Result`s and `Outcome`s.
- There are no PNG designs for these moments; component notes are in `docs/design/components/ChestOpen`, `LevelUpBurst` and `SoundToggle`.

### Phase 10 — WASM runner and more languages

> Decision 2026-11-06: proposed and accepted the same day at the author's request ("start building"), with Phase 5 still skipped. Each new language runs as an interpreter compiled to WebAssembly (WASI) on the server, so adding a language is a module plus a registry entry rather than a Docker image and a harness. Python keeps the Docker runner, which is also the tier for lessons that need compiled libraries (LangChain, numpy). Taken during the build: the modules are Go programs (yaegi for Go, goja for JS) built from pinned source instead of downloaded binaries; Run for wasm lessons goes to the server for now, because the modules are 12 MB and 21 MB; wasm jobs run in a separate, locked-down worker, never in the Docker worker that holds the Docker socket. Git and Docker courses use `project_cli` (Phase 8) and need no grader change.

Scope:
- `runners/wasm`: a Go module with `cmd/go` (yaegi; the standard library limited to the allowlist in `symbols/`, regenerated with `make -C runners/wasm symbols`) and `cmd/js` (goja plus `prelude.js`).
- `api/internal/runner/wasm`: the wazero runner.
- Queue routing by image, a `Route` runner for `cq test`, and `WORKER_RUNNER=wasm`.
- The `worker-wasm` Compose service and `runners/wasm/Dockerfile`.
- The `io_cases` and `exit_code` checks, and `code` lessons with `language: go | js`.
- Go and JavaScript highlighting in Monaco.
- Sample courses `content/go-basics` and `content/js-basics`.

Adding a language: a WASI module in `runners/wasm/cmd/<lang>` that runs `<lang> /work/<entry>` and follows the exit-code convention (3 plus `<line>:<col>: <message>` for compile errors); an entry in `wasm.Languages` and `runner.WasmLanguages`; a line in `languages` in `exercise/code`; the Monaco language; a sample lesson. No grading code changes.

- [x] `runner: wasm` grades a Go and a JS sample course with `solution` and `wrong/`. `cq test` passes all 7 lessons (7 solutions pass, 14 wrong answers fail), and `make test` runs them with `CQ_REQUIRE_WASM=1`.
- [~] Browser and server give the same rows: wasm lessons have no browser runner yet. Run uses the server and the same module as Submit. The Playwright content spec passes every Go and JS solution through the real API, Redis queue and `worker-wasm` (about 100–200 ms each).
- [x] Abuse ends within limits and the runner keeps serving (`wasm_test.go`):
  - an infinite loop over five phases ends within the 2 × `timeout_ms` job budget;
  - a single 100 MB allocation fails under a 64 MB cap and succeeds under 256 MB;
  - an output flood is cut at the output limit;
  - a normal job still runs afterwards.
- [x] Learner code can't escape its file or reach the network (`TestIsolation`):
  - it cannot read `/etc/passwd`, `/proc/self/environ` or `../` paths, and cannot write to `/work`;
  - it sees an empty environment;
  - `net`, `net/http`, `os/exec`, `syscall` and `unsafe` fail to compile;
  - JS has no `fetch` and no `require` beyond `fs` stdin.
- [x] Hidden cases stay server-side. Their stdin and stdout are not in the public view, Run never sends them, and hidden rows show no got/expected (`code/wasm_test.go`, `checks/io_cases_test.go`).
- [x] Wasm jobs never reach the Docker worker, and a job on the wrong queue gets an error reply without running (`TestQueueRoutesWasmJobs`).
- [x] `worker-wasm` starts in its container: non-root, read-only, no capabilities, internal network to Redis only. It compiles both modules in about 7 s.
- [x] Grading benchmark: 30 grades of each boss on the same machine (below).
- [x] Security sections updated (Browser sandbox and security, Grading).

Benchmark, 30 grades each, with the solution and every check on the same machine:

| Lesson | Runner | p50 | p95 | Notes |
| --- | --- | --- | --- | --- |
| `go.tower.tonnage-board` (6 io cases) | wasm | 287 ms | 311 ms | 7 module runs, about 40 ms each |
| `js.board.receipt` (5 io cases) | wasm | 136 ms | 147 ms | 6 runs, about 22 ms each |
| `fd.py.fizzbuzz-boss` | Docker | 218 ms | 248 ms | 1 container |

The wasm worker holds about 350 MB resident for all its jobs, mostly compiled modules. Cases run one after another; running them in parallel is the next lever if a lesson's grading feels slow.

Known limits:
- yaegi doesn't enforce `go build` rules such as unused variables and imports.
- Go panics read "reflect: slice index out of range" instead of the compiler's wording.
- `encoding/json` is not in the Go allowlist (Go 1.27's version doesn't extract).
- JS output formatting follows Node only for common values.

## Conventions and rules for Claude Code

- Work one phase at a time. Before coding, list the files you will create or change. After coding, run tests and report which acceptance criteria pass.
- Do not add dependencies beyond the stack table without stating why in the summary. Added so far: `zod` (sandbox message validation), `pyodide`, `sql.js` (self-hosted browser runtimes), `@playwright/test` (browser tests, run from the official Docker image so no browsers are installed on the host), `openapi-fetch` + `openapi-typescript` (typed client from `openapi.yaml`), `react-markdown` (lesson prompts without `innerHTML`), `go.yaml.in/yaml/v3` (manifests), `google/uuid`, `redis/go-redis/v9` (queue and rate limits), `monaco-editor` (self-hosted via `loader.config({ monaco })`, never the CDN; one lazy chunk), `lighthouse` 13.5.0 (dev only: the landing page audit in `make lighthouse`, using the Playwright image's Chromium), `tetratelabs/wazero` (pure-Go WASI runtime for the wasm runner; no cgo), and in the separate `runners/wasm` module `traefik/yaegi` (Go interpreter) and `dop251/goja` (JS engine), compiled to WASI. The runner image pins `langchain-core 1.7.3`, `langchain-ollama 1.1.0`, `pydantic 2.14.2`. Pinned: Tailwind 3 (the generated config is v3 format), TypeScript 5 (`openapi-typescript` needs its compiler API), Go directive 1.26 (golangci-lint's build).
- Go: standard layout above, `gofmt` + `golangci-lint`, errors wrapped with `fmt.Errorf("...: %w", err)`, `context.Context` first argument, table-driven tests, no global state except the plugin registries.
- SQL only through `sqlc`-generated queries. Every schema change is a new migration; never edit an applied one.
- React: function components, TypeScript strict mode, no `any` outside the renderer registry, server state through TanStack Query, components under 200 lines.
- New exercise type checklist: Go package under `internal/exercise/<type>`, spec schema + `ParseSpec`, `PublicView`, `Grade`, React renderer under `web/src/exercise-types/<type>`, one sample exercise with `solution` and `wrong/`, tests. Nothing else in the codebase should need to change; if it does, the abstraction is wrong — stop and flag it.
- New check kind checklist: implement `Check`, register it, document its manifest fields in this spec, add tests with a passing and a failing case.
- Security rules in Browser sandbox and security are hard requirements, not suggestions. Never relax them to make a test pass.
- Config through environment variables, loaded once in `main.go`; ship a `.env.example`.
- Keep this spec current: when a decision changes, update the relevant section in the same change.

## Open questions

- Final product name and domains (app, sandbox origin). Needed before Phase 3. The learner CLI (`codequest`, Phase 8) is renamed with it, and its release build needs the production API and app URLs and an install path (`go install …@latest` assumes a public module path).
- ~~First course topic for launch.~~ Resolved (Phase 1): LangChain, sample course `content/langchain-basics`, adapted from the author's genai lesson notes; to be refined into the full course.
- ~~LangChain `code` exercises (Phase 2).~~ Resolved: `codequest/python:3.12` ships the pinned LangChain packages and `cq.fakes`; the sample LangChain code exercises use `runner: docker-python` and fakes only.
- sql manifest (Phase 3): `prompt`, `setup` (public data, shown to the learner), `solution` (stored in `private_spec` because hidden checks compare against it), optional `starter`, `run`, `limits.timeout_ms` (default 2000). Hidden datasets are server-only.
- Phase 2 interim UI: Run became available in Phase 3; the tutor button is disabled until Phase 6; unsent answers are kept in `localStorage` (`cq-draft:<id>:v<version>`) until `PUT /exercises/{id}/draft` exists; the header's GameStats arrived in Phase 4.
- Hosting (decision 2026-11-09, at the author's request): free plans for now, as in `docs/DEPLOY.md` and `render.yaml`. One Render web service runs the API with the built web app on the same origin (`WEB_DIR`, so the session cookie stays first-party), plus Render Key Value for Redis and Neon for Postgres; every start runs `cq migrate up && cq build`. No grading workers or sandbox: only quiz and CLI lessons are gradable there, which is all the current courses use.
- Dev Compose mounts the host Docker socket into the worker (root-equivalent). Production runs the worker on a dedicated host (later gVisor).
- LangChain in the browser (Phase 3): Pyodide would need self-hosted `langchain-core` and `pydantic` wheels (CSP forbids fetching from PyPI). Decide whether these exercises are submit-only, graded on the server.
- ~~More languages (proposed 2026-11-06).~~ Resolved: Phase 10 accepted and built; Go and JS run on the `wasm` runner. Still open: Run in the browser for wasm lessons (modules are 12 MB for Go and 21 MB for JS, so it waits on a size and load-time decision), unit tests in Go/JS, and moving plain Python to WASI CPython.
- Exact prices per plan and per country tier. Needed before Phase 5.
- ~~LLM provider for the tutor.~~ Development: local Ollama (Phase 6). The production provider and model are still open, along with the guardrails that come with them.
