# Writing CodeQuest content

This guide is for anyone writing courses for CodeQuest. Courses are plain files in `content/`, versioned in Git. You never touch platform code: you write YAML manifests, Markdown prompts, starter code, a reference solution, tests and a few known-wrong answers. The build step checks all of it and loads it into the database.

`docs/SPEC.md` (section "Content format") is the source of truth. If the two ever disagree, the spec wins; please fix this guide in the same change.

Contents:

1. [How content is organised](#1-how-content-is-organised)
2. [The authoring loop](#2-the-authoring-loop)
3. [Create a course](#3-create-a-course)
4. [Create a chapter](#4-create-a-chapter)
5. [The exercise header (every lesson type)](#5-the-exercise-header-every-lesson-type)
6. [Lesson type: `mcq` (quiz)](#6-lesson-type-mcq-quiz)
7. [Lesson type: `code` (Python, Go, JavaScript)](#7-lesson-type-code-python-go-javascript)
8. [Lesson type: `sql`](#8-lesson-type-sql)
9. [Lesson type: `project_cli` (on the learner's machine)](#9-lesson-type-project_cli-on-the-learners-machine)
10. [Boss fights](#10-boss-fights)
11. [Images](#11-images)
12. [Career paths](#12-career-paths)
13. [Rules that keep answers secret](#13-rules-that-keep-answers-secret)
14. [Changing published content](#14-changing-published-content)
15. [Checklist before you open a PR](#15-checklist-before-you-open-a-pr)
16. [Common errors](#16-common-errors)

`content/feature-tour/` has one working lesson of every type. Copy from it whenever you're unsure.

---

## 1. How content is organised

```
content/
  paths.yaml                      # optional: career paths across courses
  python-basics/                  # one folder per course
    course.yaml
    01-variables/                 # one folder per chapter
      chapter.yaml
      01-hello/                   # one folder per lesson ("exercise")
        exercise.yaml
        prompt.md
        starter.py
        solution.py               # server-only
        tests_visible.py
        tests_hidden.py           # server-only
        wrong/                    # server-only: answers that must fail
          off_by_one.py
        assets/                   # optional: public images
          diagram.svg
```

- **Order comes from the manifests, not folder names.** `course.yaml` lists chapter folders and `chapter.yaml` lists lesson folders, in order. Number the folders anyway (`01-`, `02-`) so they sort sensibly in your editor.
- **Every folder must be listed, and every listed folder must exist.** A stray folder is an error, not something that gets ignored. Hidden folders (`.drafts/`) are ignored.
- **Unknown fields are errors everywhere.** A typo like `titel:` fails the build, so mistakes can't slip through silently.
- **Ids are forever.** Learners' progress and submissions are keyed on `id`. You can rename a folder or a title, but never change an `id` once a lesson is published.

## 2. The authoring loop

Run these from the repository root.

| Command | What it does | Needs |
| --- | --- | --- |
| `make validate` | Checks every manifest, file and image in `content/`. Prints each problem as `path: reason`. | Go only |
| `cd api && go run ./cmd/cq validate ../content/my-course` | Same, for one course | Go only |
| `make content-test` | Content CI: grades every reference solution (must pass) and every `wrong/` answer (must fail) on the real grader | Docker (builds the `codequest/python:3.12` image) and Go (builds the Go/JS modules) |
| `cd api && go run ./cmd/cq test ../content/my-course` | Same, for one course (run `make runner-image` once first, and `make wasm-modules` for Go or JS lessons) | Docker; Go/JS lessons need only Go |
| `make content` | Compiles `content/` into Postgres so the running app shows your changes | Docker Compose database |
| `make e2e` | Includes a browser check that each solution passes its visible tests on **Run** | Docker |

A typical loop:

1. `make up` (once) and sign in with the dev sign-in at http://localhost:5173.
2. Write or edit files.
3. `make validate` until it's clean.
4. `make content` and open the lesson in the browser. Try the starter, a wrong answer and the solution.
5. `make content-test` before you open a PR. CI runs it too.

## 3. Create a course

`content/<course-id>/course.yaml`:

```yaml
id: python-basics            # lowercase slug, unique, never changes
title: Python Basics         # required
description: >
  Variables to a tested command-line app. The best first step for everyone.
order: 1                     # position on the Courses page (ties sort by id)
published: true              # default true; false hides the course
chapters:                    # chapter folder names, in order
  - 01-variables
  - 02-conditionals

# Catalog card (all optional)
icon: py                     # 1-4 lowercase letters/digits shown in the tile; default: first 2 of the id
level: beginner              # beginner (default) | intermediate | advanced
tags: [Python]               # Courses page filters: at most 4, each 1-20 characters
topics: Variables, conditionals, loops   # one line on the card, at most 60 characters
```

You don't set where lessons run or whether the course is "Pro". The Courses page works those out from the lessons:

- **Runs in browser** shows when any lesson is not `project_cli`; **CodeQuest CLI** shows when any lesson is `project_cli`.
- **Free to start** shows when at least one lesson is `free: true`. With no free lessons the card says **Pro** and offers a preview.

## 4. Create a chapter

`content/<course>/<chapter-folder>/chapter.yaml`:

```yaml
id: py.variables             # lowercase words joined by . or -, unique across all content
title: Variables and types   # required; shown as a skill on the learner's profile
exercises:                   # lesson folder names, in order
  - 01-hello
  - 02-types-quiz
  - 03-tip-calculator        # end the chapter with a boss (see section 10)
```

Keep chapters short: 5 to 12 lessons, then a boss. A chapter's title is what the profile's Skills card shows, so name it after the concept ("Loops", not "Chapter 3").

## 5. The exercise header (every lesson type)

Every `exercise.yaml` starts with the same header, followed by a type-specific `spec:`.

```yaml
id: py.variables.hello       # required, unique, never changes
type: code                   # required: mcq | code | sql | project_cli
version: 1                   # required, 1 or higher; bump when grading changes (section 14)
title: Say hello             # required
xp: 20                       # required, 0 or more; base XP for the first pass
difficulty: 1                # 1-5, default 1 (1-2 shows as "Easy")
tags: [strings]              # default []; shown on the course map
free: true                   # default false; free lessons work without a subscription
boss: false                  # default false; only code and project_cli lessons can be bosses
spec:
  ...
```

**XP guide.** Quiz 5, small code or SQL lesson 15-20, bigger lesson 25-30, boss 50+. The game adds bonuses itself (+50% first try, +25% no hint), so don't inflate the base.

**Free tier.** The rule is "the first chapters of every course are free". Mark every lesson in those chapters `free: true`.

## 6. Lesson type: `mcq` (quiz)

```
02-types-quiz/
  exercise.yaml
  prompt.md        # the question, in Markdown
```

```yaml
spec:
  question: prompt.md
  multiple: false                # true: "pick every answer", with checkboxes
  options:                       # at least 2
    - { id: a, text: "[1, 2, 3]" }
    - { id: b, text: "None" }
    - { id: c, text: "It raises an error" }
  answer: [b]                    # option ids; exactly one when multiple is false
  explain: list.sort() sorts in place and returns None. Use sorted(nums) for a new list.
```

- Option `id`s: 1-32 characters from `a-z 0-9 _ -`, unique within the question.
- `explain` is required. It shows after **every** submit, pass or fail, so write it to teach, not just to confirm.
- `answer` and `explain` never reach the browser before a submit.
- An option can carry a picture: `{ id: a, text: "A stack", image: assets/stack.svg }`. `text` is still required and becomes the image's alt text (section 11).
- `cq test` checks your answer passes and that each single wrong option fails. You don't write `wrong/` for quizzes.

Good quiz options are plausible mistakes. "It raises an error" is a useful wrong answer when learners really do expect one.

## 7. Lesson type: `code` (Python, Go, JavaScript)

```
03-double-it/
  exercise.yaml
  prompt.md
  starter.py          # what the editor opens with
  solution.py         # reference solution (server-only)
  tests_visible.py    # shown to the learner; runs on Run and Submit
  tests_hidden.py     # server-only; runs on Submit
  wrong/              # known-bad answers; each must fail at least one check
    in_place.py
    uses_map.py
```

```yaml
spec:
  language: python
  runner: pyodide            # pyodide (graded on the server too) | docker-python (server only)
  run: browser               # browser | server | off  (where the Run button executes)
  prompt: prompt.md
  entry: main.py             # the module tests import from; the learner's file is saved under this name
  starter: starter.py
  solution: solution.py
  limits: { timeout_ms: 3000, memory_mb: 128 }   # defaults; max 10000 ms / 512 MB
  examples:                  # optional, shown as "Example n: Input / Output / Explanation"
    - { input: "nums = [1, 2, 3]", output: "[2, 4, 6]", explanation: "each number times 2" }
  constraints:               # optional
    - "0 ≤ len(nums) ≤ 10⁴"
  hints:                     # optional, revealed one at a time
    - Start with an empty list and append to it.
  checks:                    # at least one
    - kind: unit_tests
      file: tests_visible.py
      visible: true
    - kind: unit_tests
      file: tests_hidden.py
      weight: 2
    - kind: forbidden_calls
      names: [map]
      message: Try it with a loop instead of map()
```

### Where the code runs

| `runner` | Default `run` | Use it for |
| --- | --- | --- |
| `pyodide` | `browser` | Standard-library Python. **Run** is instant in the browser; **Submit** is re-graded on the server. |
| `docker-python` | `server` | Anything the browser can't do: LangChain (`cq.fakes`), pinned packages, larger limits. **Run** goes to the grading worker (visible checks only, not stored). |

The browser runs Python 3.14 (Pyodide) and the server runs Python 3.12, so keep browser lessons to features both versions share. `make e2e` fails if a solution doesn't pass its visible tests in the browser.

### Writing tests

Tests are plain Python modules with top-level `test_*` functions that import from the entry module.

```python
# tests_visible.py
from cq import expect
from main import double_all


def test_three():
    "double_all([1, 2, 3])"                # the docstring is the row name the learner sees
    expect(double_all([1, 2, 3]), [2, 4, 6])
```

- **Use `expect(got, expected)` in visible tests.** It compares exact types (so a rigged `__eq__` can't pass) and shows the learner "got … expected …" on failure.
- **Name visible tests after the call** (`double_all([1, 2, 3])`). The learner sees that name next to the result.
- **Name hidden tests in words** ("doubles negative numbers"). A hidden name that looks like code (contains `()[]{}=`) is replaced by "Hidden test N". Hidden rows never show got/expected values or exception details, only your message, "Wrong result" or "Your code raised <ExceptionType>".
- **Give hidden tests an `input:` line** so the Peek power-up can reveal what was called, and nothing more:

  ```python
  def test_negative():
      """doubles negative numbers

      input: double_all([-4, 5])
      """
      assert double_all([-4, 5]) == [-8, 10]
  ```

  Never put the expected value on the `input:` line. It is cut at `==`, `->` or `→` as a safety net, but don't rely on that.
- Hidden tests should cover what visible tests don't: empty input, negatives, mutation of the argument, big inputs and type edge cases.

### The other checks

| `kind` | Fields | Notes |
| --- | --- | --- |
| `unit_tests` | `file` | One row per test function. `visible` defaults to false. `weight` applies per row. |
| `stdout_equals` | `expected` | Runs the entry as `__main__` and compares the output. Trailing whitespace and blank lines are ignored. Visible by default, named "Program output". |
| `stdout_regex` | `pattern` (Go RE2, unanchored), `message` | Visible by default. |
| `exit_code` | `expected` (0-125) | Runs the entry as a program; passes when it exits with that code. Visible by default, named "Exit code". |
| `forbidden_calls` | `names`, `message` | Static scan for `f(...)` and `x.f(...)`. A teaching aid, not a sandbox. |

Every check also accepts `name` (except `unit_tests`), `visible` and `weight` (0 or more; the score is the weighted share of passing rows).

### `wrong/`

Add 2 or 3 realistic mistakes, one file each, written as a learner would write them: an off-by-one, mutating the input, using the banned function. `cq test` fails if any of them passes, which proves your hidden tests actually catch those mistakes.

### Starter code

The starter should run without errors and fail the tests: a function stub that returns an empty value, not `pass` or `...` that crashes the tests. A "fix the bug" lesson is just a `code` lesson whose starter has the bug in it.

### Go and JavaScript lessons

A `code` lesson can also be written in Go or JavaScript. These run on the **wasm runner**: an interpreter compiled to WebAssembly, run on the grading server in about 20-50 ms. Run and Submit both go to the server. See `content/go-basics` and `content/js-basics`.

```
02-ship-count/
  exercise.yaml
  prompt.md
  starter.go          # or starter.js
  solution.go         # server-only
  wrong/
    always_plural.go  # same extension as the entry
```

```yaml
spec:
  language: go               # go | js
  runner: wasm
  run: server                # the default; or off. There is no browser runner for these yet.
  prompt: prompt.md
  entry: main.go             # main.js for JavaScript
  starter: starter.go
  solution: solution.go
  checks:
    - kind: io_cases
      name: Count            # rows "Count 1", "Count 2", ...
      cases:
        - { stdin: "3\n", stdout: "3 ships docked\n" }
        - { stdin: "1\n", stdout: "1 ship docked\n" }
        - { stdin: "0\n", stdout: "0 ships docked\n", visible: false, name: Empty harbor }
    - kind: exit_code         # optional: the program run with empty stdin exits with this code
      expected: 0
```

These lessons are **programs, not functions**. The learner reads stdin and prints to stdout, and you grade the output. There are no test files: `unit_tests` and `forbidden_calls` are Python-only.

| `kind` | Use it for |
| --- | --- |
| `io_cases` | The main check. Each case is a stdin and the exact stdout expected (compared like `stdout_equals`) and becomes one row. 1-50 cases, at most one `io_cases` check per lesson. Hidden cases (`visible: false`) are your hidden tests: name them in words, and their input never leaves the server. |
| `stdout_equals` / `stdout_regex` | Programs that read no input. |
| `exit_code` | Error handling: "exit with 2 on a bad signal". The program runs with **empty** stdin, so pick an expected code that matches that. |

What the learner's code can use:

- **Go** runs on an interpreter (yaegi). `package main` with `func main()`, and these standard packages only: `bufio bytes cmp container/heap container/list errors fmt iter maps math math/bits math/rand os regexp slices sort strconv strings sync time unicode unicode/utf8`. Anything else, including `net` and `encoding/json`, fails with "package … is not available in lessons". The interpreter doesn't reject unused variables or imports the way `go build` does, so don't write a lesson about that error. Read input with `fmt.Scan` or `bufio.Scanner`.
- **JavaScript** is ES2015+ (goja). Node-like globals: `console.log/error`, `process.exit(code)`, `process.stdout.write`, and `require("fs").readFileSync(0, "utf8")` for all of stdin. No other modules, no npm, no `fetch`, no timers. `console.log` formats arrays and objects like Node does (`[ 1, 'a' ]`, `{ a: 1 }`) for common values; print strings yourself when the exact format matters.
- No files except the learner's own, no network, no environment variables. Limits work as for Python (`limits.timeout_ms`, `limits.memory_mb`). All runs of one submission share 2 × `timeout_ms`.

In prompts, show the input and output side by side, and say what to print for empty input. Hidden cases should cover what visible ones don't: empty input, zero, one, ties, very large values.

## 8. Lesson type: `sql`

```
05-sql-query/
  exercise.yaml
  prompt.md
  setup.sql         # public: schema and rows the learner queries
  solution.sql      # server-only
  hidden_data.sql   # server-only: another dataset for a hidden check
  starter.sql       # optional
  wrong/            # .sql files that must fail
```

```yaml
spec:
  prompt: prompt.md
  setup: setup.sql
  solution: solution.sql
  # starter: starter.sql      # optional
  # run: browser              # browser (default, sql.js) | server | off
  # limits: { timeout_ms: 2000 }
  checks:
    - kind: sql_result_equals
      name: Heroes level 5 and up, highest first
      ordered: true
      expected:                 # required on visible checks
        columns: [name, level]
        rows: [[Zara, 9], [Kai, 5]]
    - kind: sql_result_equals   # hidden: compared with the solution's result on another dataset
      setup: hidden_data.sql
      visible: false
      ordered: true
      weight: 2
```

- Only values are compared: the column count must match, but names may differ. Numbers compare with a small tolerance.
- Set `ordered: true` only when the prompt asks for an order (`ORDER BY`).
- A check with its own `setup` must be hidden, because the browser only has the public `setup.sql`.
- Learner queries are one statement. `ATTACH`, `DETACH` and `PRAGMA` are refused, and results are capped at 1000 rows.
- Use a hidden dataset with edge cases the public one lacks (ties, boundary values, NULLs). It's what stops hard-coded answers.

## 9. Lesson type: `project_cli` (on the learner's machine)

These lessons are solved in the learner's own editor and terminal. The `codequest` CLI runs the steps locally, and the server judges the output it reports.

```
01-cli-hello/
  exercise.yaml
  prompt.md
  solution.json     # a recorded passing run (for cq test)
  wrong/            # recorded failing runs, .json only
    crash.json
```

```yaml
spec:
  prompt: prompt.md
  base_url: http://localhost:8000     # optional: default for ${baseURL}
  allowed_os: [linux, darwin]         # optional: linux | darwin | windows
  solution: solution.json
  steps:                              # 1-20 steps
    - description: Run your script
      cli_command:
        command: python3 tour.py      # sh -c (cmd /C on Windows), in the learner's folder; 30 s limit
        tests:
          - exit_code: 0
          - stdout_contains_all: ["Tour complete"]
          - stdout_contains_none: ["Traceback"]
          # - stdout_lines_gt: 3
    - description: Fetch ping.txt from your server
      http_request:
        request: { method: GET, url: "${baseURL}/ping.txt" }   # also headers, body_json, follow_redirects; 10 s limit
        tests:
          - status_code: 200
          - body_contains: "pong"
          - header: { name: Content-Type, contains: text/plain }
          # - json_body: { path: .items[0].id, op: eq, value: 1 }   # op: eq | ne | gt | lt | contains
```

- "stdout" means stdout and stderr combined.
- **Captures** pass values between steps: `captures: [{ name: token, regex: "token=(\\w+)" }]` (first submatch), or `json_path` for HTTP. Later steps use `${token}`. The server recomputes captures from the reported output.
- Every test is visible to the learner. There is no hidden check here, so keep the prompt honest about what's tested.
- Recorded answers look like this:

  ```json
  {"steps": [{"cli": {"exit_code": 0, "stdout": "Tour complete\n"}}]}
  {"steps": [{"http": {"status": 200, "headers": {"Content-Type": "text/plain"}, "body": "pong\n"}}]}
  ```

- These lessons run commands on the learner's computer. Never write a step that deletes files outside the lesson folder, needs `sudo`, or touches anything the learner didn't create for the lesson.

## 10. Boss fights

A boss is a `code` or `project_cli` lesson with `boss: true`, normally the last lesson of a chapter.

- Each check row is one HP. Visible and hidden tests both count, so a boss with 3 visible and 4 hidden tests has 7 HP. Aim for 5-8.
- A `project_cli` boss has one HP per step test (every test is visible). Use it for hands-on courses such as Docker, where the chapter's skill lives on the learner's machine. The learner attacks with `codequest submit`.
- While a boss is undefeated, every later chapter of the course is locked.
- Defeating it gives +100 XP, a boss chest and Boss Slayer progress, and it appears on the learner's profile as a shipped project. Give it a project-like title ("Tip Calculator CLI", not "Exercise 12").
- It should combine the chapter's skills into one small program. Nothing new should be introduced in the boss itself.

## 11. Images

Put images in a flat `assets/` folder inside the lesson and reference them from Markdown or an mcq option:

```markdown
![nums [1, 2, 3] goes through × 2 and becomes [2, 4, 6]](assets/double-flow.svg)
```

`cq validate` fails on:

- a reference to a file that isn't there, or a file in `assets/` that nothing uses (`solution.md` doesn't count)
- external images (`http:`, `//`, `data:`)
- types other than `png`, `jpg`/`jpeg`, `webp`, `gif` and `svg`, or bytes that don't match the extension
- files over 1 MB, or over 5 MB per lesson
- subfolders, symlinks and server-only names (`solution.*`, `tests_hidden.*`, `wrong`)
- SVGs with scripts, event handlers, `foreignObject`, entities, or `href`/`url()` other than `#fragment`

Everything in `assets/` is public, so never draw the answer. Write real alt text: it's what screen readers say and what shows if the image fails. Diagrams should read in both dark and light themes (they sit on a neutral panel), so avoid black-on-transparent strokes.

## 12. Career paths

`content/paths.yaml` groups courses into paths for the Courses page, Home and the profile:

```yaml
paths:
  - id: ai-engineer                 # lowercase slug, unique
    title: AI Engineer
    description: From Python warm-ups to chat models and prompt templates.
    goal: ai                        # optional: zero | backend | interviews | ai (at most one path per goal)
    courses:                        # in order
      - { course: foundations, note: Python and SQL basics }        # note: at most 60 characters
      - { course: langchain-basics, note: Chat models and prompts }
```

- `goal` matches the choice a learner makes during onboarding. Their path is marked "Your path", and Home shows its next unstarted course as "Up next".
- Every course must exist in `content/` and be published.
- Deleting a path from the file hides it; nothing is lost.

## 13. Rules that keep answers secret

These are hard rules. The build and an integration test enforce them; don't work around them.

- `solution.*`, `tests_hidden.*` and `wrong/` never reach the browser. Visible test files do.
- Don't put the answer in anything public: `prompt.md`, examples, hints, visible test names, images, or a hidden test's `input:` line.
- Hints are nudges, not code. Point at the concept or the line ("What should `total` start at?"), never the fix.
- The AI tutor sees the prompt, the learner's code and their results, never the solution. Write prompts clearly enough that a nudge is all anyone needs.

## 14. Changing published content

| Change | What to do |
| --- | --- |
| Fix a typo in a prompt, title or hint | Just edit it. |
| Change tests, checks, the answer or limits (anything that changes grading) | Bump `version`. A submit made against the old version is refused with `409 stale_version`, and the learner reloads the lesson to get the new one. |
| Reorder lessons or chapters | Edit the lists in `chapter.yaml` / `course.yaml`. Progress follows the `id`s. |
| Remove a lesson, chapter, course or path | Remove it from its manifest (and delete the folder). `cq build` **unpublishes** it; learners' history is kept. |
| Rename an `id` | Don't. Create a new lesson with a new id and remove the old one instead. |

## 15. Checklist before you open a PR

- [ ] `make validate` is clean.
- [ ] `make content-test` passes: every solution passes, every `wrong/` answer fails.
- [ ] Browser lessons (`run: browser`): `make e2e` passes the content check.
- [ ] Opened in the app with `make content`: the starter fails cleanly, a wrong answer gets useful feedback, and the solution passes.
- [ ] Visible tests are named after the call; hidden tests are named in words and have an `input:` line.
- [ ] Nothing public gives away the answer.
- [ ] `free: true` is set on the free chapters.
- [ ] Each chapter ends with a boss of 5-8 HP.
- [ ] Images have real alt text and read in both themes.
- [ ] `version` is bumped on any lesson whose grading changed.

## 16. Common errors

| Message (after `path:`) | Fix |
| --- | --- |
| `field colour not found in type …` | Unknown or misspelled field. Check the spelling against this guide. |
| `folder is not listed in chapter.yaml exercises` | Add the folder to the list, or move it into a hidden folder like `.drafts/`. |
| `"03-x" has no exercise.yaml` | A listed folder is missing or misnamed. |
| `exercise id "…" is already used by …` | Ids are global across all courses. Pick another. |
| `boss: only code and project_cli exercises can be bosses` | Make it a `code` or `project_cli` lesson, or remove `boss: true`. |
| `answer: has 2 ids but multiple is false` | Set `multiple: true`, or give one answer. |
| `a check with its own setup must be hidden` | Add `visible: false` to the SQL check that has `setup:`. |
| `run: "…" must be browser, server or off` | Fix `run:` (or delete it to use the default). |
| `courses: "…" is not a course in this folder` (paths.yaml) | Use the course's `id`, not its folder or title. |
| `cq test`: `[solution]: fail` | The solution fails a check. Run it in the app to see which one. |
| `cq test`: `[wrong/x.py]: pass` | Your tests don't catch that mistake. Add a hidden test that does. |
