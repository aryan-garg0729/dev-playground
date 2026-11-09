# The Golden Rule

Rebase rewrites history: it throws away the old commits and makes new copies. That's great for **your own** work. It's a problem for commits other people already have.

**The golden rule of rebase: never rebase commits that someone else has already seen.** Usually that means: never rebase `main` or any branch you've shared.

Imagine you and a teammate both have the same `main`. You rebase `main` and share your new version. Your teammate keeps working on their copy.

**What goes wrong?**
