# Markers Got Committed

Conflict markers are ordinary text. If someone runs `git add` without deleting them, Git commits them without a word, and the next person sees `<<<<<<< HEAD` on the live website.

Find markers anywhere in a commit with `git grep`, which searches the files **in a commit**:

```sh
git grep -n -e '^<<<<<<<' -e '^=======$' -e '^>>>>>>>' main
# main:hobbies.md:6:<<<<<<< HEAD
```

It exits with status 1 when it finds nothing.

## Assignment

In `~/git-course/project`, on `main`, recreate a teammate's sloppy merge. Copy and paste:

```sh
printf '<<<<<<< HEAD\nFavourite food: pizza.\n=======\nFavourite food: sushi.\n>>>>>>> food\n' >> hobbies.md
git -c user.name=Teammate -c user.email=teammate@example.com commit -am "Merge favourite food"
```

1. Use `git grep` to find the markers.
2. Fix `hobbies.md` so it says exactly `Favourite food: ramen.` (a compromise), with no markers and no pizza or sushi line.
3. Commit the fix as **`G7: fix the favourite food`** on top of the bad commit. Don't rewrite it.
4. Run, then submit `gt.resolving.committed-markers`.
