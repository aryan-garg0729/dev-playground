# Your First Commit

Time to take the snapshot. `git commit` saves everything in the staging area, with a **message** that says what changed:

```sh
git commit -m "A1: add the home page"
# [main (root-commit) 9a04b31] A1: add the home page
#  1 file changed, 3 insertions(+)
```

- `-m` gives the message right on the command line. (Without it, Git opens a text editor for you to type one.)
- `9a04b31` is the start of the commit's **id**, a long unique code. Yours will be different.

Good messages say *what* and *why* in a few words: "fix typo in contact page", not "stuff".

**About the labels:** in this course every commit message starts with a label like `A1:`. That's just how the checks find your commits, since your commit ids are different from everyone else's. Real projects don't need them.

## Assignment

In `~/git-course/project`:

1. Commit the staged home page with the message **`A1: add the home page`**.
2. Run `git status`. It should say `nothing to commit, working tree clean`.
3. Run, then submit `gt.basics.first-commit`.
