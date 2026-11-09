# Tag a Release

Branches move. A **tag** is a name for a commit that **stays put**: perfect for "this is what we released as version 1.0".

There are two kinds:

```sh
git tag today                          # lightweight: just a name for a commit
git tag -a v1.0.0 -m "Version 1.0"     # annotated: also records who, when and a message
git tag -a v0.9.0 -m "..." 7f8e9d0     # tag any commit, not just the current one
```

Releases use **annotated** tags, because they record who made the release, when, and why. You can see the difference:

```sh
git cat-file -t v1.0.0
# tag
git cat-file -t today
# commit
```

`git tag` lists tags. `git show v1.0.0` shows the tag and its commit. (Tags aren't pushed by default: `git push origin v1.0.0` sends one.)

Version numbers usually follow `MAJOR.MINOR.PATCH`: a bug fix bumps the last number, a new feature the middle one, a big breaking change the first.

## Assignment

In `~/git-course/project`, the website is good enough to call 1.0:

1. Annotated tag **`v1.0.0`** on `main`'s current tip.
2. Lightweight tag **`today`** on the **`A8`** commit.
3. Looking back, the hobbies page (**`C4`**) was really the first usable version: annotated tag **`v0.9.0`** on `C4`.
4. Run, then submit `gt.toolbox.tags`.
