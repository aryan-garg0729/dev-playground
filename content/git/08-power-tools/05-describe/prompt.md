# Which Version Is This?

Someone reports "it looked wrong in the version from commit `c73f3d1`". Which release was that? **`git describe`** names any commit by the nearest **annotated** tag before it:

```sh
git describe c73f3d1
# v0.9.0-14-gc73f3d1
```

Read it in three parts:

- `v0.9.0`: the nearest annotated tag that comes **before** the commit in history
- `14`: how many commits after that tag
- `gc73f3d1`: `g` (for Git) plus the commit's short id

On a tagged commit, `describe` prints just the tag name.

## Assignment

In `~/git-course/project`:

1. Find `A5`'s id (the "mention the guitar" commit) and run `git describe` on it.
2. Before you look at the output, guess which tag it will name, and why it can't be `v1.0.0`.
3. Add the exact output to `answers.txt` as `describe=OUTPUT`.
4. Run, then submit `gt.toolbox.describe`.
