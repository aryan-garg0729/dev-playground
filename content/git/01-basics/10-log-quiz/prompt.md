# Reading the Log

`git log` shows the history. A few options make it easier to read:

```sh
git --no-pager log --oneline --parents
# 7c1e9d0 9a04b31 B1: add the about page
# 9a04b31 A1: add the home page
```

- `--no-pager` prints straight to the terminal instead of a scrolling viewer. (If you end up in the viewer, press `q` to leave.)
- `--oneline` shows one commit per line: a short id and the message.
- `--parents` adds each commit's **parent**, the commit it was built on.

Try it in your project. Your ids will be different.

**What does the output above tell you?**
