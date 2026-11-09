# Throw Away One Edit

Every Git user eventually types something they regret. This chapter is about getting back, and about knowing which "undo" you need, because Git has several and they do very different things.

The smallest one: you edited a file and want it back the way it was in the last commit.

```sh
git restore oncall.md
```

`git restore` copies a file **from a source into your working tree**. The source defaults to the staging area, which matches the last commit unless you've staged something. Only the files you name are touched.

Careful: the edits you throw away were never committed, so Git has no copy. They're gone for good.

## Assignment

Set up the situation: you're writing a new rule while Mo "helps" with the on-call page on your laptop.

```sh
echo "Ask Priya before renaming files." >> CONTRIBUTING.md
echo "TODO: lol idk" > oncall.md
git status
```

1. Throw away Mo's change to **`oncall.md`** only.
2. **Keep** your `CONTRIBUTING.md` edit. Don't commit or stage it yet; the next lesson needs it. (`git restore .` would throw that away too.)
3. Run `git diff` to see what's left, then submit.

The checks confirm `oncall.md` matches the last commit and your `CONTRIBUTING.md` line is still there, unstaged.
