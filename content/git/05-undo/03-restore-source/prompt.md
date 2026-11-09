# Bring Back an Old Version

`git restore` can bring back a file from **any** commit, not just the last one. Use `--source`:

```sh
git restore --source HEAD~1 faq.md     # faq.md as it was one commit ago
git restore --source 9a04b31 faq.md    # ... or as it was in any commit
```

The file in your folder becomes that old version. Then you commit it like any other change. History keeps both the mistake and the fix, which is honest and safe.

## Assignment

In `~/git-course/project`, on `main`, make a bad commit that wipes out the FAQ. Copy and paste:

```sh
printf '# FAQ\n\n?\n' > faq.md
git commit -am "B5: rewrite the FAQ"
```

1. Look at the damage with `git show HEAD`.
2. Bring back `faq.md` exactly as it was **in the commit before B5**.
3. Check with `git diff` that the questions are back.
4. Commit it as **`C5: bring back the FAQ`**.
5. Run, then submit `gt.undoing.restore-source`.
