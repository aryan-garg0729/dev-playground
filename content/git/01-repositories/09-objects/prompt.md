# Commit, Tree, Blob

What *is* a commit? Git stores everything as **objects** in `.git/objects`, each named by its hash. There are three kinds you need to know:

- a **commit**: who, when, the message, its parent, and one tree
- a **tree**: a folder listing; each line maps a name to a blob or another tree
- a **blob**: the raw content of one file, with no filename attached

The objects are compressed, so `cat` shows gibberish. `git cat-file -p <hash>` prints one nicely:

```sh
git cat-file -p 9a04b31
# tree d5cbb3d7b48c107b13c003d8c6c188c99fa2d135
# author Sam Rivera <sam@elixirbytes.shop> 1788087478 +0000
# committer Sam Rivera <sam@elixirbytes.shop> 1788087478 +0000
#
# A1: start the handbook

git cat-file -p d5cbb3d
# 100644 blob ecc253f458bb2e4bf7be631cd32f00d7e990ef17    index.md
# 100644 blob 1b7e0a9c2d4f6e8a0b1c3d5e7f9a2b4c6d8e0f1a    oncall.md

git cat-file -p ecc253f
# # ElixirBytes Handbook
# ...
```

So a commit is a pointer to a snapshot: **commit → tree → blob**. Two commits that contain the same file content share the same blob, so Git doesn't store unchanged files twice.

This takes a minute to click. Slowing down here is normal.

## Assignment

Priya wants proof you can read the handbook straight out of the object database, without opening the file.

1. Find `A1`'s hash with `git --no-pager log --oneline`.
2. Walk it: `git cat-file -p` the commit, then its tree. Find the line for `index.md`.
3. `git cat-file -p` that blob and confirm it's the page you wrote.
4. Add a line to `answers.txt` with the blob's **full 40-character** hash:

   ```
   a1_blob=<the blob hash>
   ```

5. Submit.

The checks work out the real blob hash for `index.md` at `A1` and look for exactly that line in `answers.txt`.
