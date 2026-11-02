# GameStats

The header's reward cluster: streak, level badge and XP bar. It is the only place pixel type appears in the header.

- **Streak**: flame icon (stroke `accent`) + count in `game` type, `accent-text`.
- **Level**: `game-sm` text "Lv 7" in `level-ink` on `level-bg`, `radius-sm`, 2px × 7px padding.
- **XP bar**: 80px × 6px track in `border`, fill in `accent`, `radius-pill`. Give it `role="img"` and an `aria-label` such as "340 of 500 XP".

The consumer provides the streak days, level, and current/next-level XP. Keep 14px between the three items.
