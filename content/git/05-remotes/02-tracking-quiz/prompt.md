# What Is origin/main?

After the push, a new name shows up:

```sh
git branch -a
# * main
#   ...
#   remotes/origin/main

git --no-pager log --oneline --decorate -n 1
# 5bef8fc (HEAD -> main, origin/main) I4: restore the front page
```

`origin/main` is a **remote-tracking branch**. Your repo stores it in `.git/refs/remotes/origin/main`.

**What does it represent?**
