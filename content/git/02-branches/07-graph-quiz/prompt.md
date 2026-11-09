# Reading the Graph

`--graph` draws history as lines, and `--parents` prints each commit's parents after its own hash:

```sh
git --no-pager log --oneline --graph --parents -n 5
# * 3e9b7d2 6d14c0a E: tidy the menu
# *   6d14c0a 51ab3f9 2c7e880 D: merge the price list
# |\
# | * 2c7e880 8c2d1e7 B: add the price list
# * | 51ab3f9 8c2d1e7 C: link the menu
# |/
# * 8c2d1e7 4ad50fb A: add the menu
```

Read it bottom to top: `A` is shared, history splits into `B` and `C`, and they join again higher up. Run the same command in your repo and compare.

**Which commit is the merge commit, and how can you tell?**
