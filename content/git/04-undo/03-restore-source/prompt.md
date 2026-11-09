# Bring Back an Old Version

This is the problem that started the whole course: Dex pasted over the deploy steps. In a doc, the old text is gone. In Git, every committed version is still there.

`--source` picks *which commit* `restore` copies from:

```sh
git restore --source HEAD~1 runbooks/deploy.md
```

`HEAD~1` is the commit before the one you're on. Any commit works: a hash, a branch, `HEAD~3`.

The file comes back in your working tree as a normal change. You review it, then commit it. The bad commit stays in history and the fix goes on top, so anyone reading the log can see what happened.

## Assignment

Recreate Dex's mistake on `main`:

```sh
printf '# Deploying brewd\n\n1. Ship it\n' > runbooks/deploy.md
git commit -am "B4: tighten the deploy runbook"
```

1. Run `git show HEAD` to see the damage.
2. Bring back `runbooks/deploy.md` exactly as it was **in the commit before `B4`**.
3. Check with `git diff` that the rollback section and all four steps are back.
4. Commit it as **`C4:`**, on top of `B4`.
5. Submit.

The checks compare `C4`'s runbook with the version from the commit before `B4` (they must be identical) and confirm `B4` is still in history with `C4` on top.
