# One Conflict per Commit

A merge conflicts **once**: Git combines the two tips and you resolve the result. A rebase replays your commits **one at a time**, so each commit that touches the contested lines can conflict separately, and each keeps its own resolved snapshot.

The loop:

```sh
git rebase main
# CONFLICT (content): Merge conflict in runbooks/deploy.md
# error: could not apply 1a2b3c4... D: ...

# fix the file, then:
git add runbooks/deploy.md
git rebase --continue          # not git commit; this replays the next commit
# ...maybe another conflict: fix, add, continue again...
```

`git status` tells you which commit is being replayed and how far along you are. `git rebase --abort` puts everything back.

Ten commits rebased over a busy file can mean ten resolutions. That's the honest argument for squashing before you rebase, which the boss uses.

## Assignment

Dex wants the canary watched longer. Meanwhile Priya wants latency watched too.

1. From `main`, create **`feat/watch-window`** and make two commits that edit step 3 of `runbooks/deploy.md`:

   ```sh
   sed 's/^3\. .*/3. Watch error rate for 20 minutes/' runbooks/deploy.md > d.tmp && mv d.tmp runbooks/deploy.md
   git commit -am "D6: watch the canary for 20 minutes"
   sed 's/^3\. .*/3. Watch error rate for 20 minutes (page if it exceeds 2%)/' runbooks/deploy.md > d.tmp && mv d.tmp runbooks/deploy.md
   git commit -am "E6: page on a 2% error rate"
   ```

2. Switch to `main` and make Priya's change:

   ```sh
   sed 's/^3\. .*/3. Watch error rate and p99 latency for 10 minutes/' runbooks/deploy.md > d.tmp && mv d.tmp runbooks/deploy.md
   git commit -am "F6: watch p99 latency too"
   ```

3. Rebase `feat/watch-window` onto `main`. It conflicts **twice**. Keep both ideas each time:
   - replaying `D6`: `3. Watch error rate and p99 latency for 20 minutes`
   - replaying `E6`: `3. Watch error rate and p99 latency for 20 minutes (page if it exceeds 2%)`
4. Fast-forward `main` to the branch and submit.

The checks confirm `D6` sits on `F6` with the first resolution, `E6` sits on `D6` at the tip of `main` with the second, and no markers are left.
