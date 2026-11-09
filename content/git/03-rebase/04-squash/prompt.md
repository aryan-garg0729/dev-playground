# Squash the Mess

Real work is messy. You write a bit, commit "wip", write more, commit "wip" again. That's fine on your own branch. But reviewers, and you in six months, want **one commit per idea** with a message that explains it.

**Interactive rebase** lets you rewrite your branch's commits before anyone sees them:

```sh
git rebase -i main
```

Git opens your editor with a to-do list of the commits it's about to replay, **oldest first**:

```
pick 1a2b3c4 wip triage
pick 5d6e7f8 wip mitigate
pick 9a0b1c2 wip comms
```

Change a word, save and close, and Git does what the list says. To fold the last two into the first, change their `pick` to `squash` (or just `s`):

```
pick 1a2b3c4 wip triage
squash 5d6e7f8 wip mitigate
squash 9a0b1c2 wip comms
```

Git then opens the editor once more for the combined commit's message. Delete the old lines and write the new one.

Not sure which editor Git opens? Set one for this repo: `git config --local core.editor nano` (or `"code --wait"` for VS Code).

## Assignment

After last week's outage, Priya wants an incident playbook. Write it the way people really do, then clean it up.

1. From `main`, create **`docs/playbook`**. Make **three** commits, each with a message starting **`wip`**:

   ```sh
   printf '# Incident playbook\n\n## Triage\n\nWho is on call? Check the pager first.\n' > playbook.md
   git add playbook.md && git commit -m "wip triage"
   printf '\n## Mitigate\n\nRoll back before you debug.\n' >> playbook.md
   git commit -am "wip mitigate"
   printf '\n## Comms\n\nPost in #status every 30 minutes.\n' >> playbook.md
   git commit -am "wip comms"
   ```

2. Squash all three into **one** commit with the message **`C3: add the incident playbook`**, using `git rebase -i main`.
3. Fast-forward `main` to the branch.
4. Run `git --no-pager log --oneline -n 3`: `C3` directly on `B3`, no `wip` anywhere.
5. Stay on `main` and submit.

The checks confirm there's exactly one `C3` commit and no `wip` commit on any branch, that `C3` sits on `B3` at the tip of `main`, and that the playbook has all three sections.
