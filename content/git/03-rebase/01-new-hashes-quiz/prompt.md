# Replayed, Not Moved

In chapter 2, a branch that fell behind `main` got a merge commit. That's honest, but on a busy repo the history turns into a tangle of merge bubbles. **Rebase** is the other option.

`git rebase main` (run on your feature branch) takes the commits your branch has that `main` doesn't, and **replays** them one by one on top of `main`'s latest commit:

```
before:  A - B - M          <- main
              \
               F1 - F2      <- feature

after:   A - B - M          <- main
                  \
                   F1' - F2'  <- feature
```

The result is a straight line, as if you'd started your work after `M`. Notice the `'` marks: `F1'` is a new commit with the same changes as `F1`.

**After the rebase, does `F1'` have the same hash as `F1`?**
