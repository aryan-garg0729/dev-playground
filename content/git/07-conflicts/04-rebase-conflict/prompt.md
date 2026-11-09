# Conflicts During a Rebase

A merge conflicts **once**: Git combines the two ends and you fix the result. A rebase replays your commits **one at a time**, so each commit that touches the contested line can conflict separately.

The loop:

```sh
git rebase main
# CONFLICT (content): Merge conflict in skills.md
# error: could not apply 1a2b3c4... D7: ...

# fix the file, then:
git add skills.md
git rebase --continue          # not git commit: this replays the next commit
# ...maybe another conflict: fix, add, continue again...
```

`git status` tells you which commit is being replayed. `git rebase --abort` puts everything back.

(If `--continue` opens an editor for the commit message, just save and close it.)

## Assignment

In `~/git-course/project`:

1. From `main`, create **`more-skills`** and make two commits that edit the same line. Copy and paste:

   ```sh
   sed 's/^- Writing$/- Writing (blog posts)/' skills.md > s.tmp && mv s.tmp skills.md
   git commit -am "D7: say what I write"
   sed 's/^- Writing (blog posts)$/- Writing (blog posts and docs)/' skills.md > s.tmp && mv s.tmp skills.md
   git commit -am "E7: mention docs too"
   ```

2. Switch to `main` and change that line too. Copy and paste:

   ```sh
   sed 's/^- Writing$/- Technical writing/' skills.md > s.tmp && mv s.tmp skills.md
   git commit -am "F7: be more specific about writing"
   ```

3. Switch to `more-skills` and rebase it onto `main`. It conflicts **twice**. Keep both ideas each time:
   - replaying D7: `- Technical writing (blog posts)`
   - replaying E7: `- Technical writing (blog posts and docs)`
4. Switch to `main`, fast-forward it to `more-skills`, then run and submit `gt.resolving.rebase-conflict`.
