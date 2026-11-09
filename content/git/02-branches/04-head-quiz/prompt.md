# What Is HEAD?

You'll see `HEAD` in Git's output constantly: `HEAD -> add-menu`, `HEAD~1`, "detached HEAD". It's another small file:

```sh
cat .git/HEAD
# ref: refs/heads/add-menu

git switch main
cat .git/HEAD
# ref: refs/heads/main
```

Try it, then switch back to `add-menu`.

**What is `HEAD`?**
