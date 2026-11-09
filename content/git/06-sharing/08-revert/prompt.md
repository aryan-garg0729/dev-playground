# Undo a Shared Commit

Once a commit is pushed, other people may already have it. Resetting it away would rewrite shared history (and a normal push would be rejected anyway). Instead, **revert** it: Git makes a **new** commit that does the exact opposite.

```sh
git revert 9a04b31
# [main 4f5e6d7] Revert "D6: publish my phone number"
```

Git opens your editor with a ready-made message, `Revert "..."`. Save and close to accept it. (Add `--no-edit` to skip the editor.)

History only grows: the bad commit stays, and the revert after it cancels it out. Everyone just pulls the fix.

## Assignment

In `~/git-course/project`, on `main`:

1. Publish something you'll regret. Copy and paste:

   ```sh
   printf '# Phone\n\nCall me: 555-0199\n' > phone.md
   git add phone.md && git commit -m "D6: publish my phone number"
   git push
   ```

2. Oops: that was your private number, and the teammate may have pulled already.
3. Undo `D6` with `git revert`, keeping Git's default `Revert "..."` message.
4. Push the fix. Run, then submit `gt.sharing.revert`.
