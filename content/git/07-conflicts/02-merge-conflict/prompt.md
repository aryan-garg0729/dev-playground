# Resolve a Merge Conflict

Resolving a conflict means **editing the file until it says what you actually want**, then telling Git you're done:

```sh
git merge update-hours
# CONFLICT (content): Merge conflict in contact.md
# Automatic merge failed; fix conflicts and then commit the result.

git status
#         both modified:   contact.md

# open contact.md in your editor: keep the right text, delete all three marker lines
git add contact.md                       # "this file is resolved"
git commit -m "C7: merge the office hours"
```

You don't have to pick one side. Often the right answer takes a bit from each.

Changed your mind halfway? `git merge --abort` puts everything back the way it was before the merge.

## Assignment

In `~/git-course/project`:

1. From `main`, create **`update-hours`**. Change the office hours line and add a section at the end, then commit as **`A7: update the office hours`**. Copy and paste:

   ```sh
   sed 's/^Office hours: .*/Office hours: Mondays, 09:30, online/' contact.md > c.tmp && mv c.tmp contact.md
   printf '\n## Response time\n\nI reply within a day.\n' >> contact.md
   git commit -am "A7: update the office hours"
   ```

2. Switch to `main` and change the **same** line differently. Copy and paste:

   ```sh
   sed 's/^Office hours: .*/Office hours: Tuesdays, 10:00/' contact.md > c.tmp && mv c.tmp contact.md
   git commit -am "B7: move office hours to Tuesday"
   ```

3. Merge `update-hours` into `main`. It conflicts. Read `git status` and the markers in `contact.md`.
4. Resolve it so the line is exactly:

   ```
   Office hours: Tuesdays, 09:30, online
   ```

   The day comes from `main`, the time and "online" from the branch. Keep the new "Response time" section.
5. Finish the merge with the message **`C7: merge the office hours`**. Run, then submit `gt.resolving.merge-conflict`.
