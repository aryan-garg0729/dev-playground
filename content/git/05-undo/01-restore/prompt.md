# Throw Away a Change

Everyone makes mistakes. Git's superpower is that almost every mistake can be undone, if you know which tool to reach for. This chapter is that toolbox.

The simplest case: you changed a file and want it back the way it was in the last commit.

```sh
git restore contact.md
```

That throws away your **uncommitted** changes to that one file. Careful: those changes were never saved in a commit, so they're gone for good.

(`git restore .` does it for every file in the folder. Use it only when you really mean "throw away everything".)

## Assignment

In `~/git-course/project`, set up the situation: you're adding a line to the about page, and something messes up the contact page at the same time. Copy and paste:

```sh
echo 'I also play the guitar.' >> about.md
echo 'oops' > contact.md
git status
```

1. Throw away the change to **`contact.md`** only.
2. **Keep** your `about.md` edit. Don't stage or commit it yet: the next lesson needs it.
3. Check with `git diff`, then run and submit `gt.undoing.restore`.
