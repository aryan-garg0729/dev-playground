# SoundToggle

Header icon button, just left of the theme toggle, that mutes and unmutes game sounds.

- Uses the 32px `IconButton`: a speaker icon when on, a speaker with a cross when muted.
- `aria-pressed` is true when muted, and the label says what a click does ("Mute sounds" / "Unmute sounds").
- On by default. The choice is saved per device in localStorage (`cq-sound`).
