# Make It a Repository

Right now `~/elixirbytes/handbook` is an ordinary folder. Git doesn't watch folders unless you ask it to. Asking is one command:

```sh
git init -b main
# Initialized empty Git repository in /home/you/elixirbytes/handbook/.git/
```

That creates a hidden `.git` folder. **That folder *is* the repository**: every snapshot, branch and setting lives inside it. The files next to it are your **working tree**, the copy you edit. Delete `.git` and you have a plain folder again (please don't).

```sh
ls -a
# .  ..  .git
```

`-b main` names the first branch `main`. Without it, older setups call it `master`. We'll cover branches properly in chapter 2. For now, a branch is "the line of history you're on".

## Assignment

1. In `~/elixirbytes/handbook`, create a repository whose first branch is **`main`**.
2. Peek inside with `ls .git`. You'll see `HEAD`, `config`, `objects` and `refs`. We'll open several of these later.
3. Submit from `~/elixirbytes/handbook`.

The checks confirm the handbook folder is the *top* of its own repository (not a folder inside some other repo, such as one in your home directory) and that you're on `main`.

If you ran `git init` in the wrong place, remove that stray `.git` folder (only that one) and try again.
