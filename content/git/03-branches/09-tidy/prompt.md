# Tidy Up Branches

Once a branch is merged, its name has done its job. Delete it:

```sh
git branch -d add-projects
# Deleted branch add-projects (was a3c9f20).
```

Deleting a branch only deletes the **name**. The commits are safe, because they're part of `main` now.

`-d` is the safe delete: Git refuses if the branch has commits that aren't merged anywhere, so you can't lose work by accident:

```sh
git branch -d experiment/old-home
# error: the branch 'experiment/old-home' is not fully merged.
```

(`-D` deletes anyway. Use it only when you're sure.)

Rename a branch with `git branch -m OLD NEW`. Names can contain `/` to group them, like `blog/posts` or `fix/typo`.

## Assignment

In `~/git-course/project`:

1. Delete **`add-projects`** and **`add-gallery`** with `-d`.
2. Try `git branch -d experiment/old-home` and read the refusal. **Keep that branch.**
3. Rename **`blog`** to **`blog/posts`**.
4. Stay on `main`. Run, then submit `gt.branching.tidy`.
