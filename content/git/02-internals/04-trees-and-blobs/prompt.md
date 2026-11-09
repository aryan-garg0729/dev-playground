# Trees and Blobs

A commit points to a **tree**. Let's follow the trail all the way down to a file.

```sh
git cat-file -p 4b825dc6     # a tree id from a commit
# 100644 blob 1f7a7a4b...    about.md
# 100644 blob 557db03d...    index.md
```

There are only a few kinds of objects:

| Object | Is like | Holds |
| --- | --- | --- |
| **commit** | a labelled photo | a tree, the parent, author, message |
| **tree** | a folder | a list of names, each pointing at a blob or another tree |
| **blob** | a file's contents | just the bytes, with no name |

So a commit points to a tree, the tree lists files, and each file points to a blob holding its contents.

```sh
git cat-file -p 557db03d     # a blob: the file's contents
# # My Website
# ...
```

## Assignment

In `~/git-course/project`, find the blob that holds **`index.md` as it was in your A1 commit**:

1. Find A1's id with `git --no-pager log --oneline`.
2. `git cat-file -p A1_ID` to find its tree.
3. `git cat-file -p TREE_ID` to find the blob for `index.md`. Check it with `git cat-file -p BLOB_ID`.
4. Add it to `answers.txt` as `a1_blob=` followed by the **full** 40-character blob id.
5. Run, then submit `gt.internals.trees-and-blobs`.

(Shortcut, once you've done it the long way: `git rev-parse A1_ID:index.md`.)
