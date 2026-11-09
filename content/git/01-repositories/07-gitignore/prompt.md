# Things Git Should Ignore

Some files belong in the folder but never in history: build output, editor junk, secrets, and your personal notes. If Git keeps listing them as untracked, sooner or later someone runs `git add .` and commits them.

A file called **`.gitignore`** lists patterns Git should pretend not to see:

```sh
cat .gitignore
# answers.txt
# *.log

git status
# (answers.txt and every .log file no longer appear)
```

`.gitignore` is an ordinary file, so you commit it like any other. Then everybody who works on the repo gets the same rules. To ask *why* a path is ignored:

```sh
git check-ignore -v answers.txt
# .gitignore:1:answers.txt        answers.txt
```

Each commit also records its **parent**: the commit that came before it. That's how a chain of snapshots becomes a history.

## Assignment

Some lessons ask you to write things down in `answers.txt`, a scratch file in the repo root that must never be committed.

1. Create `answers.txt` with `echo "# my answers" > answers.txt`.
2. Create a `.gitignore` containing the line `answers.txt`.
3. Add the deploy runbook:

   ```sh
   mkdir -p runbooks
   cat > runbooks/deploy.md <<'MD'
   # Deploying brewd

   1. Freeze the order queue
   2. Roll the canary
   3. Watch error rate for 10 minutes
   4. Promote or roll back
   MD
   ```

4. Run `git status`: you should see `.gitignore` and `runbooks/`, but not `answers.txt`.
5. Commit `.gitignore` and the runbook together, on top of `A1`, with a message starting **`B1:`**.
6. Submit.

The checks confirm `answers.txt` is ignored by your `.gitignore` and not tracked, both new files are tracked, `B1`'s parent is `A1`, and the working tree is clean.
