# Which Config Wins?

Every commit records **who** made it: a name and an email. Git looks them up in its config, which has three levels:

| Level | Applies to | Lives in | Set with |
| --- | --- | --- | --- |
| system | every user on the machine | `/etc/gitconfig` | `git config --system ...` |
| global | every repo of *your* user | `~/.gitconfig` | `git config --global ...` |
| local | just this one repo | `.git/config` | `git config --local ...` |

You can see every value and the file it came from:

```sh
git config --list --show-origin
# file:/home/you/.gitconfig   user.email=you@personal.dev
# file:.git/config            user.email=you@elixirbytes.shop
```

You set `user.email` to `you@personal.dev` with `--global` months ago. Today, inside the handbook repo, you run `git config --local user.email you@elixirbytes.shop`.

**Which email does your next commit in the handbook use?**
