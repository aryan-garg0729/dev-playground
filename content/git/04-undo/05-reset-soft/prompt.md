# Two Commits Into One

`--soft` is a neat trick for combining commits without an interactive rebase. Moving the branch back two commits while keeping everything staged leaves you exactly one `git commit` away from a single combined commit:

```sh
git reset --soft HEAD~2
git status
# Changes to be committed:
#         new file:   etiquette.md
git commit -m "one tidy commit"
```

## Assignment

Make two scrappy commits on `main`:

```sh
echo "- Answer within 5 minutes." > etiquette.md
git add etiquette.md && git commit -m "etiquette 1"
echo "- Hand over in writing." >> etiquette.md
git commit -am "etiquette 2"
```

1. Use **`git reset --soft`** to undo both commits while keeping their changes staged.
2. Commit everything as **`D4: add pager etiquette`**.
3. Submit.

The checks confirm `D4` sits directly on `C4` at the tip of `main`, the `etiquette 1` and `etiquette 2` commits are gone from every branch, your reflog shows the `reset` you used, and `etiquette.md` has both lines.
