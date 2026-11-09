# What Is origin/main?

After the push, a new name shows up:

```sh
git branch -a
# * main
#   remotes/origin/main

git --no-pager log --oneline -n 1
# 5bef8fc (HEAD -> main, origin/main) I5: undo the redesign
```

`origin/main` is called a **remote-tracking branch**.

**What does it represent?**
