# Push Rejected

Sooner or later this happens:

```sh
git push
#  ! [rejected]        main -> main (fetch first)
# error: failed to push some refs to '/home/sam/elixirbytes/handbook-origin.git'
# hint: Updates were rejected because the remote contains work that you do not
# hint: have locally.
```

Someone pushed while you were working. The remote's `main` has a commit yours doesn't, and yours has one the remote doesn't. A push only ever **fast-forwards** the remote's branch, so Git refuses rather than throw their work away.

You've solved this exact shape before: rebase your commit on top of theirs. `pull` can do the fetch and the rebase in one go:

```sh
git pull --rebase
# Successfully rebased and updated refs/heads/main.
git push
```

Your commit isn't on the remote yet, so rebasing it doesn't break the golden rule. (Plain `git pull` on diverged branches asks you to choose between merging and rebasing. Pass `--rebase` and there's nothing to choose.)

## Assignment

1. **As Dex**, in `~/elixirbytes/dex-handbook`: add his night-shift contact and push it.

   ```sh
   cd ~/elixirbytes/dex-handbook
   git pull --ff-only
   echo "- Night shift: Dex (dex@elixirbytes.shop)" >> contacts.md
   git commit -am "B5: add the night-shift contact"
   git push
   ```

2. **As you**, back in `~/elixirbytes/handbook`: link the contacts page and commit as **`C5:`**:

   ```sh
   echo "- [Contacts](contacts.md)" >> index.md
   ```

3. Run `git push` and read the rejection.
4. Rebase your commit onto Dex's, then push again.
5. Submit from `~/elixirbytes/handbook`.

The checks confirm the remote's `main` ends with your `C5` directly on Dex's `B5`, that both changes are there, and that your `main` matches `origin/main`.
