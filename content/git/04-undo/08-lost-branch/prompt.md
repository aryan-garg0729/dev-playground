# Rescue a Deleted Branch

`git branch -D` deletes a branch even if its commits aren't merged anywhere. Those commits now have nothing pointing at them, so they vanish from `git log --all`.

The reflog still has them. Every commit you made was a move of `HEAD`, so its hash is in there:

```sh
git reflog
# ...
# 2f3a4b5 HEAD@{3}: commit: E4: try a table for the menu
# ...
```

A branch is just a pointer, so put a new one on that hash:

```sh
git branch rescued 2f3a4b5
```

## Assignment

Make the mistake first:

```sh
git switch -c spike/table-format
printf '| Potion | Coins |\n| --- | --- |\n| Mana tonic | 12 |\n' > menu-table.md
git add menu-table.md && git commit -m "E4: try a table for the menu"
git switch main
git branch -D spike/table-format
```

1. Run `git --no-pager log --all --oneline`. No `E4`.
2. Find `E4`'s hash in `git reflog`.
3. Create a branch **`recovered/spike`** pointing at exactly that commit. Don't merge it, and stay on `main`.
4. Submit.

The checks confirm `recovered/spike` points at `E4`, `spike/table-format` is still gone, and `E4` isn't on `main`.
