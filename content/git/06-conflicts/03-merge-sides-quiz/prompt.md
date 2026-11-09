# Ours, Theirs, Abort

Sometimes one side is simply right for a whole file. During a merge you can take a side wholesale:

```sh
git checkout --ours <file>      # keep the version from the branch you're on
git checkout --theirs <file>    # take the version from the branch being merged in
git add <file>                  # either way, mark it resolved
git merge --abort               # or give up on the whole merge
```

You're on `main` running `git merge feat/prices`. `prices.md` conflicts, and Priya says the branch's price list is the correct one: all of it.

**What do you run?**
