# LessonRail

The chapter outline on the left of the lesson screen and the only progress indicator on it. Closed by default; opens in place and pushes the problem pane right.

**Closed** (`rail-collapsed`, 48px): an "open outline" icon button at the top, then one vertical segment per lesson (4px × 16px, `radius-pill`) inside a 32 × 22px button. Segment colors: `success` done, `progress-current` current, `border-strong` not started or locked. Each segment button's `aria-label` reads "4. Sum a list, current".

**Open** (`rail-open`, 236px): a header row "Chapter 3 · Loops · 3/12" in `caption`/`text-muted` with a close button, then one row per lesson: the same colored segment, the lesson name in `small`, and a trailing tick (done, `success`) or lock (boss, `text-faint`). The current row uses `surface-raised` (dark) or `accent-soft` (light) with a 2px inset `accent` bar on its left edge and 600 weight.

Clicking the rail, any segment or the open button toggles the panel; the close button and the chevron in the panel header close it. Remember the user's choice per device. The consumer provides the lessons array (name, status) and the current index.
