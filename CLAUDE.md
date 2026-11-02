# CodeQuest

Gamified coding-learning platform. React + TypeScript in `web/`, Go in `api/`, browser runner in `sandbox/`, courses as files in `content/`.

## Read first
- `docs/SPEC.md` is the source of truth. Build phase by phase (section "Build phases").
- UI must match `docs/design/screens/*.png` and use tokens from `docs/design/tokens.json`.
- Component rules: `docs/design/README.md` and `docs/design/components/<Name>/README.md`.

## Rules
- Never start a phase until the previous phase's acceptance criteria pass.
- Before coding, list the files you will create or change. After coding, report each acceptance criterion and how you verified it.
- No hex colors in components; use the CSS variables in `web/src/styles/tokens.css` (generated, never hand-edited).
- Grading is the only authority on pass/fail. Nothing in `game/` or `economy/` may change a `Result`.
- Never send `solution.*`, hidden tests or `wrong/` files to the browser.
- Security rules in the spec are hard requirements. Never relax them to make a test pass.
- Run `make test` before saying a task is done.
- When a decision changes, update `docs/SPEC.md` in the same change.
