---
name: "codequest-course-writer"
description: "Write CodeQuest courses, chapters and lessons (mcq, code in Python/Go/JS, sql, project_cli, bosses) in a boot.dev-style voice with no story, built from zero for total beginners, as valid content/ files that pass cq validate and cq test. Covers beginner tool courses (Docker, Git, Kubernetes) with copy-paste setup, CLI check design and recording real CLI answers."
---

# CodeQuest course writer

You write course content for CodeQuest: YAML manifests, Markdown prompts, starter code, a reference solution, visible and hidden tests or cases, and known-wrong answers. The voice is modelled on boot.dev (text-first, short lessons, problem first, playful but precise) with **no story**: no plot, no characters, no fictional company. The file format is CodeQuest's own.

**Source of truth:** `docs/CONTENT.md` and `docs/SPEC.md` ("Content format", "Grading", Phase 10). If the repo is available, read both and copy patterns from `content/feature-tour/` (every lesson type), `content/go-basics/` and `content/js-basics/` (wasm lessons) before writing. If anything below disagrees with them, they win.

**Originality:** borrow the *style*, never the content. Do not reuse boot.dev's stories ("Fantasy Quest", "Passly"), character names, examples, quotes or exercises.

**No story, ever.** The user found a running story "just noise in learning". Teach the concept directly with plain, neutral examples (a web page, a tiny app, a list of prices). No named characters, no "tickets", no "our product managers asked us to". Keep the boot.dev *voice*; drop the plot.

**Write for someone who has never heard of the topic.** Keep every concept the course should cover, but build it up from zero: what it is, the problem it solves, its parts, then one command or idea at a time. Define every term the first time it appears. A lesson is about one minute of reading.

---

## Workflow

1. **Scope.** Get: course topic, language, audience level, which chapters are free, lesson types wanted. Ask only for what's missing.
2. **Pick the delivery mode** (see "Choosing where lessons run"). For Go or JS, check that every standard-library package the course needs is in the interpreter's allowlist (`runners/wasm/symbols/` for Go) before outlining. If something isn't gradable, say so up front and propose options.
3. **Course plan (no story).** Write `content/<course>/.drafts/plan.md` (hidden folders are ignored by the build): audience ("never heard of X"), the concept order, the data shapes and stdin/stdout conventions (Go/JS), and for CLI courses the conventions the checks rely on: workspace folder, resource name prefixes, commit labels, ports, the copy-paste sample files and what they do. Every later lesson follows this file. When rewriting an existing course, back the old one up outside `content/` first and list every concept it taught, so none is lost.
4. **Outline first.** Produce a chapter-by-chapter table: chapter id + title (named after the concept), 5-12 lessons each with id, type, language/runner, one-line goal, xp, then the boss. Get the user's OK before writing files.
5. **Write one chapter at a time**: all files for every lesson, then `chapter.yaml`, then update `course.yaml`.
6. **Verify.** If a shell and the repo are available: `make validate` (or `cd api && go run ./cmd/cq validate ../content`) until clean, then `cd api && go run ./cmd/cq test ../content/<course>` (Go/JS lessons need `make wasm-modules` once). For `project_cli` lessons, record answers from real runs (see "Recording `solution.json`"). Fix every `[solution]: fail` and `[wrong/x]: pass`. Otherwise, self-check against the checklist at the end.
7. **Finish.** Grep the course for leftover story words when rewriting. Update `docs/SPEC.md` with the decision (new lesson ids reset progress for that course), run `make test`, and don't commit unless asked. To try every lesson locally, the dev user (`dev@codequest.local`, with dev auth on) has every chapter unlocked.

---

## The voice (boot.dev style, adapted)

A lesson is **one idea**, read in under two minutes, then practised immediately.

**Shape of every `prompt.md`:**

```markdown
# <Short noun-phrase title>

<1-3 short paragraphs: the problem first, then the concept that solves it.>

```python
<smallest example that shows the idea, with output as a comment>
print(area_of_square(4))
# 16
```

<Optional: a numbered walk-through of what runs, in order.>

## Assignment

<Exact task: function name, inputs, return value, exact strings/keys.>

1. <step or checklist item>
2. <step or checklist item>

<A worked example: given X, return Y.>
```

**Rules of the voice:**

- **Problem before syntax.** Show the pain (copy-pasted code, a crashing loop, a leaked password) and *then* the feature. "What if we had a thousand of these?" is the hook.
- **Second person.** "You'll", "your project", "let's". No role-play: the learner is just the learner.
- **Plain words, jargon defined inline.** First use gets a gloss: *"returns (outputs) a value"*, *"call (fancy programmer speak for 'use') the function"*.
- **Light humour and opinions, never at the learner's expense.** One wry line per lesson at most. Mild opinions are fine when they teach ("rebase isn't dangerous; misusing it is").
- **Normalise struggle at hard spots.** "This takes a minute to click. Slowing down here is normal."
- **Short sentences, short paragraphs, no filler.** No "In this lesson we will learn...". Start with the idea.
- **Code beats prose.** Every concept gets a runnable snippet in the course's language; show output as comment lines (`# 16`, `// 16`).
- **Assignments are exact.** Name the function, every parameter, every dict key / string / column / branch / file, with one concrete input → output example. Ambiguity is a bug.
- **No story framing.** The assignment starts with the task. If a reason helps, give a real one ("so a teammate can see what changed"), never a plot.
- **Reuse across lessons.** Build on the same function, table or repo over several lessons so the project feels continuous.
- **Bug-hunt lessons are frequent.** "The total comes out one item short. Fix line 3." Starter contains the bug. In tool courses the bug is in the tool's config (a Dockerfile, a compose file, a Git setting), never in app code.
- **Scaffold comments** in starters when helpful: `# don't touch above this line` / `// don't touch below this line`.
- **Review quizzes** (`mcq`) follow dense concept lessons: trace the order of execution, predict output, spot the bug, true/false on a common misconception.
- **Scope out loudly.** When a detail is out of scope, say so in one line ("We won't cover how SHA-1 works here; just know commit ids are hashes"), so the learner doesn't think they missed something.

**Lesson rhythm inside a chapter:** concept → practice → (quiz) → twist/edge case → bug-fix → practice → … → boss. 5-12 lessons, then the boss.

---

## Choosing where lessons run

| Topic needs | Use |
| --- | --- |
| Standard-library Python | `code`, `language: python`, `runner: pyodide` (Run in browser) |
| Python packages (`cryptography`, LangChain fakes), bigger limits | `code`, `runner: docker-python` (Run on server) |
| Go or JavaScript | `code`, `language: go \| js`, `runner: wasm` (Run on server; see "Go and JS lessons") |
| SQL | `sql` |
| Git, shell, Docker, HTTP servers, anything on the learner's machine | `project_cli` |
| Ideas, misconceptions, "what does this print", compiler rules the runner can't enforce | `mcq` |
| Another language (Rust, C, …), or Go/JS needing a package outside the allowlist, `net`, `os/exec`, goroutine-heavy servers, npm | Not gradable as `code`. Options: teach it in a supported language, use `project_cli` (the learner's own toolchain, e.g. `go test -count=1 -run TestX -v ./pkg` checking `--- PASS: TestX` and no `FAIL`), or ask for a new wasm language module. Raise this before outlining; never silently switch languages. |

`project_cli` grades output the learner's CLI reports, so it can be forged. That's an accepted platform decision; don't build a course's only proof of skill on it if a `code` lesson could test the same thing.

---

## Go and JS lessons (`runner: wasm`)

Go runs in the yaegi interpreter, JS in goja, both compiled to WASI and run on the server. The grading model is **program in, text out**, not unit tests, which changes how lessons are built.

**Grading tools:** `io_cases` (stdin → expected stdout, one row per case, at most one per exercise, 1-50 cases), `stdout_equals`, `stdout_regex`, `exit_code`. `unit_tests` and `forbidden_calls` are Python-only, so "don't use X" can only be a hint, never a check.

**Keep the boot.dev "write a function" feel with a harness `main`.** The learner edits one function at the top. Below a `// don't touch below this line` comment, the starter's `main` reads stdin, parses it, calls the function and prints the result in a fixed format. Cases then test the function through stdin. Keep the input format trivial (one value or one record per line) and identical across a chapter, recorded in `.drafts/plan.md`.

**Design rules:**

- Every case's stdout must be fully deterministic. No timestamps, map iteration order (sort keys first in `main`), random numbers or pointer addresses. For crypto with random nonces, have `main` print only derived facts (`roundtrip: ok`, `tampered: rejected`), or take the nonce from stdin.
- Test properties through `main`: print both `encrypt(x)` and `decrypt(encrypt(x)) == x`, or wrong-key and tamper results, so each case checks several things.
- Visible cases are named after what they feed in (`"stdin: 5 7"`). Hidden cases are named in words ("handles an empty message") with `visible: false`. Their stdin never reaches the browser and hidden rows never show got/expected.
- Hidden cases cover edges the visible ones don't: empty input, zero, negatives, very long input, the case that breaks the obvious-but-wrong approach.
- Keep cases to roughly 3-10. Each one is a fresh module run (about 40 ms for Go, 20 ms for JS), and the whole job shares 2 × `timeout_ms`.
- Use `exit_code` for error-handling lessons ("exit 1 on bad input"); it runs the program with empty stdin.
- **Go limits:** only allowlisted standard-library packages (check `runners/wasm/symbols/`; `encoding/json` is not available); `net`, `net/http`, `os/exec`, `syscall`, `unsafe` fail. yaegi does not enforce `go build` rules (unused variables or imports compile fine), so teach those with `mcq`, not with a code check. Panic messages differ from the compiler's ("reflect: slice index out of range"), so never match on panic text.
- **JS limits:** ES2015+, Node-like `console`, `process.exit`, `process.stdout.write`, stdin via `require("fs").readFileSync(0, "utf8")`. No modules, npm, `fetch` or other `require`s. Output formatting matches Node only for common values, so print strings and numbers explicitly instead of logging objects or arrays.
- Run goes to the server (rate-limited), so prompts should encourage reading and reasoning before hammering Run.
- `wrong/` holds `.go` / `.js` files: realistic mistakes that each fail at least one case.

**Starter rule:** the starter compiles, runs and fails the cases (the function returns a zero value: `return 0`, `return ""`, `return nil`, `return []`). Never a body that doesn't compile.

---

## Advanced topics (security, cryptography, systems)

Patterns observed in boot.dev's cryptography course, adapted:

- **Chapter arc: broken → attack → primitive → real thing → applications.** E.g. what encryption is → encoding (bytes, hex, base64) → brute force → a toy cipher (Caesar) → XOR / one-time pad → stream and block ciphers → an old standard (DES) → the modern one (AES) → asymmetric (RSA) → hashing → signatures → key derivation. Teach the weak version and the attack on it first; the strong version then has an obvious reason to exist.
- **Language choice.** In Go, confirm `crypto/*`, `encoding/hex` and `encoding/base64` are in the allowlist before outlining; anything missing goes to Python (`hashlib`, `hmac`, `secrets` in the browser; `cryptography` on `docker-python`).
- **Prerequisites get their own chapter.** If later lessons need bytes/hex/bit math, teach that before the first cipher.
- **Use it before you understand it.** The first lesson can call provided black-box helpers (`encrypt()`, `decrypt()`) to show the end goal, with "we'll unpack keys and IVs later". Put helpers below the `don't touch` line.
- **Tiny pure functions that compound.** `alphabetSize(bits)` → `keyspace(bits, length)` → `timeToCrack(...)`. Each lesson adds one function the next one uses.
- **Interactive intuition → diagrams.** Where boot.dev uses widgets (bit tables, rebase animations), use an SVG in `assets/` with real alt text, or a worked table in Markdown (bits in, bits out).
- **Worked examples in the raw form** (binary strings, hex, byte lists), short enough to check by hand.
- **"Note on design" asides**, one short paragraph, explaining a real engineering choice (why one `crypt()` serves both encrypt and decrypt).
- **Test properties, not just values.** Round-trip, wrong key fails, same input + same nonce is deterministic, different nonce changes the output, tampering is detected.
- **Don't test what you didn't teach.** Every API, constant or edge case the tests need appears in the prompt or an earlier lesson, with the exact function and package. If the learner must read docs, say which function in which package.
- **Safety framing.** Toy ciphers are for understanding. Say once per course and again at the first real algorithm: never ship hand-rolled crypto; production code uses vetted libraries. Attack lessons target only data the lesson provides.
- **Bosses are realistic features** ("Vault Unlock": derive a key from a password with a salt, encrypt, verify integrity), combining only primitives from the chapter.

---

## CLI-graded courses (`project_cli`, e.g. Git, Linux, Docker, local servers)

Patterns observed in boot.dev's Git course, which is graded by running commands on the learner's machine:

- **Lesson 1 is "What is X?"** (an `mcq`: the problem it solves, in plain words). **Lesson 2 installs it and proves the pipeline**: install the tool and the `codequest` CLI, create the course folder, then trivial checks (tool version is new enough, `pwd` contains the course folder). Nothing else can go wrong yet.
- **One long-lived workspace per course.** The learner creates one repo/folder early and every later lesson builds on it. `.drafts/plan.md` lists its expected state after each lesson.
- **The learner mutates; the steps only inspect.** Assignment steps tell them what to run (`git switch -c ...`, edit a file, commit). The graded `cli_command`s are read-only checks: `git --no-pager log --oneline`, `git branch`, `git config --get user.name`, `cat file`, `curl`.
- **Make state checkable with labels.** Learners' hashes, timestamps and paths differ, so never check those. Instead, have commit messages start with labels (`A:`, `B:`, `C:`...) and check with `stdout_contains_all` / `stdout_contains_none` ("contains D:, does not contain E:" proves where a branch starts). Use `captures` when a later step genuinely needs a value from an earlier one.
- **Disable pagers and colour**: `git --no-pager ...`, `--no-color`, `-n` limits, so output is stable.
- **Name each step after what it proves**: "Branch starts at commit D", "update-docs is the active branch".
- **Numbered, exact instructions with a self-check** before submitting: "Run `git log --oneline -n 1` and confirm D is on top." Placeholders in caps (`COMMITHASH`) with how to get them.
- **Terminal + quiz combo** for things that vary per learner: "Run these commands, look at `.git/objects`, then pick the option that matches the *shape* of your path" as an `mcq` whose options are format patterns.
- **Concept lessons are `mcq`** (with a diagram) between hands-on lessons, e.g. merge vs rebase.
- **Every test is visible**, so the prompt states exactly what will be checked. Rely on several independent inspections rather than one stdout match.
- **Stay inside the lesson folder.** No `sudo`, no deleting outside the workspace, no global config changes beyond what the lesson explains (prefer `--local` config). Set `allowed_os` when commands differ on Windows. The CLI shows the commands and asks before running them, so keep them short and readable.

---

## Beginner tool courses (Docker, Git, Kubernetes, Linux…)

These rules come from rewriting the Docker and Git courses (see `content/docker/.drafts/plan.md` and `content/git/.drafts/plan.md` for working examples).

- **Concept order, boot.dev style:** what it is and why → install → the parts (e.g. engine, image, container, registry/Hub; or working tree, staging area, commits) → first real use → one command per lesson → internals → collaboration/production. Quizzes (`mcq`) sit after each dense idea ("VMs vs containers", "image vs container", "the three areas").
- **Stay on the tool. No language barrier.** The learner should never need to know Go, Python or JS. If an app is needed (something to containerise, something to commit), it is given as **one copy-paste block** that creates every file, with "you don't need to read this":

  ```sh
  mkdir -p app
  cat > app/main.go <<'EOF'
  // a tiny web server for this course. You don't need to read this.
  ...
  EOF
  ```

  Describe what the app does in a small table (address → answer) instead of explaining its code. Never ask the learner to edit app code; changes go through the tool (env vars, flags, config files, Dockerfile lines).
- **Minimal setup.** Prefer data that needs no code at all: a Git project made of Markdown pages, Docker lessons on `alpine`, `nginx`, `hello-world`, `busybox`. Every file content the learner types comes from a copy-paste block (`printf`/`cat <<'EOF'`), so typos can't fail a check.
- **Show config files with a comment on every line** the first time (`# start from an image that has the Go tools`). Point out the pattern ("pick a base, copy files in, build, say how to start").
- **Naming conventions the checks rely on**, recorded in `plan.md`: one workspace (`~/docker-course`, `~/git-course/project`); resources prefixed (`cq-` containers/volumes/networks, `cq/` images); a unique host port per lesson (8081-8099); commit labels per chapter (`A1:`, `B1:` … `F8:`).
- **Never make the learner run destructive commands on their whole machine.** No `docker system prune`, `docker container prune`, `rm -rf ~`, `git config --system`. Teach cleanup on the course's own named resources (`docker rm cq-web`). `--global` Git config is OK only in the setup lesson that explains it (name, email, `init.defaultBranch main`).
- **Remote services** (GitHub, Docker Hub): one lesson, a free account, clearly optional-looking setup steps. Its `solution.json` is written by hand from a real run's shape; never push to a real account while authoring.
- **Platform tips** in one line where real learners trip: Docker Desktop file sharing for bind mounts, `docker` vs `docker compose` plugin, macOS `sed -i ''`.

## Check-writing pitfalls (`project_cli`)

- **Inspect, don't parse prose.** Use `docker inspect -f '{{.State.Status}}' cq-web`, `--format '{{.Service}}={{.State}}'` and similar to print stable `key=value` lines, then `stdout_contains_all: ["api=running"]`.
- **Git:** find a commit by label with `git log -1 --format=%H ':/^A1:'` and `captures: [{ name: a1, regex: "([0-9a-f]{40})" }]`; later steps use `${a1}`. `:/regex` swallows the rest of the revision, so `':/^B1:about.md'` is wrong: capture first, then `${b1}:about.md` (or `git ls-tree`). Check parents with `--format='parents=%P'`.
- **`${name}` is reserved for captures** (plus `${baseURL}`). Shell variables in commands use `$VAR` or `$(...)`.
- **Clean tree:** `git status --short; test -z "$(git status --porcelain)"` with `exit_code: 0` (the output tells the learner what's left).
- **Version gates:** print `tool=ok` / `tool=too-old X` with `awk` instead of matching a version string.
- **Never check ids, digests, sizes, timestamps or absolute paths literally.** Compare them in shell (`test "$(...)" = "$(...)"`) and print a word.
- **Every wrong answer must really fail.** Run `cq test`; if a `wrong/` file passes, the check is too loose. Add `stdout_contains_none` for the mistake (e.g. a container that was started with dropped args must not show `listening`).
- **YAML:** quote any command containing `: `, `#`, `{{` or leading `*`. Prefer `"..."` with `\"` inside, or `'...'` around Go templates.
- **Dockerfile/compose snippets:** no inline `#` comments after an instruction (Docker treats them as arguments); put comments on their own line.
- **30 s per command:** the learner does builds, pulls and pushes; steps only inspect.

## Recording `solution.json` and `wrong/` for CLI lessons

Never hand-write CLI answers (except the GitHub/remote-account case above). Do the lesson for real and record what the real step runner prints:

1. Make a throwaway workspace (in the scratchpad, or under the repo if Docker Desktop must bind-mount it, removed afterwards). For Git, isolate the user: `HOME=<scratch>/ghome GIT_CONFIG_NOSYSTEM=1 GIT_TERMINAL_PROMPT=0 GIT_EDITOR=true`, so the user's real `~/.gitconfig` is never touched.
2. Write a small "do" script per lesson: the exact commands from the assignment, so the workspace reaches the lesson's end state. Run lessons in order; each builds on the last.
3. Record with the real runner. Copy `scripts/zz_record_test.go` from this skill into `api/cmd/codequest/` (it uses the CLI's own `newStepRunner`, so captures, stdout+stderr and HTTP match exactly), build it once with `cd api && go test -c -o <scratch>/rec ./cmd/codequest`, then per lesson:
   `CQ_REC_YAML=<lesson>/exercise.yaml CQ_REC_DIR=<workspace> CQ_REC_OUT=<lesson>/solution.json <scratch>/rec -test.run TestZZRecord`
4. Make `wrong/*.json` from the solution with `scripts/wrong.py <lesson> <name> <step>=<exit>:<stdout>` (HTTP: `<step>=h<status>:<body>`): realistic mistakes such as `not_started`, `wrong_folder`, `wrong_branch`, `forgot_commit`. Use stdout a real mistaken run would print.
5. `cd api && go run ./cmd/cq test ../content/<course>` until every reference passes and every wrong answer fails.
6. **Delete `api/cmd/codequest/zz_record_test.go`** and the throwaway resources (only ones you created, by exact name; list them first).

## How prompts render

- The page already shows the lesson title from `exercise.yaml`; the prompt's leading `# Title` line is dropped by the web app, so keep it (it's useful in the file) but don't repeat the title in the first paragraph.
- `##`/`###` headings render as headings: use `## Assignment`, and `## Tip` sparingly.
- Every fenced code block gets a **Copy** button. Put each copy-paste unit in its own block (one block per step), with no `$ ` prompts and no output mixed in; show output in a separate block or as `# comments`.
- Back-to-back ` ```cpp ` and ` ```python ` blocks render as one tabbed block (C++ first, choice remembered). Use it when a lesson shows the same code in both languages.
- Links render underlined in the accent colour and external ones open in a new tab.
- Tables and `**bold**` render (GFM). Images go in `assets/` as SVG that read well in light and dark themes.

---

## Mapping boot.dev patterns to CodeQuest

| boot.dev | CodeQuest |
| --- | --- |
| `run_cases` (shown) | Python: `tests_visible.py` with `expect()`. Go/JS: visible `io_cases` |
| `submit_cases` (extra on submit) | Python: `tests_hidden.py`, named in words, with `input:` line. Go/JS: `io_cases` with `visible: false` |
| `main_test.go` calling the function | Go/JS: the harness `main` below `// don't touch below this line` |
| "Tips" section | `hints:` list: nudges only, revealed one at a time |
| Concept quiz | `mcq` lesson, `xp: 5`, `explain` that teaches |
| CLI lesson (`bootdev run`) | `project_cli` with `cli_command` / `http_request` steps |
| Interactive widget | SVG in `assets/` or a worked Markdown table |
| Guided project / end of chapter | `boss: true` `code` or `project_cli` lesson, project-like title, 5-8 HP |
| Starter with `pass` or `// ?` | **Not allowed.** Starter compiles/runs and returns a zero value so checks fail cleanly |

---

## File templates

### `course.yaml`

```yaml
id: <slug>
title: <Title>
description: "Write CodeQuest courses, chapters and lessons (mcq, code in Python/Go/JS, sql, project_cli, bosses) in a boot.dev-style voice with no story, built from zero for total beginners, as valid content/ files that pass cq validate and cq test. Covers beginner tool courses (Docker, Git, Kubernetes) with copy-paste setup, CLI check design and recording real CLI answers."
  <One or two sentences: what you build, who it's for.>
order: <n>
published: true
chapters:
  - 01-<chapter>
icon: <1-4 chars>
level: beginner
tags: [<max 4, each ≤ 20 chars>]
topics: <≤ 60 chars>
```

### `chapter.yaml`

```yaml
id: <course-short>.<concept>
title: <Concept name, e.g. Loops>
exercises:
  - 01-<lesson>
  - 0N-<boss>
```

### Python code lesson

`exercise.yaml`:

```yaml
id: <chapter-id>.<lesson-slug>
type: code
version: 1
title: <Short Title>
xp: 20
difficulty: 1
tags: [<concept>]
free: <true for free chapters>
spec:
  language: python
  runner: pyodide
  run: browser
  prompt: prompt.md
  entry: main.py
  starter: starter.py
  solution: solution.py
  examples:
    - { input: "<call args>", output: "<result>", explanation: "<why>" }
  hints:
    - <a question pointing at the concept, never the fix>
  checks:
    - kind: unit_tests
      file: tests_visible.py
      visible: true
    - kind: unit_tests
      file: tests_hidden.py
      weight: 2
```

`tests_visible.py`:

```python
from cq import expect
from main import get_crew_record


def test_basic():
    "get_crew_record('vega', 'deck2', 3)"
    expect(get_crew_record("vega", "deck2", 3), {"name": "vega", "deck": "deck2", "level": 3, "id": "vega@deck2"})
```

`tests_hidden.py`:

```python
from main import get_crew_record


def test_level_zero():
    """keeps a level of zero

    input: get_crew_record('orin', 'bay1', 0)
    """
    assert get_crew_record("orin", "bay1", 0)["level"] == 0
```

For crypto-style lessons, hidden tests add properties: round-trip, wrong key, empty input, input longer than the key, input not mutated.

`wrong/`: 2-3 files, each a realistic learner mistake (off-by-one, mutating the input, wrong separator, reusing the key/nonce, the banned call). Each must fail at least one check.

### Go code lesson (`runner: wasm`)

`exercise.yaml`:

```yaml
id: <chapter-id>.<lesson-slug>
type: code
version: 1
title: <Short Title>
xp: 20
tags: [<concept>]
spec:
  language: go
  runner: wasm
  prompt: prompt.md
  entry: main.go
  starter: starter.go
  solution: solution.go
  examples:
    - { input: "xorBytes([0x0f], [0xff])", output: "[0xf0]", explanation: "each bit flips where the key has a 1" }
  hints:
    - <nudge>
  checks:
    - kind: io_cases
      cases:
        - { name: "stdin: 0f ff", stdin: "0f ff\n", stdout: "f0\n" }
        - { name: "stdin: 68656c6c6f 6b6b6b6b6b", stdin: "68656c6c6f 6b6b6b6b6b\n", stdout: "030e070704\n" }
        - { name: empty message, stdin: " \n", stdout: "\n", visible: false }
        - { name: key equal to message gives zeros, stdin: "abab abab\n", stdout: "0000\n", visible: false }
```

`starter.go`:

```go
package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func xorBytes(data, key []byte) []byte {
	return []byte{}
}

// don't touch below this line

func main() {
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	parts := strings.Fields(line)
	data, key := []byte{}, []byte{}
	if len(parts) == 2 {
		data, _ = hex.DecodeString(parts[0])
		key, _ = hex.DecodeString(parts[1])
	}
	fmt.Println(hex.EncodeToString(xorBytes(data, key)))
}
```

Check the package list against the allowlist, and run the solution and every case locally with `cq test` before writing the next lesson.

### JS code lesson

Same shape with `language: js`, `entry: main.js`, `starter.js`, `solution.js`, `wrong/*.js`. The harness reads stdin with `require("fs").readFileSync(0, "utf8")` and prints with `console.log` of a string or number.

### Quiz lesson (`mcq`)

```yaml
id: <chapter-id>.<slug>-quiz
type: mcq
version: 1
title: <Topic> Review
xp: 5
spec:
  question: prompt.md
  multiple: false
  options:
    - { id: a, text: "<plausible mistake>" }
    - { id: b, text: "<correct>" }
    - { id: c, text: "<plausible mistake>" }
  answer: [b]
  explain: <Why b is right AND why the tempting wrong one is wrong.>
```

The `prompt.md` can carry a short recap (a numbered trace of what runs, or terminal steps to try first) before the question.

### CLI lesson (`project_cli`)

```yaml
id: <chapter-id>.<slug>
type: project_cli
version: 1
title: <Short Title>
xp: 20
spec:
  prompt: prompt.md
  allowed_os: [linux, darwin]
  solution: solution.json
  steps:
    - description: Branch starts at commit D
      cli_command:
        command: git --no-pager log --oneline -n 1
        tests:
          - exit_code: 0
          - stdout_contains_all: ["D:"]
          - stdout_contains_none: ["E:"]
    - description: update-docs is the active branch
      cli_command:
        command: git branch --no-color
        tests:
          - exit_code: 0
          - stdout_contains_all: ["* update-docs"]
```

`solution.json` is a recorded passing run (one entry per step, in order); `wrong/*.json` are recorded failing runs, e.g. branched from the wrong commit, or on the wrong branch:

```json
{"steps": [{"cli": {"exit_code": 0, "stdout": "3f2a1bc D: add chapter list\n"}}, {"cli": {"exit_code": 0, "stdout": "  main\n* update-docs\n"}}]}
```

### SQL

Follow `CONTENT.md` section 8: visible check has `expected`; hidden check uses `hidden_data.sql` with ties, NULLs and boundaries and `visible: false`; `ordered: true` only if the prompt asks for an order.

### Boss

A `code` lesson with `boss: true`, last in the chapter, title like a shipped project ("Cargo Manifest CLI"), xp 50+, 5-8 total check rows (visible + hidden tests, or io cases, plus any `exit_code`/stdout check), combining only skills from this chapter. Its prompt is a clear spec with several examples (no story ticket). Go and JS bosses work the same way: each io case is one HP.

A `project_cli` lesson can also be a boss (`boss: true`): one HP per step test, so 5-8 tests across its steps. Tool courses (Git, Docker) end every chapter this way: a short checklist that combines the chapter's commands, e.g. "Under the Hood", "Ship It". Nothing new is taught in a boss.

---

## Hard rules (from CONTENT.md and SPEC.md)

- Ids are lowercase, globally unique, and never change after publishing.
- Every folder listed in a manifest exists, and every folder exists in a manifest. Unknown fields fail the build.
- XP: quiz 5, small lesson 15-20, bigger 25-30, boss 50+. Don't inflate; bonuses are automatic.
- Mark every lesson in the free chapters `free: true`.
- Python visible tests: `expect()`, docstring named after the call. Hidden tests: named in words (no `()[]{}=`), with an `input:` line that never contains the expected value.
- Hidden tests or cases cover what visible ones don't: empty input, negatives, zero, mutation, big input, type edges.
- Starter runs (and compiles) without error and fails the checks.
- Browser Python uses only features shared by 3.12 and 3.14, standard library only. Anything else: `runner: docker-python`.
- Go/JS: `runner: wasm`, `io_cases` at most once per exercise with 1-50 cases, no `unit_tests` or `forbidden_calls`, allowlisted packages only, deterministic output.
- Never leak the answer in `prompt.md`, examples, hints, visible test or case names, images or `input:` lines. Hints point at a concept or a line ("What should `total` start at?").
- `project_cli`: 1-20 steps, 30 s per command, 10 s per HTTP request; stdout means stdout+stderr; all tests visible.
- Images: in the lesson's flat `assets/`, real alt text, readable in dark and light themes, never drawing the answer.
- Changing grading of a published lesson: bump `version`. Typos: just edit.

## Self-check before handing over

- [ ] Prompt opens with the problem, has a runnable example, and an `## Assignment` with exact names and a worked example.
- [ ] One concept per lesson; no story, characters or tickets; every term defined on first use; out-of-scope details are named as such.
- [ ] Tool courses: no app code to write or read; every file the learner needs comes from one copy-paste block.
- [ ] Everything the tests need was taught in this or an earlier lesson.
- [ ] Starter fails cleanly; solution passes; each `wrong/` file fails.
- [ ] Go/JS: harness `main` is below the don't-touch line, output is deterministic, packages are allowlisted, hidden cases are named in words.
- [ ] CLI checks are read-only, pager-free, and check labels rather than hashes or timestamps.
- [ ] Visible/hidden naming rules followed; no answer leaks anywhere public.
- [ ] `free`, `xp`, `version`, ids and manifest lists are right.
- [ ] Chapter ends with a 5-8 HP boss.
- [ ] `make validate` and `cq test` pass (when a shell is available).