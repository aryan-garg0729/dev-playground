# CodeQuest

CodeQuest is a calm, professional coding workspace with a thin game layer. The screen should feel like a serious tool you can work in for an hour; the reward moments (XP, level, streak) are where it gets playful.

The system has two themes from one set of tokens: **Dark (Charcoal)**, the default, and **Light**. Every color is a token with a value per theme, so components never hard-code a hex.

## Principles

- **One accent.** Amber (`accent`) is the only brand color. It marks the primary action (Submit), the current lesson, and XP. If two things on a screen are amber, one of them is wrong.
- **Color means state, not decoration.** `success` = passed / done / Easy. `danger` = failed / wrong value. `accent` = current / reward. Everything else is neutral.
- **Pixel type is for game moments only.** `pixel` (Silkscreen, chosen over Pixelify Sans because its 2/5/8 are distinct) appears on the streak count, level badge and +XP pill. Never on instructions, problem text, buttons or code.
- **The workspace stays dense and quiet.** Panes are flat `surface` blocks with a 1px `border`, separated by an 8px gap (`space-2`). No shadows, no gradients.
- **Progress lives in one place.** The outline rail on the left is the only progress indicator on the lesson screen. The header shows the lesson name and position as text, not a second bar.

## Voice

- Short, direct, second person: "Return the total of every integer in `nums`."
- Hints are questions, not answers: "What should total start at?"
- Results say what happened, with numbers: "2 of 3 checks failing · 84 ms".
- Examples follow the problem-site convention: **Example 1:** with **Input:**, **Output:**, **Explanation:** lines, then **Constraints:**.

## Color

| Role | Tokens | Rule |
| --- | --- | --- |
| Ground and panes | `bg`, `surface`, `surface-raised` | Page is `bg`; every pane is `surface`; hover and current rows use `surface-raised`. |
| Lines | `border`, `border-strong` | Pane edges use `border`; control borders, example rules and not-started rail segments use `border-strong`. |
| Text | `text`, `text-secondary`, `text-muted`, `text-faint` | `text-faint` is decorative only (line numbers). |
| Accent | `accent`, `accent-ink`, `accent-text`, `accent-soft`, `accent-line` | `accent` is a fill; put `accent-ink` on it. When amber must be text, use `accent-text`, which darkens in Light to pass contrast on white. |
| States | `success`, `success-soft`, `danger`, `danger-soft`, `progress-current` | Soft variants sit behind their matching text. |
| Code | `code-keyword`, `code-function`, `code-number` | Editor syntax only. |

## Type

- `sans` (Inter) for all interface text. Problem title `page-title` (22/30, 600); body `body` (14/24).
- `mono` (JetBrains Mono) for code, examples, constraints and check rows (`code`, `code-sm`).
- `pixel` (Silkscreen, chosen over Pixelify Sans because its 2/5/8 are distinct) for `game` and `game-sm` only.
- Fonts are hosted on Google Fonts: Inter 400/500/600/700, JetBrains Mono 400/500, Silkscreen 400/700.

## Layout: the lesson screen

1. **Header** (`header-height` 56px, `surface`): logo tile + wordmark · previous / lesson name and "Loops · lesson 4 of 12" / next · streak, level badge, XP bar, theme toggle, avatar.
2. **Outline rail** (left): closed by default at `rail-collapsed` (48px), showing one vertical segment per lesson (`success` done, `progress-current` current, `border-strong` not started). Clicking the rail or any segment opens it to `rail-open` (236px) with lesson names, ticks on done lessons and a lock on the boss lesson.
3. **Problem pane**: pills (difficulty, topic, +XP), title, statement, examples, constraints, then hints and the tutor button pinned to the bottom.
4. **Editor + checks** (takes the remaining width): file tab with Run and Submit, the editor with the failing line tinted `danger-soft`, then the checks list.

Panes wrap to a single column below about 900px; the rail stays collapsed on small screens.

## Motion and sound

Game moments get motion and sound; the workspace does not. Effects are for rare, earned moments only: lesson clear, level up, boss fight and chests. A routine submit, run or fail never animates or makes a sound, because effects that fire every few seconds feel childish and nag.

- Use the motion tokens: `motion-fast` (120ms), `motion-base` (240ms), `motion-slow` (600ms), `motion-celebrate` (1200ms), `ease-out`, `ease-spring` (a slight overshoot, for pops). Keyframes live in `web/src/styles/effects.css`.
- Loading states are the one routine animation: the `Spinner` ring (`border-strong` track, `accent` arc) with a screen-reader label, never "Loading…" text. It stays still under reduced motion.
- Under `prefers-reduced-motion` nothing moves, but every end state still shows (the opened chest, the final HP, the new level).
- Sound is short (≤ 1 s) and quiet, and on by default. The header speaker toggle mutes it per device. A sound is never the only cue: each one has a visible counterpart.
- One sound per moment: a level-up fanfare replaces the lesson-clear chime, and the boss fanfare replaces both.
- Every overlay can be skipped (Escape, backdrop, button) and returns focus.

## Lesson images

Diagrams and pictures in prompts and quiz options sit on `surface-raised` with a 1px `border` and `radius-md`, so transparent SVGs read in both themes. Draw them in a mid-grey that works on both grounds, with amber only to point at the thing that matters. They span the prompt width up to 420px tall (160px tall in quiz options) and open full size on click. Alt text says what the image shows.

## Iconography

Stroke icons on a 24px grid, 2px stroke (2.5–3 for tiny chevrons and ticks), round caps and joins, colored with the text token they sit beside. No filled icons, no emoji. Decorative icons are `aria-hidden`; icon-only buttons carry an `aria-label`.

## Accessibility

- Text meets 4.5:1 on its ground in both themes, except `text-faint` (decorative).
- Amber is never text on white: use `accent-text`.
- Status is never color alone: passed/failed checks also use a tick/cross or the words "expected / got".
- Every interactive element is a real `button` or `a` with at least a 32px box and a 44px target where space allows (`touch-min`).
