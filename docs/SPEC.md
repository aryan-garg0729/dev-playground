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
    cmd/worker/main.go      # server-side grading worker
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
  sandbox/                  # static runner page (served from separate origin)
    index.html  worker.ts   # Pyodide / sql.js / JS execution
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

var registry = map[string]Type{}

func Register(t Type) { registry[t.Name()] = t }   // called from each plugin's init()
func Get(name string) (Type, bool) { t, ok := registry[name]; return t, ok }
```

React registry (`web/src/exercise-types/index.ts`):

```ts
export interface ExerciseRendererProps<V = unknown, A = unknown> {
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
| `sql` | `{ query: string }` | `sqljs` |

Later types (do not build in v1): `fill_blank`, `predict_output`, `order_steps`, `project_cli` (learner runs a local CLI that posts results), `open_answer` (AI-graded).

Prefer configuration over new types. "Fix the bug" is a `code` exercise with buggy starter code. "Make it faster" is `code` plus a `max_runtime` check. Only add a type when the UI is genuinely different.

## Grading

The server is the only authority on whether a submission passed. The browser runs code for instant feedback; the server re-grades every submission before awarding XP.

### Two paths

1. **Run (client, ungraded).** The learner clicks Run. The sandbox executes their code plus visible tests and shows output. Nothing is stored.
2. **Submit (server, authoritative).** The browser POSTs the answer. The API calls the type's `Grade`, which executes all checks (visible and hidden) on a server runner, stores the result, and returns it. Only this path awards XP.

Client runners map to server runners: `pyodide` → `docker-python`, `js` → `docker-node`, `sqljs` → in-process SQLite (`modernc.org/sqlite`). `mcq` grades in-process with no runner. The same test files must produce the same result on both.

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
}
```

`status` is `pass` only when every check passes. Game, tutor and analytics code must depend only on `Result`, never on the exercise type.

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
| `sql_result_equals` | Compares result set to the solution's; `ordered` flag | sql |

```go
type Check interface {
    Kind() string
    Run(ctx context.Context, in CheckInput) ([]CheckResult, error)
}
```

### Runner interface (`api/internal/runner`)

```go
type Runner interface {
    Run(ctx context.Context, job Job) (Output, error)
}
type Job struct {
    Image   string            // e.g. codequest/python:3.12
    Files   map[string]string // learner files + test files
    Cmd     []string
    Limits  Limits            // timeout, memory, pids, output bytes
}
type Output struct { Stdout, Stderr string; ExitCode, RuntimeMs int; TimedOut bool }
```

Docker runner requirements: `--network none`, read-only root filesystem with a tmpfs work dir, non-root user, `--memory`, `--pids-limit 64`, CPU quota, hard wall-clock kill, output capped. Test harnesses emit results as one JSON line on a dedicated file descriptor so learner `print` calls cannot spoof them.

Server grading runs in `cmd/worker`, which pulls jobs from a Redis queue. The API waits up to 15 s for the result, otherwise returns `202` with a submission id for polling.

### Content CI (`cq test`)

For every exercise, CI must verify:

- [ ] Manifest parses and matches its type's spec schema.
- [ ] `solution.*` passes every check on the server runner.
- [ ] Every file in `wrong/` fails at least one check.
- [ ] Visible tests pass in the browser runner for the solution (headless run via Playwright, Phase 3).

## Browser sandbox and security

Learner code never runs in the main app's page. It runs in a Web Worker inside a sandboxed iframe served from a separate origin, and talks to the app only through `postMessage`.

- **Separate origin.** Serve `sandbox/` from its own domain (e.g. `run.<other-domain>`; in dev, `localhost:5174`). Embed it with `<iframe sandbox="allow-scripts">` and never add `allow-same-origin`.
- **Web Worker.** Pyodide, sql.js and JS execution happen in a worker. The iframe kills it with `worker.terminate()` after `timeout_ms` and reports `timeout`. Recreate the worker for the next run; keep a warm spare to hide Pyodide startup.
- **Message protocol.** App → sandbox: `{ type: "run", runId, runner, files, tests, timeoutMs }`. Sandbox → app: `{ type: "result", runId, stdout, stderr, checks, runtimeMs, timedOut }`. Both sides check `event.origin` and validate the shape with zod; unknown messages are dropped.
- **CSP on the sandbox page:** `default-src 'none'; script-src 'self' 'wasm-unsafe-eval' <cdn>; connect-src 'self' <cdn>; worker-src 'self' blob:`. Learner code must not reach arbitrary hosts. Self-host Pyodide if possible so `connect-src` can be `'self'` only.
- **Output.** Render stdout and stderr with `textContent` only. No `innerHTML` or `dangerouslySetInnerHTML` anywhere for learner-derived text.
- **Never ship secrets to the browser.** `PublicView` strips `solution`, hidden tests and `wrong/`. Add an integration test that asserts no server-only file content appears in any API response.

Server-side rules:

- Never trust a client-reported pass. XP, streaks, leagues and course completion are written only by the server grading path.
- Rate-limit submissions per user (e.g. 30/min) and per IP in Redis.
- Docker runner hardening as listed under Grading. The worker host runs no other services and has no credentials to the database beyond writing results.
- Store passwords nowhere (magic link + OAuth only). CSRF protection on all state-changing endpoints. Validate all request bodies.

## Data model and API

Postgres holds users, compiled content, submissions, progress and game state. Redis holds leaderboards, rate limits and the job queue.

### Tables

| Table | Key columns |
| --- | --- |
| `users` | id (uuid), email, name, avatar, country, timezone, created_at |
| `sessions` | id, user_id, expires_at |
| `courses` | id (slug), title, description, order_index, published |
| `chapters` | id, course_id, title, order_index, published |
| `exercises` | id (manifest id), chapter_id, type, version, title, xp, difficulty, free, boss, tags (text[]), order_index, public_view (jsonb), private_spec (jsonb, server-only), published |
| `submissions` | id, user_id, exercise_id, exercise_version, answer (jsonb), status, score, result (jsonb), created_at |
| `progress` | user_id, exercise_id, status (started, passed), first_passed_at, attempts, used_hint |
| `xp_events` | id, user_id, amount, reason, ref_id, created_at (append-only ledger) |
| `user_stats` | user_id, total_xp, level, bytes, streak_days, streak_freezes, last_active_date, league_tier, title |
| `subscriptions` | user_id, provider, provider_ref, plan, status, current_period_end |
| `tutor_messages` | id, user_id, exercise_id, role, content, created_at |

Game economy tables (medals, Bytes, inventory, quests) are listed under Game economy. `total_xp` and `bytes` are always derivable from their ledgers; treat `user_stats` as a cache.

### Endpoints (all JSON under `/api/v1`)

| Method | Path | Purpose |
| --- | --- | --- |
| POST | `/auth/magic-link` | Send login email |
| GET | `/auth/callback` | Magic link / OAuth callback, sets session |
| GET | `/me` | Current user, stats, entitlement |
| GET | `/courses` | Course list with progress |
| GET | `/courses/{id}` | Chapters and exercises with lock state |
| GET | `/exercises/{id}` | `public_view`, type, saved draft, last result; 402 if locked |
| PUT | `/exercises/{id}/draft` | Autosave answer |
| POST | `/exercises/{id}/submit` | Grade; returns `Result` + XP awarded + level/streak changes |
| GET | `/submissions/{id}` | Poll a pending submission |
| POST | `/exercises/{id}/hint` | AI tutor message |
| GET | `/leagues/current` | This week's league standings |
| GET | `/medals` | Earned and available medals with progress |
| POST | `/billing/checkout` | Start Stripe or Razorpay checkout |
| POST | `/billing/webhook/{provider}` | Provider webhooks |

Errors always use `{ "error": { "code": "locked", "message": "..." } }`, including 402. Codes in use: `not_found`, `bad_request`, `invalid_answer` (400, answer doesn't fit the type's shape; not stored), `unauthorized`, `forbidden_origin`, `unsupported_media_type`, `stale_version` (409, the submitted `version` is not the current one), `internal`. Generate TypeScript types from `api/openapi.yaml` with `openapi-typescript` (`make api-types`) and call the API with `openapi-fetch`.

CSRF (until real sessions in Phase 4): every non-GET request must be `Content-Type: application/json`, and a browser `Origin` header, when present, must be in `ALLOWED_ORIGINS`.

Phase 1 has no auth. When `DEV_USER_EMAIL` is set, the API upserts that user at startup and every request acts as it; submissions and progress are stored, and no XP is awarded. Phase 4 removes this.

## Gamification engine

All game logic lives in `api/internal/game` and reacts to one event: a graded submission. It reads only `Result`, the exercise's `xp`, and the user's history, so it works for every exercise type.

```go
func (e *Engine) OnGraded(ctx context.Context, tx pgx.Tx, u UserID, ex Exercise, r Result, p Progress) (Outcome, error)
// Outcome: XPAwarded, NewLevel, StreakDays, MedalsEarned, QuestsCompleted, LeagueRankDelta
```

Rules for v1 (keep all numbers in `game/config.go`):

- **XP.** Award only on the first pass of an exercise. Base = exercise `xp`. Bonuses: +50% first-try pass, +25% no hint used. Re-solving a passed exercise gives 0 XP.
- **Levels.** Level n needs `100 * n^1.5` cumulative XP (rounded). Return level-ups in `Outcome` for the UI celebration.
- **Streak.** A day counts when the user passes at least one exercise, in the user's own timezone. Missing a day consumes a Streak Freeze if one is held, otherwise resets to 0. Earn 1 freeze per 7-day streak, max 2 held.
- **Leagues.** Weekly, Monday 00:00 UTC. Users are bucketed into groups of up to 30 within a tier (Bronze, Silver, Gold, Diamond). Redis sorted set per group, scored by weekly XP. A weekly job promotes the top 5 and demotes the bottom 5.

All writes for one submission (submission row, progress, xp_events, user_stats, medals, quests) happen in a single DB transaction. The leaderboard update to Redis happens after commit and is idempotent (keyed by submission id).

Frontend: XP bar and level in the header, streak flame with day count, a short celebration on pass (see `screens/pass-celebration.png`), league page, profile medals tab. Celebrations must be skippable and respect `prefers-reduced-motion`.

## Game economy: Bytes, power-ups, medals, quests

The economy rewards learning but never changes grading. Rule for every feature here: **nothing earned or bought can mark a check as passed or reveal a solution.** `grading` never imports `economy`; the game engine reads power-ups only when it computes bonuses.

Designed screens: `screens/home.png` (daily quests), `screens/boss-fight.png`, `screens/power-ups-shop.png`, `screens/profile-medals.png`.

### Bytes (soft currency)

- Earn only, never sold. Pro grants a small monthly amount.
- Ledger table `byte_events(id, user_id, amount, reason, ref_id, created_at)`, append-only like `xp_events`; balance cached in `user_stats.bytes`.
- Starting values (in `economy/config.go`): daily chest 40–80, chapter complete 150, boss defeated 300, medal 50–500 by tier, league top 5: 200.

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

- Tables: `inventory(user_id, powerup, count)`, `powerup_uses(id, user_id, powerup, ref_id, created_at)`, `active_effects(user_id, effect, expires_at)`, `bets(id, user_id, stake, starts_on, status)`.
- Endpoints: `GET /shop`, `POST /shop/{powerup}/buy`, `POST /powerups/{powerup}/use`, `GET /inventory`. Buying and using are server-side, in one transaction with the Byte ledger, idempotent by request id.

### Medals

- About 22 medals in families: Consistency, Mastery, Craft, Memory, League, Secret. Tiers Bronze → Silver → Gold, a few Diamond. Secret medals show "???" until earned.
- `medals(id, family, name, secret, tiers jsonb)` where each tier holds a rule (`{"kind":"streak","gte":30}`), Bytes reward and optional title or cosmetic. `user_medals(user_id, medal_id, tier, earned_at)`.
- Evaluated by the game engine after each graded submission, review, league close and daily rollover.
- Medals grant titles (shown on profile and league rows) and cosmetics (editor themes, avatar frames). Cosmetics are never sold.
- Medal colors are design tokens: `medal-bronze`, `medal-silver`, `medal-gold`, `medal-diamond`.

### Quests and chests

- Three daily quests drawn from a pool (pass 2 lessons, one first-try pass, finish a review…) plus one weekly quest. `quests(id, cadence, rule jsonb, reward jsonb)`, `user_quests(user_id, quest_id, period, progress, done_at)`.
- Completing all three opens the daily chest: Bytes plus a random power-up. Odds live in config and are shown to users.

### Boss fights

- A boss is a `code` exercise with `boss: true`. Its checks are its HP; each passing check removes one HP. Same grader and `Result`; only the UI and rewards differ.
- Defeat reward: +100 XP, boss chest, Boss Slayer progress, unlocks the next chapter.
- Primary button reads "Attack"; results show HIT / BLOCK rows. HP bar uses the `boss-hp` token.

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
5. Below 900px the panes stack in one column and the rail stays collapsed.

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
- Do not add dependencies beyond the stack table without stating why in the summary. Added so far: `openapi-fetch` + `openapi-typescript` (typed client from `openapi.yaml`), `react-markdown` (lesson prompts without `innerHTML`), `go.yaml.in/yaml/v3` (manifests), `google/uuid`. Pinned: Tailwind 3 (the generated config is v3 format), TypeScript 5 (`openapi-typescript` needs its compiler API), Go directive 1.26 (golangci-lint's build).
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
- LangChain `code` exercises (Phase 2): the `docker-python` image needs `langchain-core`, `langchain-ollama` and fake chat models; grading uses fakes only, never a live model or the network.
- LangChain in the browser (Phase 3): Pyodide would need self-hosted `langchain-core` and `pydantic` wheels (CSP forbids fetching from PyPI). Decide whether these exercises are submit-only, graded on the server.
- Exact prices per plan and per country tier. Needed before Phase 5.
- LLM provider for the tutor. Needed before Phase 6.
