# Make a Branch

Commands that manage branches without moving you anywhere:

```sh
git branch              # list branches; * marks the one you're on
git branch <name>       # create a branch at the current commit
git branch -m old new   # rename
git branch -d <name>    # delete
```

Creating a branch doesn't switch to it. You stay exactly where you were:

```sh
git branch potions
git branch
# * main
#   potions
```

Both now point at the same commit. They'll only drift apart when someone commits on one of them.

## Assignment

Priya will move the potion recipes into the handbook next month. Reserve a branch for it.

1. Create a branch called **`recipes`** at your current commit. Stay on `main`.
2. Run `git branch` and check the `*` is still next to `main`.
3. Run `cat .git/refs/heads/recipes` and compare it with `cat .git/refs/heads/main`.
4. Submit.

The checks confirm `recipes` exists, you're still on `main`, and both point at the same commit.
