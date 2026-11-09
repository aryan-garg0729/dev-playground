# Catch Up

You've had a look; now take the teammate's commit. Your `main` has no new commits of its own, so this can be a fast-forward. `--ff-only` makes sure it is (and refuses if a real merge would be needed, so nothing surprising happens):

```sh
git pull --ff-only
# Fast-forward
#  guestbook.md | 3 +++
```

## Assignment

In `~/git-course/project`:

1. Pull `origin/main` into your `main` with `--ff-only`.
2. Check `git status -sb` says `## main...origin/main`, with nothing in brackets.
3. Run, then submit `gt.sharing.pull`.
