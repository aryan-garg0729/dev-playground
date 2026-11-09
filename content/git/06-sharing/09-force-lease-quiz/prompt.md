# Force, Carefully

You pushed your own feature branch for review, then rebased it onto the latest `main`. The rebase gave your commits new ids, so a normal push is rejected: the remote's branch isn't part of your history anymore.

Rewriting a branch that only **you** work on is fine, but it needs a forced push:

```sh
git push --force                  # overwrite the remote branch
git push --force-with-lease       # overwrite it, but only if nobody else changed it
```

Never force-push `main` or any branch others build on: that's the golden rule again.

Your teammate sometimes pushes small fixes to other people's branches.

**Which do you use, and why?**
