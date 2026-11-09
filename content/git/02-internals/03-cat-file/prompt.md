# Look Inside a Commit

Time for some plumbing. Everything Git saves is stored as an **object** in `.git/objects`, named by its hash. `git cat-file` shows any object:

```sh
git cat-file -t HEAD    # -t: what type of object is it?
# commit
git cat-file -p HEAD    # -p: print it nicely
# tree 4b825dc642cb6eb9a060e54bf8d69288fbee4904
# parent 7c1e9d0a2b3c4d5e6f708192a3b4c5d6e7f80912
# author Your Name <you@example.com> 1731000000 +0100
# committer Your Name <you@example.com> 1731000000 +0100
#
# A2: ignore answers and logs
```

(`HEAD` means "the commit you're on right now". More on it in the branches chapter.)

A commit is surprisingly small. It's just a few lines:

- **tree**: the id of the snapshot of the project's files
- **parent**: the id of the commit before it
- **author** and **committer**: who and when
- the message

## Assignment

In `~/git-course/project`:

1. Run `git cat-file -p HEAD` and read it. (HEAD should be your A2 commit.)
2. Add two lines to `answers.txt`, with the **full** 40-character ids from that output:

   ```sh
   echo "a2_tree=TREE_ID" >> answers.txt
   echo "a2_parent=PARENT_ID" >> answers.txt
   ```

3. Run, then submit `gt.internals.cat-file`.
