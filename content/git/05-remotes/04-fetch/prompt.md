# Look Before You Pull

Dex pushed, but your handbook doesn't know yet. Nothing in Git updates by itself. Check:

```sh
git --no-pager log --oneline -n 1 origin/main
# 5bef8fc I4: restore the front page        (old news)
```

`git fetch` downloads new commits from the remote and moves your **remote-tracking** branches. It never touches your own branches or your files:

```sh
git fetch
# From /home/sam/elixirbytes/handbook-origin
#    5bef8fc..9d0e1f2  main       -> origin/main

git status -sb
# ## main...origin/main [behind 1]
```

Now you can see what's coming before it lands:

```sh
git --no-pager log --oneline main..origin/main    # commits on origin/main that main doesn't have
git diff main origin/main                         # what they change
```

## Assignment

1. In `~/elixirbytes/handbook`, fetch from `origin`.
2. Read Dex's commit with `git --no-pager log --oneline main..origin/main` and `git diff main origin/main`.
3. **Don't merge or pull yet.** Submit.

The checks confirm `origin/main` has `A5`, your `main` doesn't yet, and `git status -sb` reports you're exactly one commit behind.
