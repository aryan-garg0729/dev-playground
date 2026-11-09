# You Published Something Wrong

Undoing a commit nobody else has is easy: reset, amend or rebase it away. Once it's pushed, every one of those rewrites history other people already have, and the remote will reject your push anyway.

The safe undo for published work is **`git revert`**. It works out the opposite of a commit's changes and records that as a **new** commit:

```sh
git revert 9f8e7d6
# (your editor opens with the message: Revert "D5: ...". Save and close.)
# [main 1a2b3c4] Revert "D5: add supplier pricing"
```

History only grows, so anyone who pulled the bad commit simply pulls the fix next.

One honest warning: a revert removes the content from the **current** files, not from history. Anyone can still run `git show` on the old commit. If a real secret is pushed, rotate the secret immediately; scrubbing history is a separate job and out of scope here.

## Assignment

1. Publish a page you'll regret. In `~/elixirbytes/handbook` on `main`:

   ```sh
   printf '# Supplier pricing\n\n- Bogwater Ltd: 3 coins per newt eye\n' > pricing.md
   git add pricing.md && git commit -m "D5: add supplier pricing"
   git push
   ```

2. An hour later Priya calls: those prices are under NDA. Dex has probably pulled already.
3. Undo `D5` without rewriting history. Keep Git's default `Revert "..."` message.
4. Push the fix and submit.

The checks confirm the remote's `main` has `D5` with a `Revert "D5: ..."` commit on top, `pricing.md` is gone from the remote's current files, and your `main` matches `origin/main`.
