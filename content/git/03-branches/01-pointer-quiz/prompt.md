# What Is a Branch?

So far every commit went in a straight line on `main`. **Branches** let you work on something on the side (a new page, an experiment, a bug fix) without touching the main version until it's ready.

```
          A3  B3        <- my-idea
         /
A1--B1--C1--D1          <- main
```

Here's the surprise: a branch isn't a copy of anything. **A branch is just a name pointing at a commit.** Look:

```sh
cat .git/refs/heads/main
# 5118e65b19c1d0a3...
```

That's the whole branch: one file holding a commit id. When you make a new commit on a branch, Git moves that name forward to point at the new commit.

**What is a branch, really?**
