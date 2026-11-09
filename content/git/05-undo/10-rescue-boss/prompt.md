# Rescue Mission

Boss fight! A bad evening: a branch got deleted, `main` was reset, and a terrible home page redesign is in there somewhere. Use this chapter's tools to fix it all.

In `~/git-course/project`, on a clean `main`, recreate the disaster. Copy and paste:

```sh
git switch -c add-resume
printf '# Resume\n\n- Learning Git\n' > resume.md
git add resume.md && git commit -m "H5: add my resume"
git switch main
printf '\n## Say hello\n\nI reply to every email.\n' >> index.md
git commit -am "F5: add a hello section"
printf '# WEBSITE!!!\n\nclick stuff\n' > index.md
git commit -am "G5: redesign the home page"
git reset --hard HEAD~2
git branch -D add-resume
```

## Assignment

1. Put `main` back where it was before the reset, so `F5` and `G5` are on it again.
2. Bring back **`add-resume`**, pointing at `H5`. Don't merge it.
3. The redesign stays in history, but the home page must go back to exactly how it was at **`F5`**. Commit the fix as **`I5: undo the redesign`** on top of `G5`.
4. Stay on `main` with nothing left uncommitted, then submit `gt.undoing.rescue-boss`.
