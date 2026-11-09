# Put Work on a Shelf

You're halfway through a new page when you spot something that has to be fixed on `main` **right now**. Your half-written page isn't ready to commit.

**`git stash`** puts your uncommitted changes on a shelf and gives you a clean folder:

```sh
git stash -u           # shelve changes, including new (untracked) files
git stash list
# stash@{0}: WIP on add-events: 3bcc3a1 J7: merge the reviews page
# ...switch branches, fix, commit...
git stash pop          # take the work back off the shelf
```

Without `-u`, stash only takes changes to files Git already tracks. A brand-new file would stay behind and follow you around.

## Assignment

In `~/git-course/project`:

1. From `main`, create **`add-events`** and start a page. Don't commit it:

   ```sh
   echo "Next event: Git meetup" > events.md
   ```

2. Interruption! Shelve your work (it's a new file, so use `-u`), switch to `main`, and make the fix:

   ```sh
   sed 's/I reply within a day./I reply within an hour./' contact.md > c.tmp && mv c.tmp contact.md
   git commit -am "A8: reply faster"
   ```

3. Back on `add-events`, take your work off the shelf, add a line and commit both as **`B8: add the events page`**:

   ```sh
   echo "Bring a laptop." >> events.md
   ```

4. Put `B8` on top of `A8` in a straight line (rebase), then fast-forward `main` to it.
5. Run, then submit `gt.toolbox.stash`.
