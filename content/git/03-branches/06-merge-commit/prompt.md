# A Real Merge

Fast-forward only works when `main` hasn't moved. Usually it has: you work on a branch while other commits land on `main`. Now both sides have new commits:

```
          B3          <- add-gallery
         /
...--A3--C3           <- main
```

Merging now makes a **merge commit**: a new commit with **two parents**, one from each side, that joins the histories:

```
          B3
         /  \
...--A3--C3--D3       <- main
```

Git combines the changes automatically when they touch different files or different parts of a file. (When they touch the same lines, you get a conflict. That has its own chapter.)

To name the merge commit, pass a message: `git merge add-gallery -m "D3: merge the gallery"`.

## Assignment

In `~/git-course/project`:

1. From `main`, create and switch to **`add-gallery`**. Create a gallery page and commit it as **`B3: add the gallery`**:

   ```sh
   printf '# Gallery\n\nPhotos coming soon.\n' > gallery.md
   ```

2. Switch back to `main`. Link the projects page from the home page and commit it as **`C3: link the projects page`**:

   ```sh
   echo '- [Projects](projects.md)' >> index.md
   ```

3. Merge `add-gallery` into `main` with the message **`D3: merge the gallery`**.
4. See the shape: `git --no-pager log --oneline --graph`.
5. Stay on `main`. Run, then submit `gt.branching.merge-commit`.
