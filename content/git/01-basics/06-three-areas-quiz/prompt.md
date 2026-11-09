# The Three Areas

Taking a snapshot in Git is a two-step move, like taking a group photo: first you choose who's in the picture, then you press the button.

![Three boxes: your working folder, then the staging area (the next snapshot), then the history of commits. git add moves changes from the folder to the staging area, git commit saves the staging area as a commit.](assets/areas.svg)

1. **Working folder**: your files, as you edit them.
2. **Staging area**: the changes you've picked for the next snapshot. You put them there with `git add`.
3. **History**: the snapshots (commits) you've saved, with `git commit`.

`git status` tells you where everything is:

```
Changes to be committed:        <- staged: will go into the next commit
        new file:   index.md
Untracked files:                <- Git has never saved these
        about.md
```

You create `index.md` and `about.md`, run `git add index.md`, then `git commit`.

**What goes into the commit?**
