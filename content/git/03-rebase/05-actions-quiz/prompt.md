# Pick, Squash, Fixup

Each line of an interactive rebase's to-do list starts with an action:

| Action | What it does |
| --- | --- |
| `pick` | keep the commit as it is |
| `reword` | keep the changes, edit the message |
| `edit` | stop after this commit so you can change it |
| `squash` | fold into the commit above, and combine both messages |
| `fixup` | fold into the commit above, and discard this one's message |
| `drop` | remove the commit entirely |

Your branch has two commits:

```
pick 3c1d2e4 A: add the potion menu
pick 7f8a9b0 fix typo in menu
```

**You want one commit with the message "A: add the potion menu" and the typo fixed. What do you change?**
