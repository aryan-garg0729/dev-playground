# Right Fix, Wrong Branch

You're on a feature branch and notice a bug in some unrelated file. You fix it there and commit. Now ops needs that fix on `main` today, but the feature is nowhere near done.

**`git cherry-pick`** copies one commit's changes onto the branch you're on, as a **new commit** (new hash, same message):

```sh
git switch main
git cherry-pick 7e8f9a0
# [main 2b3c4d5] D: fix the escalation window
```

The original stays where it was. Find hashes with `git --no-pager log --oneline <branch>`.

## Assignment

1. From `main`, create **`feat/metrics`** and make two commits:

   ```sh
   printf '# Metrics dashboard\n\nTODO: pick the graphs.\n' > dashboard.md
   git add dashboard.md && git commit -m "C7: start the metrics dashboard"
   sed 's/^Escalation window: .*/Escalation window: 15 minutes/' oncall.md > o.tmp && mv o.tmp oncall.md
   git commit -am "D7: fix the escalation window"
   ```

2. Get **only** `D7` onto `main`. Don't merge the branch, and don't bring `C7`.
3. Leave `feat/metrics` alive and unmerged. Submit from `main`.

The checks confirm there are two `D7` commits (original and copy), `main`'s tip is a `D7` without `C7`, `feat/metrics` still has both, and `main` has the 15-minute window.
