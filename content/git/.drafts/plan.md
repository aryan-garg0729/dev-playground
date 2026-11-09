# Git course: authoring notes (no story)

Beginner course in boot.dev style (Learn Git 1 + 2 order): one idea per lesson, then do it.
No story or characters. The learner never writes code: the project is a tiny website made of
Markdown pages, created with copy-paste commands.

## Workspace
- Course folder `~/git-course`. The repo is `~/git-course/project`; run `codequest` from there.
- Remote (sharing chapter): `~/git-course/project-origin.git` (bare). A second clone playing
  a teammate: `~/git-course/teammate`.
- `answers.txt` in the repo root is ignored (chapter 2) and holds `key=value` lines.

## Checks
- Commits are labelled `<letter><chapter>:` (A1:, B1: ...). Find them with
  `git log -1 --format=%H ':/^A1:'`; pass hashes with captures. Never check a hash literally.
- `:/regex` eats the rest of the revision: capture the hash, then use `${hash}:file`.
- `${VAR}` is reserved for captures; use `$VAR` or `$(...)` in commands.
- Clean tree: `git status --short; test -z "$(git status --porcelain)"` with exit_code 0.

## Files
- index.md (A1): `# My Website` / blank / `Welcome to my website.`
- about.md (B1), contact.md (C1, `Email: hello@example.com`, `Office hours: Mondays, 10:00`).
