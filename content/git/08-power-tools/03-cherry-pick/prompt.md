# Pick One Commit

You're on a branch for an experiment, and you notice a small mistake somewhere unrelated. You fix it there and commit. The fix should go to `main` today, but the experiment is nowhere near ready.

**`git cherry-pick`** copies one commit's changes onto the branch you're on, as a **new commit** (new id, same message):

```sh
git switch main
git cherry-pick 7e8f9a0
# [main 2b3c4d5] D8: say I learn every day
```

The original stays where it was. Find commit ids with `git --no-pager log --oneline BRANCH`.

## Assignment

In `~/git-course/project`:

1. From `main`, create **`try/footer`** and make two commits. Copy and paste:

   ```sh
   printf 'Made with Git\n' > footer.md
   git add footer.md && git commit -m "C8: try a footer"
   sed 's/I am learning Git./I am learning Git every day./' about.md > a.tmp && mv a.tmp about.md
   git commit -am "D8: say I learn every day"
   ```

2. Get **only** `D8` onto `main`. Don't merge the branch, and don't bring `C8`.
3. Leave `try/footer` alone. Run, then submit `gt.toolbox.cherry-pick` from `main`.
