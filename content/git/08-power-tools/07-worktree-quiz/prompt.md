# Two Branches at Once

Stash is fine for a five-minute interruption. For a fix that takes an afternoon, while your own branch is full of half-done edits, there's a better tool: **worktrees**.

```sh
git worktree add ../project-hotfix -b hotfix/urgent v1.0.0
# Preparing worktree (new branch 'hotfix/urgent')

git worktree list
# /Users/you/git-course/project          3bcc3a1 [main]
# /Users/you/git-course/project-hotfix   2b3c4d5 [hotfix/urgent]

git worktree remove ../project-hotfix
```

Now you can have two editor windows open, one per branch, and switch between them without stashing anything.

**What is `../project-hotfix`?**
