# Talking to Git

Every Git command has the same shape:

```sh
git <command> [options]
```

For example:

```sh
git version          # command: version
git status --short   # command: status, option: --short
```

Git has about 150 commands. You'll use around 15 of them daily. The rest are low-level building blocks that the everyday commands use under the hood. People call the friendly everyday ones **porcelain** (like a shiny sink) and the low-level ones **plumbing** (the pipes behind the wall). You'll peek at a little plumbing in the next chapter.

Nobody remembers every option. Git has a manual built in:

```sh
git help             # the most common commands
git <command> --help # the full manual for one command (press q to quit)
```

**How do you open the manual for `git status`?**
