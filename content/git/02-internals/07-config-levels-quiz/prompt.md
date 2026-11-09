# Which Setting Wins?

`git config` settings live at three levels:

| Level | Flag | File | Applies to |
| --- | --- | --- | --- |
| system | `--system` | `/etc/gitconfig` | every user on the computer |
| global | `--global` | `~/.gitconfig` | every repo of yours |
| local | `--local` (the default) | `.git/config` | just this repo |

A common case: you use your personal email everywhere, but one project is for work.

```sh
git config --global user.email "you@personal.example"
cd ~/work-project
git config --local user.email "you@work.example"
```

**You commit in `~/work-project`. Which email does the commit get?**
