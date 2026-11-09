# See What's New

The teammate pushed. Your repo doesn't know yet: `origin/main` only updates when you talk to the remote.

**`git fetch`** downloads new commits and updates `origin/main`, but **never touches your own branches or files**. It's always safe to run.

```sh
git fetch
git status -sb
# ## main...origin/main [behind 1]
```

Look at what came in before you take it:

```sh
git --no-pager log --oneline main..origin/main   # commits on origin/main that main doesn't have
git diff main origin/main                         # what they change
```

## Assignment

In `~/git-course/project`:

1. Fetch from `origin`.
2. Read the teammate's commit with the two commands above.
3. **Don't merge or pull yet.** Run, then submit `gt.sharing.fetch`.
