# When Both Sides Moved

A fast-forward only works when the branch you're merging into hasn't moved. Usually it has: someone else committed to `main` while you were working.

```
         B2            <- add-prices
        /
... - A2 - C2          <- main
```

Now neither branch contains the other. Git has to combine them, so it makes a **merge commit**: a commit with **two parents**, one for each side.

```
         B2 ------
        /         \
... - A2 - C2 - D2     <- main
```

Git works out the combined snapshot itself by comparing both sides with their **merge base**, the last commit they share (`A2` here). If the two sides changed different things, it merges silently. If they changed the same lines, you get a conflict; that's chapter 6.

```sh
git merge add-prices -m "D2: merge the price list"
# Merge made by the 'ort' strategy.
```

Without `-m`, Git opens your editor with a default message.

## Assignment

1. From `main`, create and switch to **`add-prices`**. Add a price list and commit it as **`B2:`**:

   ```sh
   cat > prices.md <<'MD'
   # Wholesale prices

   - Mana tonic: 7 coins
   - Frost draught: 5 coins
   MD
   ```

2. Switch back to `main`. Link the menu from the front page and commit as **`C2:`**:

   ```sh
   echo "- [Menu](menu.md)" >> index.md
   ```

3. Merge `add-prices` into `main` with a message starting **`D2:`**.
4. Run `git --no-pager log --oneline --graph` and find the two lines joining at `D2`.
5. Stay on `main` and submit.

The checks confirm `D2`'s two parents are exactly `C2` and `B2`, and that `main` has `prices.md`, `menu.md` and the new link.
