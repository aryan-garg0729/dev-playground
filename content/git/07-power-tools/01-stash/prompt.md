# The Pager Goes Off

You're halfway through a new page when Dex pings: the playbook says to post status updates every 30 minutes, the team agreed on 15, and an incident is running *right now*. The fix has to go out from `main`. Your half-written page isn't ready to commit.

**`git stash`** puts your uncommitted changes on a shelf and gives you a clean working tree:

```sh
git stash -u           # shelve changes, including untracked (new) files
git stash list
# stash@{0}: WIP on docs/status-page: 3bcc3a1 J6: merge ...
# ...switch, fix, commit...
git stash pop          # put the top entry back and remove it from the shelf
```

Without `-u`, `stash` only takes changes to files Git already tracks. A brand-new file would stay behind and follow you to `main`.

## Assignment

1. From `main`, create **`docs/status-page`** and start the page. Don't commit it:

   ```sh
   echo "Status page: status.elixirbytes.shop" > status.md
   ```

2. The pager fires. Shelve your work (it's a new file, so think about `-u`), switch to `main`, and make the hotfix:

   ```sh
   sed 's/every 30 minutes/every 15 minutes/' playbook.md > p.tmp && mv p.tmp playbook.md
   git commit -am "A7: post status updates every 15 minutes"
   ```

3. Back on `docs/status-page`, bring your work back, add a second line and commit both as **`B7:`**:

   ```sh
   echo "Update it within 10 minutes of a SEV1." >> status.md
   ```

4. Put `B7` on top of the hotfix in a straight line, then fast-forward `main` to it.
5. Submit.

The checks confirm the stash is empty, `B7` sits directly on `A7` at the tip of `main`, the hotfix is in the playbook, and `status.md` has both lines.
