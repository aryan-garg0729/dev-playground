# Where Objects Live

Try this in your repo:

```sh
git log -1 --format=%H ':/^A1:'
# (prints A1's full 40-character hash)

ls .git/objects
# 1b  7c  9a  d5  ec  info  pack
```

`':/^A1:'` means "the newest commit whose message starts with `A1:`", which is how the checks in this course find your commits.

Now list one of those two-character folders and compare what you see with A1's hash.

**Which pattern matches where Git stored commit `9a04b31c5e...`?**
