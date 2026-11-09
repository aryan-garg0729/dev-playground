# The Undo Trail

Friday, 17:55. You mean to drop one commit and type a number three times too big:

```sh
git reset --hard HEAD~3
git --no-pager log --oneline -n 1
# 5c076f2 A4: ask before renaming files
```

Three commits gone from `git log`. But they aren't deleted: nothing points at them any more, that's all. Git keeps a private diary of everywhere `HEAD` has been, called the **reflog**:

```sh
git reflog
# 5c076f2 HEAD@{0}: reset: moving to HEAD~3
# 9e1f2a3 HEAD@{1}: commit: D4: add pager etiquette
# 4b5c6d7 HEAD@{2}: reset: moving to HEAD~2
# ...
```

Newest at the top. `HEAD@{1}` means "where HEAD was one move ago". Find the entry from just before the mistake and point your branch back at it:

```sh
git reset --hard HEAD@{1}      # or the hash from that line
```

`git log` is the history of the project. `git reflog` is the history of **you**.

## Assignment

This lesson only counts if you make the mistake first. On a clean `main`:

1. Run `git reset --hard HEAD~3`. Look at `git --no-pager log --oneline -n 3`. `D4`, `C4` and `B4` are gone.
2. Read `git reflog` and find where `main` was before the reset.
3. Put `main` back exactly there.
4. Submit.

The checks look for the `reset: moving to HEAD~3` entry in your reflog, then confirm `D4`, `C4` and `B4` are back on `main` and the working tree is clean.
