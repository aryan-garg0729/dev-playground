# Hotfix Release

Final boss! The blog link in version 1.0.0 is broken. Ship a 1.0.1 fix from a separate worktree, tag it and merge it, all with this chapter's tools.

## Assignment

In `~/git-course/project`:

1. Create a second working folder at **`~/git-course/project-hotfix`** on a new branch **`hotfix/v1.0.1`**, starting from the **`v1.0.0` tag**.
2. In that folder, fix the link and commit it as **`E8: fix the blog link`**. Then tag that commit as annotated **`v1.0.1`**:

   ```sh
   cd ~/git-course/project-hotfix
   sed 's#Blog: example.com/blog#Blog: blog.example.com#' contact.md > c.tmp && mv c.tmp contact.md
   ```

3. Back in `~/git-course/project`, merge `hotfix/v1.0.1` into `main` as a **real merge commit** with the message **`F8: merge the 1.0.1 hotfix`**. (`main` hasn't moved since 1.0, so Git would happily fast-forward. Insist on the merge commit.)
4. Remove the extra worktree the proper way, so Git forgets it too.
5. Submit `gt.toolbox.hotfix-boss` from `~/git-course/project`.

When you've won, run `git --no-pager log --oneline --graph --all` and read your whole journey back: branches, merges, rebases, reverts, conflict fixes and the hotfix. `git reflog` still remembers every mistake along the way. That's Git. 🎉
