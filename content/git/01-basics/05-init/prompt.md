# Your First Repository

A **repository** (or "repo") is a project folder that Git is keeping track of. Turn any folder into one with `git init`:

```sh
git init
# Initialized empty Git repository in /Users/you/git-course/project/.git/
```

That created a hidden folder called **`.git`**. It's Git's whole memory of the project: every snapshot, every branch, every setting. Your own files stay where they are, and Git just watches them.

```sh
ls -a
# .  ..  .git
```

Never edit files inside `.git` by hand. Delete that folder and the project's whole history is gone (your current files stay).

Ask Git what it sees, any time:

```sh
git status
# On branch main
# No commits yet
# nothing to commit
```

## Assignment

1. Create the project folder inside `~/git-course` and make it a repository:

   ```sh
   mkdir -p ~/git-course/project && cd ~/git-course/project
   git init
   ```

2. Run `git status`. It should say `On branch main`. (If it says `master`, you skipped the last lesson: run `git branch -m main` to rename it.)
3. Run, then submit `gt.basics.init`, **from `~/git-course/project`**. Every lesson from now on runs from this folder.
