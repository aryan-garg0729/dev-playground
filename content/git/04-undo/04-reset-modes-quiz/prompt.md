# Soft, Mixed, Hard

`git reset <commit>` moves **the branch you're on** to another commit. Its mode decides what else it touches:

| Mode | Branch pointer | Staging area | Working tree (your files) |
| --- | --- | --- | --- |
| `--soft` | moved | unchanged | unchanged |
| `--mixed` (default) | moved | reset to match | unchanged |
| `--hard` | moved | reset to match | reset to match |

```sh
git reset --soft HEAD~1    # undo the last commit, keep its changes staged
git reset HEAD~1           # undo the last commit, keep its changes unstaged
git reset --hard HEAD~1    # undo the last commit and throw its changes away
```

Unlike `restore`, which changes files, `reset` moves a branch. That rewrites history, so the golden rule applies: only reset commits nobody else has.

**You run `git reset ??? HEAD~1`. Afterwards the last commit is no longer on your branch and its changes are nowhere in your files. Which mode did you use?**
