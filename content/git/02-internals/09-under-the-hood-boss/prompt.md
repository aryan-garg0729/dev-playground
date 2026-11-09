# Under the Hood

Boss fight! Keep junk out of the history, make a clean commit, then dig into what you made.

## Assignment

In `~/git-course/project`:

1. Create the FAQ page, plus some drafts that must never be committed. Copy and paste:

   ```sh
   printf '# FAQ\n\nAsk me anything.\n' > faq.md
   echo "half-baked idea" > ideas.draft
   mkdir -p drafts && echo "secret plan" > drafts/plan.md
   ```

2. Change `.gitignore` so that **every file ending in `.draft`** and the **whole `drafts/` folder** are ignored (keep the old rules).
3. Commit `faq.md` and `.gitignore` together as **`C2: add the FAQ page`**.
4. Find C2's **tree** id and add it to `answers.txt` as `c2_tree=TREE_ID` (the full 40 characters).
5. Finish with nothing left uncommitted, then submit `gt.internals.under-the-hood-boss`.
