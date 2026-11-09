# Three Places a Change Can Be

A file change in Git lives in one of three places:

1. **Working tree**: the files on disk you edit.
2. **Staging area** (also called the *index*): the changes you've picked for the next commit. `git add <file>` puts them there.
3. **Repository**: snapshots that are committed for good. `git commit` turns the staging area into a new commit.

`git status` tells you where everything is:

```sh
echo "# ElixirBytes Handbook" > index.md
echo "# On-call" > oncall.md
git add index.md
git status
# On branch main
# No commits yet
# Changes to be committed:
#         new file:   index.md
# Untracked files:
#         oncall.md
```

An **untracked** file is one Git has never been told to record.

**You now run `git commit -m "start"`. What goes into the commit?**
