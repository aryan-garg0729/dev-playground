# Push Rejected

What happens when you and a teammate both commit, and they push first?

```sh
git push
#  ! [rejected]        main -> main (fetch first)
# error: failed to push some refs
# hint: Updates were rejected because the remote contains work that you do not have locally.
```

Git refuses, because your push would throw away their commit. You have to get their work first. The tidy way is to pull with rebase, so your commit is replayed on top of theirs and history stays a straight line:

```sh
git pull --rebase
git push
```

## Assignment

1. **As the teammate**, in `~/git-course/teammate`, add a line and push. Copy and paste:

   ```sh
   cd ~/git-course/teammate
   git pull --ff-only
   echo '- Nice site!' >> guestbook.md
   git commit -am "B6: sign the guestbook"
   git push
   ```

2. **As you**, back in `~/git-course/project`, link the guestbook and commit it as **`C6: link the guestbook`**:

   ```sh
   cd ~/git-course/project
   echo '- [Guestbook](guestbook.md)' >> index.md
   ```

3. Run `git push` and read the rejection.
4. Pull with rebase, then push again.
5. Run, then submit `gt.sharing.rejected`.
