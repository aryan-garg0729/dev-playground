# Switch Branches

Move to another branch with `git switch`:

```sh
git switch blog          # go to an existing branch
git switch -c new-idea   # -c: create it and switch to it, in one go
```

When you switch, Git changes the files in your folder to match that branch's commit. Files that only exist on the other branch disappear from your folder, and come back when you switch back. Nothing is lost: it's all in the commits.

(You'll also see `git checkout` in older tutorials. It does this too, plus many other things. `git switch` is newer and does just one job.)

## Assignment

In `~/git-course/project`:

1. From `main`, create and switch to a branch called **`add-projects`**.
2. Create a projects page and commit it as **`A3: add the projects page`**. Copy and paste:

   ```sh
   printf '# Projects\n\n- A website made with Git\n' > projects.md
   ```

3. Run `git switch main` and `ls`: `projects.md` disappears, because `main` doesn't have it. Then `git switch add-projects`, and it's back.
4. Stay on **`add-projects`**. Run, then submit `gt.branching.switch`.
