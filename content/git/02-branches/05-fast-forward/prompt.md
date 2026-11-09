# Fast-Forward

Mo approved the menu. To bring a branch's work into the branch you're on, use `git merge`:

```sh
git switch main          # go to the branch that should receive the work
git merge add-menu       # bring add-menu's commits in
# Updating 4ad50fb..8c2d1e7
# Fast-forward
#  menu.md | 4 ++++
```

`git merge X` always merges X **into the branch you're on**. Get that backwards and you merge `main` into your feature instead.

**Fast-forward** means `main` hadn't moved since `add-menu` was created. `add-menu`'s commits already sit on top of it, so there's nothing to combine: Git just slides the `main` sticky note forward to the same commit.

```
before:  A1 - B1 - C1          <- main
                     \
                      A2        <- add-menu

after:   A1 - B1 - C1 - A2      <- main, add-menu
```

No new commit is created.

## Assignment

1. Switch to `main` and merge `add-menu` into it.
2. Look at the output: it should say `Fast-forward`.
3. Run `git --no-pager log --oneline --decorate -n 2`. Both branch names sit on `A2`.
4. Stay on `main` and submit.

The checks confirm `main` and `add-menu` point at the same commit and that `main`'s history has no merge commits.
