# Work on a Branch

To *move* between branches, use `git switch`:

```sh
git switch recipes       # move to an existing branch
git switch -c add-menu   # create a branch here and move to it, in one step
git switch main          # back to main
```

Switching changes two things: which branch your next commit moves, and the files in your working tree, which Git rewrites to match that branch's latest commit.

You'll see `git checkout` in older tutorials. It does this too, plus several unrelated jobs, which is why Git added the clearer `switch` (and `restore`, in chapter 4).

## Assignment

Mo wants the potion menu in the handbook. Do it on a branch, so `main` stays untouched until it's reviewed.

1. From `main`, create and switch to a branch called **`add-menu`**.
2. Create the menu:

   ```sh
   cat > menu.md <<'MD'
   # Potion menu

   - Mana tonic: 12 coins
   - Frost draught: 9 coins
   MD
   ```

3. Commit it as **`A2:`**.
4. Run `git switch main` and `ls`. `menu.md` disappears, because `main` doesn't have it. Then `git switch add-menu` and it's back.
5. Stay on **`add-menu`** and submit.

The checks confirm you're on `add-menu`, that `A2` sits directly on `main`'s commit and contains `menu.md`, and that `main` doesn't have `A2`.
