# LevelUpBurst

Shown inside the lesson-clear dialog when a pass levels the learner up.

- "LEVEL UP" in `game` (22/26), letter-spaced, `accent-text`, popping in (`cq-pop`).
- Below it, the old level badge (`surface-raised`, `text-muted`), an arrow, then the new badge (`level-bg` / `level-ink`), which pops 200 ms after the banner.
- One confetti burst over the dialog (`motion-celebrate`, about 80 pieces in `accent`, `success`, `medal-gold`, `level-bg` and `accent-line`, read from CSS variables). Decorative, `aria-hidden`, and absent under reduced motion.
- Sound: the dialog plays the `levelUp` fanfare instead of the lesson-clear chime, never both.
- The banner is a `status` region. The badges are labelled "Level n to level n+1".
