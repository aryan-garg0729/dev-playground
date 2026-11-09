# Create a Branch

Create a branch with `git branch NAME`. List branches with plain `git branch`:

```sh
git branch blog
git branch
#   blog
# * main
```

The `*` marks the branch you're on. Creating a branch does **not** switch to it: you're still on `main`. And the new branch points at the same commit you're on right now.

## Assignment

In `~/git-course/project`:

1. Create a branch called **`blog`**. Stay on `main`.
2. Run `git branch` and check the `*` is still next to `main`.
3. Compare `cat .git/refs/heads/blog` with `cat .git/refs/heads/main`. Same commit id.
4. Run, then submit `gt.branching.create`.
