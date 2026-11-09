# Main Moved Without You

The loop you'll run most often:

```sh
git switch my-feature
git rebase main          # replay my-feature on top of main's tip
git switch main
git merge my-feature     # now a fast-forward: no merge commit
```

After the rebase, `my-feature` sits on top of `main`, so the merge just slides `main` forward.

## Assignment

Dex wants a rollback section in the deploy runbook. Halfway through, you spot that the front page doesn't link the runbook and fix that on `main` first. Your branch is now based on an older `main`.

1. From `main`, create **`docs/rollback`**. Append a rollback section and commit it as **`B3:`**:

   ```sh
   printf '\n## Rollback\n\nRedeploy the previous image tag, then unfreeze the queue.\n' >> runbooks/deploy.md
   ```

   (Yes, `B3` before `A3`. That's the point: your branch starts from the older tip.)

2. Switch to `main`. Link the runbook and commit as **`A3:`**:

   ```sh
   echo "- [Deploying brewd](runbooks/deploy.md)" >> index.md
   ```

3. Rebase `docs/rollback` onto `main`, then fast-forward `main` to it.
4. Run `git --no-pager log --oneline --graph -n 4`: one straight line, `B3` on top of `A3`.
5. Stay on `main` and submit.

The checks confirm `B3`'s parent is `A3`, `main` is at `B3`, there's no merge commit since `A3`, and the rollback section is on `main`.
