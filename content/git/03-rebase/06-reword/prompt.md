# Say What You Did

A commit message is a note to the next person who runs `git log` while something is on fire. "stuff" doesn't help them.

To fix the message of the **latest** commit, `git commit --amend` is enough. For an older one, use interactive rebase with `reword`. `HEAD~2` means "two commits before HEAD", so this lists the last two:

```sh
git rebase -i HEAD~2
# reword 4e5f6a7 D3: stuff
# pick   8b9c0d1 E3: link the refund policy
```

Save, and Git opens the editor for the new message of `D3`. `E3` is replayed on top with a new hash, since its parent changed.

## Assignment

Mo committed the refund policy straight onto your `main` with a useless message. Recreate his mess on `main` first:

```sh
printf '# Refunds\n\nPotions that fizz out within 7 days are refunded.\n' > refunds.md
git add refunds.md && git commit -m "D3: stuff"
echo "- [Refunds](refunds.md)" >> index.md
git commit -am "E3: link the refund policy"
```

Nobody else has these commits yet, so rewriting them is safe.

1. Change `D3`'s message to exactly **`D3: add the refund policy`**. Keep `E3` as it is.
2. Run `git --no-pager log --oneline -n 3` to check.
3. Stay on `main` and submit.

The checks confirm the last two commits on `main` are `E3: link the refund policy` on top of `D3: add the refund policy`, that `D3: stuff` is gone from every branch, and that `refunds.md` survived.
