# Snapshots, Not Copies

Every commit is a **complete snapshot** of the project, not just the changes. So doesn't the `.git` folder grow huge?

No, because objects are named by their contents. `git ls-tree` lists a commit's tree. Compare the last two commits (`HEAD~1` means "the commit before HEAD"):

```sh
git ls-tree HEAD~1 about.md
# 100644 blob 1f7a7a4b2c...    about.md
git ls-tree HEAD about.md
# 100644 blob 1f7a7a4b2c...    about.md
```

Same blob id: `about.md` didn't change between those commits, so Git stored it **once**, and both snapshots point at the same blob.

(For really big histories, Git also squeezes objects together into "pack files". That's out of scope here.)

You make ten commits. You never touch `about.md` in any of them.

**How many copies of `about.md`'s contents does Git store for those ten commits?**
