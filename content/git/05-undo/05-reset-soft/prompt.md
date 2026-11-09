# Redo the Last Commits

`git reset --soft` is great for "I made a few scrappy commits, let me redo them as one". The commits disappear, but all their changes stay staged, ready for one clean commit:

```sh
git reset --soft HEAD~2     # undo the last two commits, keep everything staged
git commit -m "one tidy commit"
```

(It does the same job as squashing with `rebase -i`, but for the last few commits it's quicker.)

## Assignment

In `~/git-course/project`, on `main`, make two scrappy commits. Copy and paste:

```sh
echo '- Add a photo' > todo.md
git add todo.md && git commit -m "todo 1"
echo '- Fix the footer' >> todo.md
git commit -am "todo 2"
```

1. Use **`git reset --soft`** to undo both commits while keeping their changes staged.
2. Commit everything as **`D5: add a todo list`**.
3. Run, then submit `gt.undoing.reset-soft`.
