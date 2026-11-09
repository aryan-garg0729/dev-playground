# A Branch Is a Sticky Note

Mo wants to rework the handbook's front page "without breaking anything". The usual answer is a **branch**: a separate line of history you can commit to while `main` stays as it is.

People often imagine a branch as a copy of the project. Look at what Git actually stores:

```sh
ls .git/refs/heads
# main

cat .git/refs/heads/main
# 4ad50fb9e1c03b7f6a2d8e5c1b0a9f7e3d2c6b18

git log -1 --format=%H
# 4ad50fb9e1c03b7f6a2d8e5c1b0a9f7e3d2c6b18
```

Try it in your repo. The hash in the file is the hash of your latest commit.

**What does creating a new branch store?**
