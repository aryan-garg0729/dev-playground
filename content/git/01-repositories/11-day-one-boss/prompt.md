# Handbook Day One

Priya's first ticket for you:

> **HB-1: Make the handbook safe for other people**
>
> Mo keeps leaving `.draft` files lying around and I *know* one of them is going to end up in history. Add contribution rules, make drafts invisible to Git, and send me the tree hash of the commit so I can check it landed. Same rules as always: work email, labelled commits.

## Assignment

Everything you need is from this chapter.

1. Create **`CONTRIBUTING.md`** with a `# Contributing` heading and at least one rule of your own (the label rule is a good one).
2. Mo's draft is already here. Create it: `echo "Wednesday special: fizzing frogspawn" > specials.draft`.
3. Make Git ignore **every** file ending in `.draft`, not just this one.
4. Commit `CONTRIBUTING.md` and your `.gitignore` change as **`C1:`**, directly on top of `B1`.
5. Find `C1`'s tree hash (`git cat-file -p` on the commit shows it) and add it to `answers.txt`:

   ```
   c1_tree=<the full 40-character tree hash>
   ```

6. Before you submit: `git status` should say the working tree is clean, and `specials.draft` must still be on disk.

Submit from `~/elixirbytes/handbook` to attack. Each check is a hit on the boss:

- `C1`'s parent is `B1`, and its author email ends in `@elixirbytes.shop`
- `C1` contains `CONTRIBUTING.md` and no `.draft` file
- `specials.draft` is ignored
- `answers.txt` has `C1`'s real tree hash
- nothing is left uncommitted

Defeat it to unlock **Branches**.
