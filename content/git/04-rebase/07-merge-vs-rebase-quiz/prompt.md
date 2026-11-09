# Merge or Rebase?

Both bring changes together. They just tell different stories:

| | Merge | Rebase |
| --- | --- | --- |
| History | Keeps what really happened, with a merge commit where branches join | Rewrites it into one straight line |
| Rewrites commits? | No | Yes: makes new copies |
| Safe on shared branches? | Yes | No (the golden rule) |

Teams argue about this a lot. Most settle on something like: rebase your own branch while you work, then merge (or fast-forward) it into `main`.

**Which of these are good uses? Pick all that apply.**
