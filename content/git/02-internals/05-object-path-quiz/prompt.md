# Where Objects Live

Every object is a file inside `.git/objects`. Have a look (it's safe to look, just don't edit):

```sh
ls .git/objects
# 1f  4b  55  7c  9a  info  pack
ls .git/objects/9a
# 04b31c5e2d8f7a6b1c0e9d8f7a6b5c4d3e2f10
```

The contents are compressed, so `cat` shows gibberish. That's what `git cat-file` is for.

**Where is the object with the id `9a04b31c5e2d8f7a6b1c0e9d8f7a6b5c4d3e2f10` stored?**
