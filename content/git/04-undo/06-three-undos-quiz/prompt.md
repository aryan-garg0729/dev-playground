# Restore, Reset or Revert?

Three commands with similar names:

| Command | What it changes | Rewrites history? |
| --- | --- | --- |
| `git restore <file>` | file contents, in your working tree or staging area | no |
| `git reset <commit>` | where your branch points (and maybe staging and files) | yes |
| `git revert <commit>` | adds a **new** commit that undoes an old one | no |

`git revert` works out the opposite of a commit's changes and commits that:

```sh
git revert 7a8b9c0
# [main 1d2e3f4] Revert "D: add the beta prices"
```

Both the mistake and the correction stay in history. You'll use it for real in chapter 5.

**A bad commit was pushed to the shared `main` an hour ago, and Dex has already pulled it. What's the right undo?**
