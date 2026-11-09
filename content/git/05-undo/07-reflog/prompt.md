# The Reflog Saves You

You ran `git reset --hard` and lost three commits. Panic? No. Git keeps a diary of **everywhere HEAD has been**, called the **reflog**:

```sh
git reflog
# 1a2b3c4 HEAD@{0}: reset: moving to HEAD~3
# 9f8e7d6 HEAD@{1}: commit: D5: add a todo list
# 5c4b3a2 HEAD@{2}: commit: C5: bring back the FAQ
```

Read it from the top: the newest entry is the reset. The line just below it shows where you were **before** the reset. That commit still exists; no branch points at it anymore, that's all. Put the branch back:

```sh
git reset --hard 9f8e7d6
```

## Assignment

This lesson only counts if you make the mistake first. In `~/git-course/project`, on a clean `main`:

1. Run `git reset --hard HEAD~3`. Look at `git --no-pager log --oneline -n 3`: `D5`, `C5` and `B5` are gone.
2. Read `git reflog` and find where `main` was just before the reset.
3. Put `main` back exactly there.
4. Run, then submit `gt.undoing.reflog`.
