# Merge, the Easy Way

The projects page is done. Bring it into `main` with `git merge`: switch to the branch that should **receive** the work, then name the branch to take it from:

```sh
git switch main
git merge add-projects
# Updating 5118e65..a3c9f20
# Fast-forward
#  projects.md | 3 +++
```

**Fast-forward** is the simplest kind of merge. `main` hasn't changed since `add-projects` split off, so there's nothing to combine: Git just slides the `main` pointer forward to the same commit.

```
before:  A1..C2 (main) -- A3 (add-projects)
after:   A1..C2 -- A3 (main, add-projects)
```

## Assignment

In `~/git-course/project`:

1. Switch to `main` and merge `add-projects` into it.
2. Look for `Fast-forward` in the output.
3. Run `git --no-pager log --oneline -n 2`: both branch names sit on `A3`.
4. Stay on `main`. Run, then submit `gt.branching.fast-forward`.
