# Reading the Graph

`--graph` draws the history as lines. Add `--parents` to see each commit's parents:

```
git --no-pager log --oneline --graph --parents
* 3e9b7d2 6d14c0a Update FAQ
*   6d14c0a 51ab3f9 2c7e880 Merge branch 'add-gallery'
|\
| * 2c7e880 a1b2c3d Add gallery
* | 51ab3f9 a1b2c3d Link projects
|/
* a1b2c3d 9f8e7d6 Add projects page
```

Each `*` is a commit. `|` lines connect a commit to its parents, and where two lines meet, two histories join.

**Which commit is the merge commit?**
