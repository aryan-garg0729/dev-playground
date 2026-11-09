# Git course: story bible (ElixirBytes)

## World
Same company as the Docker course: **ElixirBytes**, a small online potion shop whose backend
is **brewd**. Its ops handbook (on-call, deploy runbook, menu, contacts) lives in a shared
Google Doc that keeps getting overwritten. You move it into a Git repo and become its
maintainer. Everything is local: the "server" in chapter 5 is a bare repo in a folder.

## Characters (same as Docker)
- **Priya**: platform lead, writes the tickets. Precise, dry.
- **Dex**: on-call SRE. Grumpy, right. Has his own clone from chapter 5.
- **Mo**: product manager, "can we just..." questions, commits things he shouldn't.

## Workspace
- Repo: `~/elixirbytes/handbook`. Learners run `codequest` from there, all course.
- Remote (ch5+): `~/elixirbytes/handbook-origin.git` (bare). Dex's clone: `~/elixirbytes/dex-handbook`.
- Bisect test script (ch7): `~/elixirbytes/bisect-test.sh` (outside the repo on purpose).
- Hotfix worktree (ch7 boss): `~/elixirbytes/handbook-hotfix` (removed at the end).
- `answers.txt` in the repo root is ignored (`.gitignore`) and holds `key=value` lines.

## Check style
- Commits are labelled `<letter><chapter>:` (A1:, B1: ... A2: ...). Checks find them with
  `git log -1 --format=%H ':/^A1:'` and pass hashes between steps with captures. Never
  check a hash, date or path literally.
- `:/regex` eats the rest of the revision, so never add `^`, `~` or `:path` after it:
  capture the hash first, then use `${hash}:file` or `${hash}^`.
- Shell `${VAR}` is reserved for captures; use `$VAR` or `$(...)` in commands.
- Parents: `git log -1 --format='parents=%P' <rev>` contains `${a}` and `${b}`.
- Counts: `echo "x=$(...)"`, test `x=N\n`.
- Clean tree: `git status --short; test -z "$(git status --porcelain)"` with `exit_code: 0`.

## Files (exact lines the conflict lessons rely on)
- `index.md`: `# ElixirBytes Handbook`, a line, then links appended over time.
- `oncall.md` (A1): `Primary pager: +1-555-0142`, `Escalation window: 5 minutes`,
  `Handover: Mondays, 10:00`, then a `## Rotation` section (keeps the Handover line
  away from the end of the file, so ch6's checklist doesn't touch the conflict hunk).
- `runbooks/deploy.md` (B1): step `3. Watch error rate for 10 minutes`.
- `contacts.md` (A5, Dex): `- Brewing floor: ext. 4410` (fixed to 4401 in the ch7 boss).

## State per chapter
- ch1: A1 (index, oncall), B1 (.gitignore, runbooks/deploy.md), C1 (CONTRIBUTING.md, *.draft ignored).
- ch2: branch recipes -> renamed recipes/potions; A2 menu.md (ff); B2 prices.md on add-prices,
  C2 index link on main, D2 merge; E2 on experiment/old-index (from detached A1);
  boss: F2, G2 on launch/specials, H2 on main, I2 merge, branch deleted.
- ch3: A3 on main, B3 rollback rebased onto it; C3 squashed playbook; D3 reworded, E3;
  boss: F3 on main, G3 squashed glossary rebased, docs/glossary deleted.
- ch4: A4 (CONTRIBUTING line, secrets.txt never committed), B4 bad runbook, C4 restored runbook;
  D4 from reset --soft; reflog reset HEAD~3 recovered; E4 on recovered/spike;
  boss: F4, G4 on main, H4 on docs/suppliers (recovered), I4 restores index.md from F4.
- ch5: origin set up, upstream main; A5 (Dex, contacts.md); B5 (Dex), C5 rebased on it;
  D5 + Revert "D5:"; boss: E5 (Dex), G5 (Mo, codes) pushed, H5 on docs/faq (pushed -u),
  I5 rebased, Revert "G5:".
- ch6: A6 (branch), B6 (main), C6 merge; F6 main, D6/E6 rebased with two resolutions;
  Mo's marker commit fixed by G6; boss: H6 main, I6 squash+rebase, J6 --no-ff merge.
- ch7: A7 hotfix on main, B7 from stash; C7, D7 on feat/metrics, D7 cherry-picked to main;
  tags v1.0.0 (main), sprint-3 (A7, lightweight), v0.9.0 (C3); describe of A4;
  audit/history (churn 01..15) bisected; boss: E7 on hotfix/v1.0.1 from v1.0.0,
  tag v1.0.1, F7 --no-ff merge, worktree removed.
