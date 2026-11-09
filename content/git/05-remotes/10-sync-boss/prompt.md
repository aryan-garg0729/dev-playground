# The Discount Code Leak

Priya, Monday morning:

> **HB-21: Monday fires**
>
> 1. Mo pushed the wholesale discount codes to the handbook. Every reseller can read that repo. Get them out of main **without rewriting history**; Dex and half the floor have already pulled. I'm rotating the codes myself.
> 2. Your FAQ page needs review. Push it on its own branch, `docs/faq`, and make sure it tracks the remote so you can push fixes later.
> 3. Your glossary commit still has to go out. Dex pushed over the weekend too, so expect to be behind. Keep main a straight line.

Run this to recreate the weekend from Dex's clone:

```sh
cd ~/elixirbytes/dex-handbook
git pull --ff-only
printf '\n## Severity\n\nSEV1 wakes everyone. SEV3 waits for morning.\n' >> playbook.md
git commit -am "E5: add a severity scale"
printf '# Discount codes\n\nFRIENDS50, NEWT4FREE\n' > codes.md
git add codes.md
git -c user.name=Mo -c user.email=mo@elixirbytes.shop commit -m "G5: add the wholesale discount codes"
git push
cd ~/elixirbytes/handbook
```

## Assignment

Everything you need is from this chapter. In `~/elixirbytes/handbook`:

1. From `main`, create **`docs/faq`**, add an `faq.md` with a `# FAQ` heading and one question, and commit it as **`H5:`**. Push it with upstream tracking. Switch back to `main`.
2. On `main`, add a term to the glossary and commit it as **`I5:`**:

   ```sh
   echo "- SEV: how bad an incident is, 1 being worst." >> glossary.md
   ```

3. Get `I5` on top of everything Dex and Mo pushed, without a merge commit.
4. Undo Mo's **`G5`** the safe way, keeping the default `Revert "..."` message, on top of `I5`.
5. Push `main` and submit.

Each check is a hit on the boss:

- `docs/faq` is on the remote with `H5`, and yours tracks `origin/docs/faq`
- `I5`'s only parent is `G5`
- the remote's `main` ends with `Revert "G5: ..."`, with `I5` and `G5` below it
- `codes.md` is gone from the remote's current files
- your `main` matches `origin/main`

Defeat it to unlock **Conflicts**.
