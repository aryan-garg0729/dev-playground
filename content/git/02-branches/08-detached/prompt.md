# Detached HEAD

You can switch to any **commit**, not just a branch. That's handy for looking at how things used to be:

```sh
git switch --detach 75d6690
# HEAD is now at 75d6690 A1: start the handbook
cat .git/HEAD
# 75d6690c4e...        (a hash, not "ref: refs/heads/...")
```

This is **detached HEAD**: `HEAD` points straight at a commit instead of at a branch. Nothing is wrong, but there's a trap. If you commit now, the new commit has no branch pointing at it. Switch away and nothing remembers where it is, so it's easy to lose.

The fix is to put a branch on it before you leave:

```sh
git switch -c experiment/try-this    # a new branch right here, and you're on it
```

Branch names can contain `/`. Teams use it to group branches: `feature/...`, `fix/...`.

## Assignment

Mo asks what the front page looked like before the menu link, and wants to try a different version of it from there.

1. Find `A1`'s hash with `git --no-pager log --oneline` and switch to it, detached.
2. Run `cat index.md` and `git status`. Git tells you `HEAD detached at ...`.
3. Add a line to the old page and commit it as **`E2:`**:

   ```sh
   echo "Open every day except Sundays." >> index.md
   ```

4. Save it: create a branch **`experiment/old-index`** right there.
5. Switch back to **`main`** and submit.

The checks confirm `experiment/old-index` holds `E2` with `A1` as its parent, you're back on `main`, and `main` doesn't contain `E2`.
