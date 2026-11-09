# Unstage a File

Staged something by mistake? Take it back out of the staging area with `--staged`:

```sh
git restore --staged passwords.txt
```

This only **unstages** it. The file and its contents stay in your folder, untouched, so nothing is lost.

`git status` even reminds you of this: under "Changes to be committed" it says `(use "git restore --staged <file>..." to unstage)`.

## Assignment

In `~/git-course/project`, a file that should never be committed sneaks in. Copy and paste:

```sh
echo "my password: hunter2" > passwords.txt
git add -A
git status
```

(`git add -A` stages **everything**, which is exactly how secrets get committed.)

1. Unstage **`passwords.txt`**. Leave `about.md` staged.
2. Commit as **`A5: mention the guitar`**.
3. Delete the file: `rm passwords.txt`.
4. Run, then submit `gt.undoing.unstage`.
