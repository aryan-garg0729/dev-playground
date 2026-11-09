# Release Day

Branches move. A **tag** is a name for a commit that stays put: perfect for "this is what we shipped as 1.0".

There are two kinds:

```sh
git tag sprint-3                          # lightweight: just a name pointing at a commit
git tag -a v1.0.0 -m "Handbook 1.0"       # annotated: a real object with tagger, date, message
git tag -a v0.9.0 -m "..." 7f8e9d0        # tag any commit, not just HEAD
```

Releases use **annotated** tags, because they record who cut the release, when, and why. You can see the difference in the object database:

```sh
git cat-file -t v1.0.0
# tag
git cat-file -t sprint-3
# commit
```

`git tag` lists tags. `git show v1.0.0` shows the tag message and the commit.

## Assignment

The handbook is good enough to call 1.0.

1. Annotated tag **`v1.0.0`** on `main`'s current tip.
2. Lightweight tag **`sprint-3`** on the **`A7`** hotfix commit, for the sprint demo.
3. Retroactively, annotated tag **`v0.9.0`** on **`C3`**, the playbook commit, which really was the first usable version.
4. Submit.

The checks confirm each tag's type (annotated `tag` objects for the versions, a plain `commit` reference for `sprint-3`) and which commit each one marks.
