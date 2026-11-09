# Fix an Old Message

Fixing the message of your **last** commit is easy: `git commit --amend -m "better message"`. To fix an **older** one, use interactive rebase with `reword`.

`HEAD~2` means "two commits before HEAD". So `git rebase -i HEAD~2` lists just your last two commits:

```
pick 4d5e6f7 D4: stuff
pick 8a9b0c1 E4: link the privacy page
```

Change the first `pick` to `reword` (or `r`), save and close, and Git opens the editor for that message.

## Assignment

In `~/git-course/project`, on `main`:

1. Make two commits, the first with a useless message. Copy and paste:

   ```sh
   printf '# Privacy\n\nThis site uses no cookies.\n' > privacy.md
   git add privacy.md && git commit -m "D4: stuff"
   echo '- [Privacy](privacy.md)' >> index.md
   git commit -am "E4: link the privacy page"
   ```

2. Nobody else has these commits, so rewriting them is safe. Change `D4`'s message to exactly **`D4: add the privacy page`**. Keep `E4` as it is.
3. Check with `git --no-pager log --oneline -n 3`.
4. Run, then submit `gt.rebasing.reword`.
