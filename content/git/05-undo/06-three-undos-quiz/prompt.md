# Restore, Reset or Revert?

You now know three undo tools, plus one more you'll use in the sharing chapter:

| Tool | Undoes | Changes history? |
| --- | --- | --- |
| `git restore` | changes in your **files** (or the staging area) | No |
| `git reset` | **commits**, by moving the branch back | Yes: commits vanish from the branch |
| `git revert` | a **commit**, by adding a new commit that does the opposite | No: history only grows |

```sh
git revert 9a04b31
# [main 4f5e6d7] Revert "Add a broken link"
```

The rule of thumb: rewriting history (reset, rebase) is fine for commits only **you** have. Once others have a commit, undo it with a new commit instead.

A commit with a broken link went out yesterday, and your teammates already have it.

**How do you undo it?**
