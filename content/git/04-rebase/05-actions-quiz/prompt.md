# Pick, Squash, Fixup

Each line of the interactive rebase to-do list starts with an action:

| Action | What it does |
| --- | --- |
| `pick` | keep the commit as it is |
| `reword` | keep the changes, edit the message |
| `edit` | stop after this commit so you can change it |
| `squash` | fold into the commit above, and combine both messages |
| `fixup` | fold into the commit above, and throw this one's message away |
| `drop` | remove the commit entirely |

Your branch has two commits:

```
pick 3c1d2e4 A: add the menu page
pick 7f8a9b0 fix typo in menu
```

**You want one commit called "A: add the menu page", with the typo fixed. What do you change?**
