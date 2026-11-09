# Squash Commits

Real work is messy. You write a bit, commit "wip" (work in progress), write more, commit "wip" again. That's fine on your own branch. But reviewers, and you in six months, want **one commit per idea** with a clear message.

**Interactive rebase** lets you rewrite your branch's commits before anyone sees them:

```sh
git rebase -i main
```

Git opens a text editor with a to-do list of the commits it will replay, **oldest first**:

```
pick 1a2b3c4 wip cooking
pick 5d6e7f8 wip hiking
pick 9a0b1c2 wip chess
```

Change words, then save and close the editor, and Git does what the list says. To fold the last two into the first, change their `pick` to `squash` (or just `s`):

```
pick 1a2b3c4 wip cooking
squash 5d6e7f8 wip hiking
squash 9a0b1c2 wip chess
```

Git then opens the editor once more for the combined commit's message. Delete the old lines and write the new one.

Not sure which editor Git opens? Pick one for this repo: `git config --local core.editor nano` (or `"code --wait"` for VS Code). In nano, save with Ctrl+O, Enter, then quit with Ctrl+X.

## Assignment

In `~/git-course/project`:

1. From `main`, create **`add-hobbies`** and make **three** messy commits. Copy and paste:

   ```sh
   printf '# Hobbies\n\n- Cooking\n' > hobbies.md
   git add hobbies.md && git commit -m "wip cooking"
   echo '- Hiking' >> hobbies.md && git commit -am "wip hiking"
   echo '- Chess' >> hobbies.md && git commit -am "wip chess"
   ```

   (`commit -a` stages every change to files Git already tracks, so you can skip `git add`.)

2. Squash all three into **one** commit with the message **`C4: add the hobbies page`**, using `git rebase -i main`.
3. Switch to `main` and merge `add-hobbies` (a fast-forward).
4. Run `git --no-pager log --oneline -n 3`: `C4` directly on `B4`, with no `wip` anywhere.
5. Stay on `main`. Run, then submit `gt.rebasing.squash`.

If something goes wrong in the middle, `git rebase --abort` puts everything back the way it was before the rebase.
