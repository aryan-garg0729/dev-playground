# Rebase a Branch

Here's the everyday use of rebase: you start a branch, `main` moves on, and you want your branch to sit on top of the latest `main`.

```sh
git switch my-branch
git rebase main           # replay my-branch on top of main
git switch main
git merge my-branch       # now a fast-forward: one straight line
```

## Assignment

In `~/git-course/project`:

1. From `main`, create **`add-skills`**. Add a skills page and commit it as **`B4: add the skills page`**:

   ```sh
   printf '# Skills\n\n- Git\n- Writing\n' > skills.md
   ```

   (Yes, `B4` before `A4`. That's the point: your branch starts from the older `main`.)

2. Switch to `main`. Link the gallery and commit it as **`A4: link the gallery`**:

   ```sh
   echo '- [Gallery](gallery.md)' >> index.md
   ```

3. Switch to `add-skills` and rebase it onto `main`. Then switch to `main` and merge `add-skills` (it's a fast-forward now).
4. Run `git --no-pager log --oneline --graph -n 4`: one straight line, `B4` on top of `A4`.
5. Stay on `main`. Run, then submit `gt.rebasing.onto`.
