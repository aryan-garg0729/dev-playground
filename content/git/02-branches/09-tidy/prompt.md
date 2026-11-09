# Tidy Up Branches

Once a branch is merged, its commits are part of `main`. The branch is just a leftover sticky note, and leftovers pile up fast. Delete them:

```sh
git branch -d add-menu
# Deleted branch add-menu (was 8c2d1e7).
```

Deleting a branch deletes the **pointer**, not the commits. They're still on `main`.

`-d` is the safe delete: it refuses if the branch has commits that aren't merged anywhere, because those *would* become hard to find:

```sh
git branch -d experiment/old-index
# error: the branch 'experiment/old-index' is not fully merged.
```

`-D` deletes anyway. Chapter 4 shows how to rescue things after a `-D` you regret.

## Assignment

1. Delete **`add-menu`** and **`add-prices`** with the safe `-d`. Both are merged, so Git agrees.
2. Try `git branch -d experiment/old-index` and read the refusal. **Keep that branch.**
3. Priya wants topic folders: rename **`recipes`** to **`recipes/potions`**.
4. Stay on `main` and submit.

The checks list your branches (`main`, `recipes/potions` and `experiment/old-index` remain, the other two are gone) and confirm `A2`, `B2` and `D2` are still on `main`.
