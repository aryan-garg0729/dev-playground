# CodeQuest Build Spec — React + Go

## How to use this doc

This is the source of truth for building CodeQuest (working name), a gamified, hands-on coding-learning platform in the style of boot.dev, with a React frontend and a Go backend. `CLAUDE.md` at the repo root points here.

Prompt Claude Code one phase at a time, for example: "Implement Phase 1 from docs/SPEC.md. Stop when its acceptance criteria pass and summarize what you built."

Design references live in `docs/design/`: `tokens.json` (all design tokens), `README.md` (the design system's rules), `components/` (component guidelines and HTML previews) and `screens/` (a PNG of every designed screen).

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
| Server-side execution | Docker containers (later gVisor) | Worker pool behind a Redis job queue |
| Auth | Email magic link + GitHub OAuth | Session cookie, `HttpOnly`, `SameSite=Lax` |
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
    internal/
      content/              # load + validate exercise packages
      exercise/             # type registry + plugins (code, mcq, sql)
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
  content/                  # courses as files (see Content format)
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
  runner: pyodide           # pyodide | sqljs | js | docker-python | docker-go
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

# chapter.yaml
id: lc.chat-models          # lowercase words joined by . or -, globally unique
title: Chat models & messages
exercises: [01-message-roles, 02-runnable-methods] # folder names, in order
```

- Every folder must be listed, and every listed folder must exist with its manifest; both are errors otherwise. Hidden folders (`.x`) are ignored; folders inside an exercise (e.g. `wrong/`) belong to that exercise.
- Exercise header: `id`, `type`, `version` (≥ 1), `title` and `xp` (≥ 0) are required; `difficulty` defaults to 1 (1–2 shows as Easy), `tags` to `[]`, `free` and `boss` to false. Only `code` exercises can be bosses.
- `cq validate <dir>` and `cq build <dir>` print every problem as `path: reason` and exit 1. `<dir>` is a folder of courses or a single course folder.
- `cq build` writes everything in one transaction, upserting by id. Courses, chapters and exercises no longer in the content are **unpublished, never deleted**, so progress and submissions keyed on old ids survive.
- mcq: `question` names a markdown file; `answer` and `explain` live only in `private_spec`. `explain` is returned after every submit (pass or fail) as the `Feedback` of the implicit `exact_choice` check.
- Run location (decided in Phase 3): `code` and `sql` specs take `run: browser | server | off`. `browser` runs in the sandbox iframe (Pyodide for Python, sql.js for SQL); `server` sends the answer to `POST /exercises/{id}/run`, which runs it on the grading worker with visible checks only, ungraded and not stored; `off` hides Run. Defaults: `pyodide` → browser, `docker-python` → server, `sql` → browser. Changing it is a manifest edit plus `cq build`; content CI checks the solution's visible checks wherever `run:` points, so a lesson moved to `browser` whose packages don't load there fails CI.
- code (decided in Phase 2): `language: python`; `runner: docker-python` (submit-only) or `pyodide` (server grades it on `docker-python`); `prompt` (markdown), `entry` (a module file name such as `main.py`), `starter` and `solution` are required files; `limits` default to 3000 ms / 128 MB, capped at 10000 ms / 512 MB; `examples`, `constraints` and `hints` are optional and public. At least one check. `solution.*` and `wrong/*.py` are read by `cq build`/`cq test` but never stored in `private_spec`; hidden test sources are stored there (server-only). Visible test files are sent to the browser as `visible_tests`.
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

Later types (do not build in v1): `fill_blank`, `predict_output`, `order_steps`, `project_cli` (learner runs a local CLI that posts results), `open_answer` (AI-graded).

Prefer configuration over new types. "Fix the bug" is a `code` exercise with buggy starter code. "Make it faster" is `code` plus a `max_runtime` check. Only add a type when the UI is genuinely different.

## Grading

The server is the only authority on whether a submission passed. The browser runs code for instant feedback; the server re-grades every submission before awarding XP.

### Two paths

1. **Run (client, ungraded).** The learner clicks Run. The sandbox executes their code plus visible tests and shows output. Nothing is stored.
2. **Submit (server, authoritative).** The browser POSTs the answer. The API calls the type's `Grade`, which executes all checks (visible and hidden) on a server runner, stores the result, and returns it. Only this path awards XP.

Client runners map to server runners: `pyodide` → `docker-python`, `js` → `docker-node`, `sqljs` → SQLite in the `docker-python` grading container (decided in Phase 3, instead of in-process `modernc.org/sqlite`, so untrusted SQL never runs inside the API process). `mcq` grades in-process with no runner. The same test files must produce the same result on both.

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

Docker runner (docker CLI): `--network none`, `--read-only` with tmpfs `/work` (16 MB) and `/tmp` (8 MB, noexec), `--user 10001`, `--cap-drop ALL`, `no-new-privileges`, `--memory` = `--memory-swap`, `--pids-limit 64`, `--cpus 1`, nofile/fsize ulimits, no logging, stdout capped at 256 KB and stderr at 16 KB, and a hard kill (`docker rm -f`) at 2 × timeout + 5 s.

Python harness (`runners/python/harness`), PID 1 in the container and non-dumpable (no ptrace or `/proc/1/mem` from learner code):

1. Reads the job from stdin, writes only the learner files to `/work`. Test sources are never written to disk.
2. Statically scans learner files for syntax errors and forbidden calls (AST).
3. If a stdout check exists, runs the entry as `__main__` in a child process (program phase).
4. Runs the tests in a second child that gets the test sources and a random nonce on stdin, and reports one JSON line per test on a dedicated pipe fd. Records without the nonce, for unknown tests or repeated are dropped.
5. Each phase is killed (whole process group) at `timeout_ms`. Child stdout/stderr are pipes the harness reads, so learner `print` output can only ever appear inside the report's `stdout` field.
6. Writes exactly one JSON report line to its own stdout; the Go side parses it and runs the checks.

Known limit: tests execute in the learner's process, so a determined learner could introspect or monkey-patch test objects in memory. Hidden tests deter casual peeking; they are not a vault. The grade itself still comes only from the server.

Server grading runs in `cmd/worker`, which pulls jobs from Redis (`cq:jobs` list; reply on `cq:result:<id>`, 2 min TTL; jobs older than 60 s in the queue are skipped). The worker has no database access. The API grades in a goroutine and waits up to 15 s: a result in time is stored with progress in one transaction (200); otherwise a `pending` submission is stored, the API returns `202 {submission_id, status: "pending"}`, and the background grade finishes it. `GET /submissions/{id}` polls (owner only). A runner outage before the wait ends is `503 grader_unavailable` and nothing is stored; after it, the submission becomes `error` and does not count as an attempt. At startup, submissions pending for over 5 minutes are marked `error`.

### Content CI (`cq test`)

`cq test <dir>` (needs Docker and the runner image; part of `make test`) grades every `Testable` exercise's reference answer and wrong answers on the Docker runner and prints one line per case. For mcq, the wrong answers are each single non-answer option. For every exercise, CI must verify:

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
- Rate-limit submissions per user (30/min) and per IP (120/min) in Redis, fixed one-minute windows; over the limit is `429 rate_limited` with `Retry-After`. If Redis is down, submits are refused (503), never let through unlimited.
- Docker runner hardening as listed under Grading. The worker host runs no other services and has no credentials to the database beyond writing results.
- Store passwords nowhere (magic link + OAuth only). CSRF protection on all state-changing endpoints. Validate all request bodies.

## Data model and API

Postgres holds users, compiled content, submissions, progress and game state. Redis holds leaderboards, rate limits and the job queue.

### Tables

| Table | Key columns |
| --- | --- |
| `users` | id (uuid), email, name, avatar, country, timezone, created_at |
| `sessions` | id (SHA-256 of the cookie value, never the value), user_id, expires_at, created_at |
| `magic_links` | token_hash, email, next, expires_at, used_at (single use, 15 minutes) |
| `oauth_identities` | provider, provider_user_id, user_id |
| `courses` | id (slug), title, description, order_index, published |
| `chapters` | id, course_id, title, order_index, published |
| `exercises` | id (manifest id), chapter_id, type, version, title, xp, difficulty, free, boss, tags (text[]), order_index, public_view (jsonb), private_spec (jsonb, server-only), published |
| `submissions` | id, user_id, exercise_id, exercise_version, answer (jsonb), status, score, result (jsonb), outcome (jsonb, the game Outcome), created_at |
| `progress` | user_id, exercise_id, status (started, passed), first_passed_at, attempts, used_hint |
| `xp_events` | id, user_id, amount, reason, ref_id, created_at (append-only ledger; unique (user_id, reason, ref_id), so a lesson's pass XP can exist once) |
| `user_stats` | user_id, total_xp, level, bytes, streak_days, streak_freezes, last_active_date, league_tier, title, daily_goal |
| `day_passes` | user_id, day (the user's local date), exercise_id: lessons passed per day (streak week, daily goal) |
| `league_members` | week (Monday, UTC), user_id, tier, group_no, final_rank, result |
| `league_weeks` | week, closed_at (a row means the week's results are final) |
| `subscriptions` | user_id, provider, provider_ref, plan, status, current_period_end |
| `tutor_messages` | id, user_id, exercise_id, role, content, created_at |

Game economy tables (medals, Bytes, inventory, quests) are listed under Game economy. `total_xp` and `bytes` are always derivable from their ledgers; treat `user_stats` as a cache.

### Endpoints (all JSON under `/api/v1`)

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/auth/providers` | Enabled sign-in methods (`magic_link`, `github`, `dev`) |
| POST | `/auth/magic-link` | Email a sign-in link (same answer whether or not the account exists) |
| POST | `/auth/magic-link/verify` | Redeem the link's token, set the session cookie |
| GET | `/auth/github/start`, `/auth/github/callback` | GitHub OAuth web flow, sets the session cookie |
| POST | `/auth/dev` | Sign in without email; only when `DEV_AUTH=1` (404 otherwise) |
| POST | `/auth/logout` | End the session |
| GET | `/me` | Current user and stats (XP, level, streak, league tier, daily goal); entitlement in Phase 5 |
| PATCH | `/me` | Name, timezone, daily goal. XP, level, streak and league fields are rejected |
| GET | `/home` | Continue card, today (goal, streak week), paths, league summary |
| GET | `/courses` | Course list with progress |
| GET | `/courses/{id}` | Chapters and exercises with lock state |
| GET | `/exercises/{id}` | `public_view`, type, saved draft, last result; 402 if locked |
| PUT | `/exercises/{id}/draft` | Autosave answer |
| POST | `/exercises/{id}/run` | Ungraded server Run (`run: server`): visible checks only, nothing stored; 60/min per user |
| POST | `/exercises/{id}/submit` | Grade; returns `Result` + XP awarded + level/streak changes |
| GET | `/submissions/{id}` | Poll a pending submission (owner only; `result` is null while pending) |
| POST | `/exercises/{id}/hint` | AI tutor message |
| GET | `/leagues/current` | This week's league group standings (names only, never other users' emails) |
| GET | `/medals` | Earned and available medals with progress |
| POST | `/billing/checkout` | Start Stripe or Razorpay checkout |
| POST | `/billing/webhook/{provider}` | Provider webhooks |

Errors always use `{ "error": { "code": "locked", "message": "..." } }`, including 402. Codes in use: `not_found`, `bad_request`, `invalid_answer` (400, answer doesn't fit the type's shape; not stored), `unauthorized`, `forbidden_origin`, `unsupported_media_type`, `stale_version` (409, the submitted `version` is not the current one), `run_unavailable` (400, the exercise does not run on the server), `invalid_link` (400, a magic link that is unknown, used or expired), `chapter_locked` (403, a lesson behind an earlier chapter's undefeated boss: GET, Run, Submit and power-ups), `insufficient_bytes`, `limit_reached`, `not_owned`, `nothing_to_peek`, `request_reused` (409, economy actions; state unchanged), `rate_limited` (429), `grader_unavailable` (503), `unavailable` (503), `internal`. Request bodies are capped at 128 KB (a 50 KB code file can double when JSON-escaped). Generate TypeScript types from `api/openapi.yaml` with `openapi-typescript` (`make api-types`) and call the API with `openapi-fetch`.

CSRF: every non-GET request must carry an `Origin` header in `ALLOWED_ORIGINS` (browsers always send one on these methods; a missing one is refused) and `Content-Type: application/json` (a cross-site form cannot send it without a CORS preflight). The session cookie is also `SameSite=Lax`.

### Accounts and sessions

- **Sessions.** A random 32-byte value in the `cq_session` cookie (`HttpOnly`, `SameSite=Lax`, `Secure` when `APP_URL` is https, 30 days). Postgres stores only its SHA-256. Sign-out deletes the row. Expired sessions and links are deleted hourly.
- **Magic links.** `POST /auth/magic-link` emails `APP_URL/sign-in/verify#token=…`. The token is in the URL fragment, so it never reaches a server or proxy log, and the page redeems it with `POST /auth/magic-link/verify` only after a click, so mail scanners that open links cannot use it up. Tokens are single use, expire after 15 minutes and are stored hashed. Limits: 5 emails per address per hour, 20 sign-in attempts per IP per hour. Without `SMTP_ADDR` (development) the email is written to the API log; production requires SMTP. The first sign-in saves the browser's timezone (streak days follow it).
- **GitHub.** Hand-rolled OAuth web flow (no new dependency): a random `state` in a 10-minute `HttpOnly` cookie, compared in constant time; scopes `read:user user:email`; the primary verified email is required. A GitHub identity links to the account with that email, or creates one. Enabled when `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` are set.
- **Post-sign-in redirect** (`next`) is accepted only as a same-site path (`/…`, not `//`, no backslashes, not under `/api/`).
- **Dev sign-in** replaces Phase 1's `DEV_USER_EMAIL` (decision 2026-11-03): `DEV_AUTH=1` (on in Compose) enables `POST /auth/dev`, which creates a real session for `dev@codequest.local` (or a given email) without email. Browser tests use it. With `APP_ENV=production` the API refuses to start if `DEV_AUTH` is set, `SMTP_ADDR` is missing or `APP_URL` is not https.

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

Frontend: XP bar and level in the header, streak flame with day count, a short celebration on pass (see `screens/pass-celebration.png`), league page, profile medals tab. Celebrations must be skippable and respect `prefers-reduced-motion`. The celebration shows on a first pass (when XP was awarded); a re-solve shows the normal result. It omits the design's "Compare with solution" (solutions never reach the browser) and "Share" (no share target yet). Home showed Daily quests and Review as disabled "Coming soon" placeholders (decision 2026-11-03); quests arrived with Phase 7, Review stays a placeholder until reviews exist.

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
| Streak Freeze | Covers one missed day, used automatically | Hold max 2; 1 free per 7-day streak | 200 |
| Streak Repair | Restores a streak lost in the last 48 h | Once a month | 1,500 |
| Holiday Mode | Pauses the streak up to 7 days | Pro only; scheduled 3+ days ahead | Free (Pro) |
| Focus Boost | 2× XP for 15 minutes | 1 per day | 150 |
| Extra Hint | One tutor hint past the daily limit | Sets `used_hint` | 100 |
| Peek | Shows the input of one failed hidden check, never the expected output | Removes the first-try bonus | 120 |
| Second Wind | Keeps the first-try bonus after one failed submit | 1 per chapter | 250 |
| League Shield | No demotion this week | Not two weeks in a row | 400 |
| Double or Nothing | Bet that the daily goal is hit 7 days running; win 2× | One active bet; 100–500 | Stake |

- Tables: `inventory(user_id, powerup, count)`, `powerup_uses(id, user_id, powerup, ref_id, created_at)`, `active_effects(user_id, effect, ref_id, expires_at)`, `bets(id, user_id, stake, goal, starts_on, status)`, `idempotency(user_id, request_id, action, response json)`. Held Streak Freezes stay in `user_stats.streak_freezes` (`CHECK` 0–2).
- Endpoints: `GET /shop`, `POST /shop/{powerup}/buy`, `POST /powerups/{powerup}/use`, `GET /inventory` (same body as `/shop`). Buying and using are server-side, in one transaction with the Byte ledger, idempotent by request id: the first response is stored (as `json`, byte-for-byte) and replayed; a failed request stores nothing; reusing an id for another action is `409 request_reused`. The client sends a new UUID per click.
- Rules as built (Phase 7):
  - **Streak Freeze** is spent by the daily rollover on each missed day (or lazily by the next pass if the job has not run); a streak the freezes cannot cover is lost (recorded for Streak Repair, freezes kept).
  - **Streak Repair** is used, not bought: 1,500 Bytes when a streak was lost within the last 2 days; it restores the lost length (plus any days since) and is allowed once per 30 days.
  - **Focus Boost**: one use per local day, 15 minutes; the boost equals the rest of the award (2×).
  - **Second Wind** is armed on an exercise right after exactly one failed submit (one per chapter); the next submit keeps the first-try bonus. **Peek** sets `progress.peeked`, which removes the first-try bonus even with Second Wind. Peek with no failed hidden check that declares an input is refused and costs nothing.
  - **League Shield** is used for the current week (not if the previous week was shielded); at the close a demotion becomes a stay.
  - **Double or Nothing**: stake 100–500, uses the daily goal at the time of the bet, starts today; settled at a pass, the rollover or a shop visit: lost once a finished day misses the goal, won (2× the stake back) when all 7 days meet it.
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
- Completing all three opens the daily chest: Bytes plus a random power-up. Odds live in config and are shown to users.
- As built: daily pool of five (pass 2 lessons, a first-try pass, earn 50 XP, pass a coding lesson, pass a quiz; "finish a review" waits for reviews), three per user per local day chosen by a hash of user and date (stable all day, nothing stored until progress); weekly pool of three (10 lessons, goal on 5 days, 300 XP), 150 Bytes on completion. Progress counts a lesson once per day (re-solving the same one does not add). The chest is opened with `POST /quests/chest/open` (once per local day, `chests` key): 40–80 Bytes and one power-up drawn from Focus Boost 40, Peek 30, Second Wind 20, League Shield 10 (weights, shown on Home). Randomness is server-side (`crypto/rand`).

### Boss fights

- A boss is a `code` exercise with `boss: true`. Its checks are its HP; each passing check removes one HP. Same grader and `Result`; only the UI and rewards differ.
- Defeat reward: +100 XP, boss chest, Boss Slayer progress, unlocks the next chapter.
- Primary button reads "Attack"; results show HIT / BLOCK rows. HP bar uses the `boss-hp` token.
- As built: a code exercise's public view carries `check_counts {visible, hidden}` (rows grading will produce; `cq test` fails if they differ), so the HP bar is right before the first attack. A chapter is locked while any earlier chapter of the course has a boss the user has not passed: its lessons answer `403 chapter_locked` (GET, Run, Submit, power-ups) and the course map shows it locked. The first defeat opens the boss chest (300 Bytes and a power-up) once.

## AI tutor

The tutor gives hints, never full solutions. It lives behind a `tutor.Provider` interface so the LLM vendor can change.

- **Context sent per request:** lesson `prompt.md`, the learner's current code, the latest `Result` (failing check names and feedback), the last 6 tutor messages for this exercise. Never send `solution.*` or hidden test source.
- **System prompt rules:** answer with questions and small nudges; point to the line or concept, not the fix; never output more than 3 consecutive lines of code; if asked for the answer, refuse briefly and give the next hint.
- **Limits:** free users 5 hints/day, paid 50/day, enforced in Redis. Max 400 output tokens. Log token usage per user.
- Using a hint sets `progress.used_hint = true` (removes the no-hint XP bonus).
- Stream responses to the browser with Server-Sent Events.

## Payments and entitlements

One `billing` package with a `Provider` interface (`CreateCheckout`, `HandleWebhook`, `CancelSubscription`), implemented for Stripe and Razorpay.

- **Plans:** monthly and annual. Prices live in config, not code. Show purchasing-power-adjusted prices by country (country from IP at checkout, confirmed by billing address).
- **Provider choice:** Razorpay for India, Stripe for everyone else.
- **Entitlement check:** `billing.HasAccess(user, exercise)` returns true if the exercise is `free`, or the user has an active subscription. `GET /exercises/{id}` returns 402 with `{ code: "locked" }` otherwise; the UI shows the content read-only with an upgrade prompt.
- **Webhooks** are the only source of truth for subscription status. Verify signatures. Handle events idempotently by provider event id.
- No card data ever touches our servers; use hosted checkout pages.

## Design system and UI

The UI follows the CodeQuest design system in `docs/design/`: one set of color tokens with a Dark (Charcoal, default) and a Light theme, Inter for interface text, JetBrains Mono for code, and Pixelify Sans only for game moments (streak, level, +XP, Bytes, medals). Every designed screen is in `docs/design/screens/`; match them.

Implementation rules:

- `web/src/styles/tokens.css` and `web/tailwind.config.ts` are generated from `docs/design/tokens.json` and already in the repo. Add a `make tokens` script that regenerates them; never edit them by hand. No hex values in components.
- Theme: default dark; the header toggle sets `data-theme` on `<html>` and saves the choice in localStorage; respect `prefers-color-scheme` on first visit.
- Amber (`accent`) is the only accent. Use it as a fill with `accent-ink` text; when amber must be text, use `accent-text`.
- Load fonts from Google Fonts (Inter 400–800, JetBrains Mono 400/500, Pixelify Sans 500/700) or self-host them with `@fontsource`.
- Build these components first, matching `docs/design/components/<Name>/README.md`: `Button`, `Pill`, `GameStats`, `LessonRail`, `ProblemExample`, `CheckList`.

Lesson screen layout (`/learn/:courseId/:exerciseId`; `screens/lesson-*.png`):

1. Header, 56px: logo + wordmark · previous / lesson name and "Chapter · lesson n of m" / next · streak, level, XP bar, theme toggle, avatar. No progress bar in the header.
2. `LessonRail` on the left: collapsed (48px, one vertical segment per lesson) by default; clicking it opens the 236px outline. It is the only progress indicator. Remember open/closed per device.
3. Problem pane: pills, title, statement, `ProblemExample` blocks (Example n: Input / Output / Explanation), Constraints, then hints and the tutor button at the bottom.
4. Editor pane (remaining width): file tab with Run and Submit, Monaco editor themed from the code tokens with the failing line tinted `danger-soft`, then `CheckList` rendering the grader `Result`.
5. Resizable panes (decision 2026-11-03): a drag handle between the problem and editor panes (problem width 240–720px, default 320, the editor keeps at least 360px) and one between the editor and the checks once a result is shown (checks height 96–720px, default 240, the editor keeps at least 160px). The handles are WAI-ARIA window splitters: pointer drag, arrow keys (Shift for larger steps), Home/End, and Enter or double-click to reset. Sizes are remembered per device (`cq-panes` in localStorage).
6. Below 900px the panes stack in one column, the handles are hidden and the rail stays collapsed.

Other screens and their references: Home (`home.png`), Course map (`course-map.png`), Pass celebration overlay (`pass-celebration.png`), Onboarding (`onboarding.png`), Landing (`landing.png`), Pricing (`pricing.png`), Boss fight (`boss-fight.png`), Power-ups (`power-ups-shop.png`), Profile medals (`profile-medals.png`). Only the lesson screen is drawn in both themes; build every screen from tokens so Light works everywhere.

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

- [ ] Tutor never returns the reference solution for any sample exercise (scripted eval of 20 "give me the answer" prompts).
- [ ] Daily hint limits enforced per plan.
- [ ] Lighthouse score above 90 on the landing page.

> Decision 2026-11-03: Phases 5 (payments) and 6 (AI tutor) are skipped for now and Phase 7 is built first, at the author's request. Phase 7 features that need them (Holiday Mode, Extra Hint, Pro's monthly Bytes, review quests and the Memory medals) show as unavailable or "coming soon" until those phases ship.

### Phase 7 — Game economy

Scope: Byte ledger, shop and inventory, power-up effects, medals with tiers and titles, daily quests and chests, boss fight UI, Home quests card, Profile medals tab.

- [ ] Byte balance always equals the sum of `byte_events`; buying with too few Bytes is refused and nothing changes.
- [ ] Replaying a buy or use request with the same request id changes state once.
- [ ] No power-up can change a `Result`: an integration test submits wrong code with every power-up active and still gets `fail`.
- [ ] Peek returns only the input of a failed hidden check, never expected output or test source.
- [ ] Streak Freeze is consumed automatically on a missed day; the streak survives; a third freeze cannot be held.
- [ ] Medal rules are evaluated after submission, review, league close and daily rollover; each tier is awarded once.
- [ ] A boss shows HP = number of checks; the HP bar drops by one per passing check; defeat unlocks the next chapter.

## Conventions and rules for Claude Code

- Work one phase at a time. Before coding, list the files you will create or change. After coding, run tests and report which acceptance criteria pass.
- Do not add dependencies beyond the stack table without stating why in the summary. Added so far: `zod` (sandbox message validation), `pyodide`, `sql.js` (self-hosted browser runtimes), `@playwright/test` (browser tests, run from the official Docker image so no browsers are installed on the host), `openapi-fetch` + `openapi-typescript` (typed client from `openapi.yaml`), `react-markdown` (lesson prompts without `innerHTML`), `go.yaml.in/yaml/v3` (manifests), `google/uuid`, `redis/go-redis/v9` (queue and rate limits), `monaco-editor` (self-hosted via `loader.config({ monaco })`, never the CDN; one lazy chunk). The runner image pins `langchain-core 1.7.3`, `langchain-ollama 1.1.0`, `pydantic 2.14.2`. Pinned: Tailwind 3 (the generated config is v3 format), TypeScript 5 (`openapi-typescript` needs its compiler API), Go directive 1.26 (golangci-lint's build).
- Go: standard layout above, `gofmt` + `golangci-lint`, errors wrapped with `fmt.Errorf("...: %w", err)`, `context.Context` first argument, table-driven tests, no global state except the plugin registries.
- SQL only through `sqlc`-generated queries. Every schema change is a new migration; never edit an applied one.
- React: function components, TypeScript strict mode, no `any` outside the renderer registry, server state through TanStack Query, components under 200 lines.
- New exercise type checklist: Go package under `internal/exercise/<type>`, spec schema + `ParseSpec`, `PublicView`, `Grade`, React renderer under `web/src/exercise-types/<type>`, one sample exercise with `solution` and `wrong/`, tests. Nothing else in the codebase should need to change; if it does, the abstraction is wrong — stop and flag it.
- New check kind checklist: implement `Check`, register it, document its manifest fields in this spec, add tests with a passing and a failing case.
- Security rules in Browser sandbox and security are hard requirements, not suggestions. Never relax them to make a test pass.
- Config through environment variables, loaded once in `main.go`; ship a `.env.example`.
- Keep this spec current: when a decision changes, update the relevant section in the same change.

## Open questions

- Final product name and domains (app, sandbox origin). Needed before Phase 3.
- ~~First course topic for launch.~~ Resolved (Phase 1): LangChain, sample course `content/langchain-basics`, adapted from the author's genai lesson notes; to be refined into the full course.
- ~~LangChain `code` exercises (Phase 2).~~ Resolved: `codequest/python:3.12` ships the pinned LangChain packages and `cq.fakes`; the sample LangChain code exercises use `runner: docker-python` and fakes only.
- sql manifest (Phase 3): `prompt`, `setup` (public data, shown to the learner), `solution` (stored in `private_spec` because hidden checks compare against it), optional `starter`, `run`, `limits.timeout_ms` (default 2000). Hidden datasets are server-only.
- Phase 2 interim UI: Run became available in Phase 3; the tutor button is disabled until Phase 6; unsent answers are kept in `localStorage` (`cq-draft:<id>:v<version>`) until `PUT /exercises/{id}/draft` exists; the header's GameStats arrived in Phase 4.
- Dev Compose mounts the host Docker socket into the worker (root-equivalent). Production runs the worker on a dedicated host (later gVisor).
- LangChain in the browser (Phase 3): Pyodide would need self-hosted `langchain-core` and `pydantic` wheels (CSP forbids fetching from PyPI). Decide whether these exercises are submit-only, graded on the server.
- Exact prices per plan and per country tier. Needed before Phase 5.
- LLM provider for the tutor. Needed before Phase 6.
