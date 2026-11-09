# A Teammate Joins

A new teammate gets the project with **`git clone`**: it copies the whole remote, history and all, and sets up `origin` for them automatically.

```sh
git clone ~/git-course/project-origin.git ~/git-course/teammate
```

To see teamwork without a second computer, you'll **play the teammate** in a second folder for a few lessons. Give that folder its own identity with local config, so its commits are clearly theirs.

## Assignment

1. Clone `~/git-course/project-origin.git` into **`~/git-course/teammate`**.
2. Inside `teammate`, set local config: name `Teammate`, email **`teammate@example.com`**:

   ```sh
   cd ~/git-course/teammate
   git config --local user.name "Teammate"
   git config --local user.email "teammate@example.com"
   ```

3. As the teammate, add a guestbook page, commit it and push. Copy and paste:

   ```sh
   printf '# Guestbook\n\n- Hello from the teammate!\n' > guestbook.md
   git add guestbook.md && git commit -m "A6: add the guestbook"
   git push
   ```

   (The clone's `main` already tracks `origin/main`, so plain `git push` works.)

4. Go back to `~/git-course/project`. **Don't pull yet.** Run, then submit `gt.sharing.clone`.
