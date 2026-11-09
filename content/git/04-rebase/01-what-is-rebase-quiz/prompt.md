# What Rebase Does

Merging joins two histories with a merge commit. **Rebase** is the other way to bring a branch up to date: it takes your branch's commits and **replays them on top of** the latest `main`, as if you'd started your work just now.

Before:

```
          B          <- my-branch
         /
A1--A2--A3           <- main
```

`git rebase main` while on `my-branch`:

```
A1--A2--A3--B'       <- my-branch
```

The result is one straight line: no merge commit, easy to read. Notice the `'` on `B'`: it's a **copy** of B, with the same changes but on a new parent.

Remember what goes into a commit's id: the snapshot, the **parent**, the author, the date and the message.

**After the rebase, does `B'` have the same commit id as the old `B`?**
