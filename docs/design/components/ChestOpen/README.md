# ChestOpen

Modal for opening a chest: the daily chest (Home, "Open") and the boss chest (lesson clear dialog, "Open boss chest").

- **Sequence.** 1. The closed chest wobbles (`cq-chest-wobble`) while the server decides the loot, for at least `motion-slow` so a fast reply still reads as opening. 2. The lid swings open on its back hinge with `ease-spring`, a soft `accent-soft` glow and four short `accent` rays appear, and `chestOpen` plays once. 3. "+N Bytes" pops in (`game`, `accent-text`) and counts up, then the power-up card rises (`surface-raised`, power-up icon in `accent-text`, name in `section`).
- **Chest art** (`ChestSprite`): the body is `accent-line` with an `accent` lock, the lid is `accent`, and the outlines are `accent-text`. All are token variables, so it works in both themes.
- **Title** in `game-sm`, letter-spaced, `accent-text`: "DAILY CHEST" or "BOSS CHEST".
- **Button**: "Skip" while opening, "Collect" once opened. Escape and the backdrop also close it, and focus returns to the opener.
- **Error**: the message in `danger` replaces the loot; the chest stops wobbling.
- **Reduced motion**: no wobble, no lid swing and no count-up. The chest shows open with the loot.
