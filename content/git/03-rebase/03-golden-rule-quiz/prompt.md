# The Golden Rule

Rebase isn't dangerous. Misusing it is.

Imagine `main` is shared (chapter 5 sets this up for real). Dex already has commits `X` and `Y` from it, and he's building his next change on top of `Y`. Then you rebase `main` to tidy it up, and push. On your side, `X` and `Y` became `X'` and `Y'`.

**What happens to Dex?**
