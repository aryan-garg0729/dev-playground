# Catch Up

You've seen what Dex changed. Bring it into your `main`:

```sh
git pull --ff-only
# Updating 5bef8fc..9d0e1f2
# Fast-forward
#  contacts.md | 4 ++++
```

You haven't committed anything since your last push, so this is a fast-forward. `--ff-only` is a safety catch: if a fast-forward isn't possible, `pull` stops and explains, instead of making a merge commit you didn't ask for. The next lesson is that exact situation.

## Assignment

1. Pull `origin/main` into your `main` with `--ff-only`.
2. Check `git status -sb` says `## main...origin/main` with nothing in brackets.
3. Submit.

The checks confirm your `main` is at Dex's `A5`, it matches `origin/main`, and no merge commit was made.
