# Visit the Past

You can switch to **any commit**, not just a branch. That's how you look at the project as it was weeks ago:

```sh
git switch --detach 9a04b31
# HEAD is now at 9a04b31 A1: add the home page
```

Now HEAD points straight at a commit instead of a branch. Git calls this **detached HEAD**. It sounds scary, but it's fine for looking around.

One catch: if you make commits here, **no branch points at them**. Switch away and they become hard to find. Give them a branch right away:

```sh
git switch -c my-experiment     # a new branch, right here
```

## Assignment

In `~/git-course/project`:

1. Find `A1`'s id with `git --no-pager log --oneline` and switch to it, detached.
2. Run `cat index.md` (the very first version) and `git status` (`HEAD detached at ...`).
3. Try a different home page and commit it as **`E3: try a shorter home page`**:

   ```sh
   echo 'Short and sweet.' >> index.md
   ```

4. Save it: create a branch called **`experiment/old-home`** right there.
5. Switch back to **`main`**. Run, then submit `gt.branching.detached`.
