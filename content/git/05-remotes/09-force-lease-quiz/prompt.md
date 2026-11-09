# Force, Carefully

You pushed your feature branch for review, then rebased it onto the latest `main`. The rebase gave your commits new hashes, so a normal push is rejected: the remote's branch isn't an ancestor of yours any more.

Rewriting a branch that only *you* work on is acceptable, and it needs a forced push:

```sh
git push --force                  # overwrite the remote branch
git push --force-with-lease       # overwrite it, but only if nobody else changed it
```

Dex sometimes pushes small fixes to other people's branches.

**Which do you use, and why?**
