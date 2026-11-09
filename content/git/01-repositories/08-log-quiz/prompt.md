# Reading the Log

`git log` shows history. The full output is long, so a few flags make it readable:

```sh
git --no-pager log --oneline --parents
# 7c1e9d0 9a04b31 B1: add the deploy runbook and ignore answers.txt
# 9a04b31 A1: start the handbook
```

- `--no-pager` prints straight to the terminal instead of opening a scrolling viewer (press `q` to leave the viewer if you forget).
- `--oneline` shows one commit per line: a short hash and the first line of the message.
- `--parents` adds each commit's parent hash after its own.

Try `git --no-pager log --oneline --parents` in your repo now. Your hashes will differ.

**What does the output above tell you?**
