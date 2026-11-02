# CodeQuest starter

Everything Claude Code needs to build CodeQuest, ready to drop into a new repository.

| Path | What it is |
| --- | --- |
| `CLAUDE.md` | Standing rules Claude Code reads every session |
| `docs/SPEC.md` | The build spec: architecture, contracts, data model, game economy, build phases |
| `docs/design/tokens.json` | Every design token, dark and light |
| `docs/design/README.md` | Design system rules: color roles, type, layout, accessibility |
| `docs/design/components/` | Guidelines and an HTML preview for each core component |
| `docs/design/screens/` | PNG of every designed screen (2x) |
| `web/src/styles/tokens.css` | CSS variables generated from tokens.json |
| `web/tailwind.config.ts` | Tailwind theme mapped to those variables |

## Getting started

1. `git init`, copy these files in, and commit.
2. Open Claude Code in the folder and start with a read-only planning prompt:

   > Read CLAUDE.md, docs/SPEC.md and docs/design/. Don't write code yet. Give me (1) anything unclear or contradictory in the spec, (2) your plan for Phase 1 as a list of files, (3) the questions you need me to answer.

3. Then build one phase per session: plan, implement, verify against the acceptance criteria and the screen PNGs, commit, clear context, next phase.

## Screens

| File | Screen |
| --- | --- |
| `lesson-dark.png`, `lesson-dark-outline-open.png` | Lesson, dark theme (outline closed / open) |
| `lesson-light.png`, `lesson-light-outline-open.png` | Lesson, light theme |
| `home.png` | Home with daily quests |
| `course-map.png` | Course map |
| `pass-celebration.png` | Pass celebration overlay |
| `onboarding.png` | Onboarding: goal and daily time |
| `landing.png` | Landing page |
| `pricing.png` | Pricing (annual selected) |
| `boss-fight.png` | Boss fight |
| `power-ups-shop.png` | Power-ups and inventory |
| `profile-medals.png` | Profile, medals tab |
