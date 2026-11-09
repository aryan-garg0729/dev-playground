# Hotfix v1.0.1

Priya, 16:50 on a Friday (it's always a Friday):

> **HB-40: v1.0.1, today**
>
> 1.0 ships the wrong extension for the brewing floor. It says 4410; it's **4401**. The night shift is calling a broom cupboard. I need a 1.0.1 built from exactly what we released, not from whatever's on main, tagged properly, and merged back so main has it too. Don't make a mess of your own working tree to do it.

## Assignment

Everything you need is from this chapter.

1. Create a second working directory at **`~/elixirbytes/handbook-hotfix`** on a new branch **`hotfix/v1.0.1`**, starting from the **`v1.0.0` tag**.
2. In that directory: change `ext. 4410` to `ext. 4401` in `contacts.md` and commit it as **`E7:`**. Tag that commit as annotated **`v1.0.1`**.
3. Back in `~/elixirbytes/handbook`, merge `hotfix/v1.0.1` into `main` as a **real merge commit** with a message starting **`F7:`**. `main` hasn't moved since 1.0, so Git would happily fast-forward. Insist on the merge commit.
4. Remove the extra worktree the proper way, so Git forgets it too.
5. Submit from `~/elixirbytes/handbook`.

Each check is a hit on the boss:

- only one worktree is registered and `../handbook-hotfix` is gone
- `E7` sits directly on the `v1.0.0` commit
- `v1.0.1` is an annotated tag on `E7`
- `F7` merges `v1.0.0` and `E7` at the tip of `main`
- `main` has `ext. 4401`

Defeat it to finish the course. Then run `git --no-pager log --oneline --graph --all` and read the whole story back: the squashes, the rebases, the revert pairs, both conflict resolutions, the cherry-picked twin and the hotfix merge. `git reflog` still lists every mistake you made along the way. That's the safety net.
