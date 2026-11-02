# Button

Three button styles: primary (Submit), secondary (Run), and dashed accent (Ask the tutor).

- **Primary**: `accent` fill, `accent-ink` text, `control` type, `radius-md`, 34px tall in the editor toolbar. One per screen.
- **Secondary**: `surface-raised` fill, `text` label; in Light add a 1px `border`.
- **Tutor**: transparent, `accent-text` label, 1px dashed `accent-line` border, full width in the problem pane, 40px tall.
- **Icon button**: 32px square, transparent, 1px `border`, icon in `text-secondary`; always has an `aria-label`.

Use a real `<button type="button">`. Disabled buttons use `text-faint` and no fill change; never hide the label.
