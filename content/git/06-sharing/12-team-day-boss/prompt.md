# Team Day

Boss fight! The teammate pushed over the weekend, including a file with a password in it. Share your own work, catch up, and undo their mistake safely.

Recreate the weekend from the teammate's clone. Copy and paste:

```sh
cd ~/git-course/teammate
git pull --ff-only
echo '- Welcome back!' >> guestbook.md
git commit -am "E6: welcome visitors back"
printf '# Secret\n\nadmin password: hunter2\n' > secret.md
git add secret.md && git commit -m "G6: add the admin password"
git push
cd ~/git-course/project
```

## Assignment

In `~/git-course/project`:

1. From `main`, create **`add-blog`**, add a blog page and commit it as **`H6: add the blog`**:

   ```sh
   printf '# Blog\n\nFirst post coming soon.\n' > blog.md
   ```

   Push it to `origin` **with upstream tracking**. Then switch back to `main`.

2. On `main`, add a line to the about page and commit it as **`I6: mention the blog`**:

   ```sh
   echo 'I write a blog now.' >> about.md
   ```

3. Get `I6` on top of everything the teammate pushed, **without a merge commit**.
4. Undo the teammate's **`G6`** the safe way (keep the default `Revert "..."` message), on top of `I6`.
5. Push `main`, then submit `gt.sharing.team-day-boss`.
