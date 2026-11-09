# Two Branches at Once

Stash works for a five-minute interruption. For a hotfix that takes an afternoon, while your own branch is a mess of half-done edits, there's a better tool:

```sh
git worktree add ../handbook-hotfix -b hotfix/urgent v1.0.0
# Preparing worktree (new branch 'hotfix/urgent')
# HEAD is now at 2b3c4d5 D7: fix the escalation window

git worktree list
# /home/sam/elixirbytes/handbook          3bcc3a1 [main]
# /home/sam/elixirbytes/handbook-hotfix   2b3c4d5 [hotfix/urgent]

git worktree remove ../handbook-hotfix
```

**What is `../handbook-hotfix`?**
