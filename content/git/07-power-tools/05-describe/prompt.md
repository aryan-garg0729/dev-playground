# Which Version Is This?

A bug report says "it was wrong in the build from commit `c73f3d1`". Which release was that? **`git describe`** names any commit by the nearest **annotated** tag behind it:

```sh
git describe c73f3d1
# v0.9.0-14-gc73f3d1
```

Read it in three parts:

- `v0.9.0`: the nearest annotated tag that is an ancestor of the commit
- `14`: how many commits after that tag
- `gc73f3d1`: `g` (for Git) plus the commit's short hash

Lightweight tags are ignored unless you add `--tags`. On a tagged commit, `describe` prints just the tag name.

## Assignment

1. Find `A4`'s hash (the "ask before renaming files" commit) and run `git describe` on it.
2. Before you look at the count, guess which tag it will name, and why `v1.0.0` can't be the answer.
3. Add the exact output to `answers.txt`:

   ```
   describe=<the output of git describe for A4>
   ```

4. Submit.

The checks run `git describe` on `A4` themselves and look for exactly that line in `answers.txt`.
