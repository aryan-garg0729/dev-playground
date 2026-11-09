# Unstage a Secret

`git add -A` (or `git add .`) stages *everything*. Usually that's what you want. Sometimes it scoops up a file that must never be committed.

To take a file **out of the staging area** but leave it on disk:

```sh
git restore --staged secrets.txt
git status
# Changes to be committed:
#         modified:   CONTRIBUTING.md
# Untracked files:
#         secrets.txt
```

`--staged` changes the destination: instead of copying into your working tree, `restore` copies the last commit's version into the staging area. For a new file, that means it isn't staged any more.

Unstaging is the cheap fix. Once a secret is **committed**, it's in history, and history gets copied to every clone. Getting it out properly is a much bigger job (chapter 5 shows why).

## Assignment

Mo left something on your laptop:

```sh
echo "pager password: hunter2" > secrets.txt
git add -A
git status
```

Both your `CONTRIBUTING.md` line and `secrets.txt` are now staged.

1. Unstage **`secrets.txt`**. Leave `CONTRIBUTING.md` staged.
2. Commit as **`A4:`**, for example `A4: ask before renaming files`.
3. Tell Mo to use a password manager, then delete the file: `rm secrets.txt`.
4. Submit.

The checks confirm `A4` only touches `CONTRIBUTING.md`, no commit on any branch has ever contained `secrets.txt`, and the working tree is clean with `secrets.txt` gone.
