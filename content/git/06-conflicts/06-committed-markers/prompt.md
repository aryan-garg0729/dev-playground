# Markers in Production

Conflict markers are ordinary text. If you `git add` a file without deleting them, Git commits them without a word, and the next person sees `<<<<<<< HEAD` in the published page.

To find them anywhere in a commit:

```sh
git grep -n -e '^<<<<<<<' -e '^=======$' -e '^>>>>>>>' main
# main:index.md:12:<<<<<<< HEAD
```

`git grep` searches the files *in a commit*, not your working tree, and exits with status 1 when it finds nothing.

## Assignment

Mo resolved a conflict on your laptop by adding the file and committing. Recreate it on `main`:

```sh
printf '<<<<<<< HEAD\nOpening hours: Monday to Saturday, 09:00-17:00.\n=======\nOpening hours: every day, 08:00-18:00.\n>>>>>>> mo/hours\n' >> index.md
git -c user.name=Mo -c user.email=mo@elixirbytes.shop commit -am "Merge Mo's opening hours"
```

1. Use `git grep` to find the markers.
2. Priya settles it: the shop is open **Monday to Saturday, 08:00-18:00**. Make the front page say exactly `Opening hours: Monday to Saturday, 08:00-18:00.` with no markers and no other opening-hours line from Mo's mess.
3. Commit the fix as **`G6:`** on top of Mo's commit. Don't rewrite his commit.
4. Submit.

The checks confirm `git grep` finds no markers on `main`, the front page has the agreed line (and not Mo's "every day" one), and `G6` sits on Mo's commit.
