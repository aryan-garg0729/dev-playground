# Friday Disaster Recovery

Dex, on the incident channel, 18:40:

> ok so. I was "cleaning up" on Sam's laptop. I reset main back two commits and I deleted docs/suppliers because it looked old. Also Mo's front page redesign is in there and it's terrible. Please fix before Priya sees.

Run this from `~/elixirbytes/handbook` (on a clean `main`) to recreate the evening:

```sh
git switch -c docs/suppliers
printf '# Suppliers\n\n- Newt eyes: Bogwater Ltd\n' > suppliers.md
git add suppliers.md && git commit -m "H4: add the supplier list"
git switch main
printf '\n## Opening hours\n\nMonday to Saturday, 09:00-17:00.\n' >> index.md
git commit -am "F4: add opening hours"
printf '# HANDBOOK!!!\n\nclick the links lol\n' > index.md
git commit -am "G4: redesign the front page"
git reset --hard HEAD~2
git branch -D docs/suppliers
```

## Assignment

Everything you need is from this chapter.

1. Put `main` back where it was before Dex's reset, so `F4` and `G4` are on it again.
2. Bring back **`docs/suppliers`**, pointing at `H4`. Don't merge it.
3. Mo's redesign stays in history (he'll want to see it), but the front page must go back to how it was at **`F4`**. Commit the fix as **`I4:`** on top of `G4`. Any undo from this chapter that adds a commit will do.
4. Stay on `main` with a clean working tree, and submit.

Each check is a hit on the boss:

- `main` has `F4`, `G4` and `I4`
- the front page on `main` is exactly the `F4` version
- `docs/suppliers` points at `H4`, and `H4` isn't on `main`
- the working tree is clean

Defeat it to unlock **Remotes**.
