# Three Kinds of Reset

`git reset` moves your branch back to an earlier commit, as if the later commits never happened. `HEAD~1` means "one commit back":

```sh
git reset --soft  HEAD~1
git reset --mixed HEAD~1   # the default if you leave out the mode
git reset --hard  HEAD~1
```

They all move the branch. They differ in what happens to the changes from the commits you undid:

| Mode | Branch moves back | Changes are... |
| --- | --- | --- |
| `--soft` | Yes | still **staged**, ready to commit again |
| `--mixed` | Yes | still in your files, **unstaged** |
| `--hard` | Yes | **gone** from your files too |

`--hard` is the dangerous one: any uncommitted work in your files is lost.

**Which reset leaves your files exactly as they were in the earlier commit, as if you never made the later one?**
