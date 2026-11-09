# Merge or Rebase?

|  | Merge | Rebase |
| --- | --- | --- |
| History | kept as it happened, with merge commits | rewritten into a straight line |
| Hashes | existing commits untouched | replayed commits get new hashes |
| Safe on shared branches | yes | no |
| Conflicts | resolved once, in the merge commit | possibly once per replayed commit |

Teams usually combine them: **rebase your own branch** to stay current and tidy, then **merge it into `main`** (often through a pull request).

**Which of these are good practice? Pick all that apply.**
