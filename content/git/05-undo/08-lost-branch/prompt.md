# Rescue a Deleted Branch

Deleted a branch with `-D` and then realised you needed it? Its commits still exist. The reflog remembers every commit HEAD was on, including the ones on that branch:

```sh
git reflog
# ... HEAD@{1}: commit: E5: try a dark theme
```

Point a new branch at it, and it's back:

```sh
git branch rescued-work 7a8b9c0
```

## Assignment

In `~/git-course/project`, make the mistake first. Copy and paste:

```sh
git switch -c try/dark-theme
printf 'background: black\ncolor: white\n' > theme.txt
git add theme.txt && git commit -m "E5: try a dark theme"
git switch main
git branch -D try/dark-theme
```

1. Run `git --no-pager log --all --oneline`. No `E5` anywhere.
2. Find `E5`'s id in `git reflog`.
3. Create a branch **`recovered/dark-theme`** pointing at exactly that commit. Don't merge it, and stay on `main`.
4. Run, then submit `gt.undoing.lost-branch`.
