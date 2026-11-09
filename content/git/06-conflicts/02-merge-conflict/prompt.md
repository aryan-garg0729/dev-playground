# Two People, One Line

Resolving a conflict means **editing the file until it says what you actually want**, then telling Git you're done:

```sh
git merge feat/handover
# CONFLICT (content): Merge conflict in oncall.md
# Automatic merge failed; fix conflicts and then commit the result.

git status
# Unmerged paths:
#         both modified:   oncall.md

# edit oncall.md: keep the right text, delete all three marker lines
git add oncall.md                       # "this file is resolved"
git commit -m "C: merge the handover changes"
```

You don't have to pick one side. Often the right answer takes something from each.

Changed your mind halfway? `git merge --abort` puts everything back to before the merge.

## Assignment

Two handover fixes land at once.

1. From `main`, create **`feat/handover`**. Change the handover line and add a checklist, then commit as **`A6:`**:

   ```sh
   sed 's/^Handover: .*/Handover: Mondays, 09:30, in #ops/' oncall.md > oncall.tmp && mv oncall.tmp oncall.md
   printf '\n## Handover checklist\n\n- Open incidents\n- Pending deploys\n' >> oncall.md
   ```

2. Switch to `main` and change the **same** line differently, then commit as **`B6:`**:

   ```sh
   sed 's/^Handover: .*/Handover: Tuesdays, 10:00 (moved from Monday)/' oncall.md > oncall.tmp && mv oncall.tmp oncall.md
   ```

3. Merge `feat/handover` into `main`. It conflicts. Read `git status` and the markers in `oncall.md` until you can say which side is which.
4. Resolve it so the line is exactly:

   ```
   Handover: Tuesdays, 09:30, in #ops
   ```

   The day comes from `main`, the time and channel from the branch. Keep the checklist.
5. Finish the merge with a message starting **`C6:`**, and submit.

The checks confirm `C6` is a merge of `B6` and `A6` at the tip of `main`, the handover line is the combined one, the checklist is there, and no marker or leftover text remains.
