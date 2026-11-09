# Ours, Theirs, Abort

Sometimes one side is simply right for a whole file. During a merge you can take a side wholesale:

```sh
git checkout --ours FILE      # keep the version from the branch you're on
git checkout --theirs FILE    # take the version from the branch being merged in
git add FILE                  # either way, mark it resolved
git merge --abort             # or give up on the whole merge
```

You're on `main` running `git merge new-prices`. `prices.md` conflicts, and the branch's price list is the correct one: all of it.

**What do you run?**
