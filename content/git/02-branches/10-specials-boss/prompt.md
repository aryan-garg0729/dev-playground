# Specials Menu Launch

A ticket from Priya:

> **HB-7: Specials page for Wednesday**
>
> Mo wants a specials page live by Wednesday. Build it on a branch called `launch/specials` in two commits: the page, then the link from the menu. While you're on it, I'll need one more rule in CONTRIBUTING on `main`, so main *will* move. Land the specials with a proper merge commit, then clean up the branch.

## Assignment

Everything you need is from this chapter.

1. From `main`, create **`launch/specials`** and make two commits on it:
   - **`F2:`** create `specials.md`:

     ```sh
     printf '# Specials\n\n- Wednesday: fizzing frogspawn\n' > specials.md
     ```

   - **`G2:`** link it from the menu: `echo "- [Specials](specials.md)" >> menu.md`
2. Switch to `main` and add a rule to the end of `CONTRIBUTING.md`, for example `echo "- One topic per commit." >> CONTRIBUTING.md`. Commit as **`H2:`**.
3. Merge `launch/specials` into `main` with a message starting **`I2:`**.
4. Delete `launch/specials` with the safe delete.
5. Stay on `main` with a clean working tree, and submit.

Each check is a hit on the boss:

- `G2`'s parent is `F2`
- `I2`'s two parents are `H2` and `G2`
- `main` is at `I2` and has `specials.md`
- `launch/specials` is gone and you're on `main`
- the working tree is clean

Defeat it to unlock **Rebase**.
