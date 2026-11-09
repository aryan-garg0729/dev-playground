# Where Is HEAD?

You've seen `HEAD` a few times. It's Git's **"you are here"** marker:

```sh
cat .git/HEAD
# ref: refs/heads/add-projects
```

So HEAD points at a **branch**, and the branch points at a **commit**:

```
HEAD -> add-projects -> A3
```

When you commit, the branch moves forward to the new commit, and HEAD comes along because it points at the branch. When you `git switch main`, HEAD changes to `ref: refs/heads/main`.

**What is HEAD?**
