# A Server in a Folder

So far the handbook lives only on your laptop. Priya wants Dex and Mo working on it too, which means a shared copy everyone sends commits to and gets commits from: a **remote**.

On a real team that's GitHub, GitLab or a company server. Underneath, it's just another Git repository, usually a **bare** one: a `.git` folder with no working tree, because nobody edits files on the server. We'll make one in a folder next to the handbook, so nothing in this chapter needs the internet.

```sh
git init --bare -b main ~/elixirbytes/handbook-origin.git
```

Then tell your repo about it, under the conventional name `origin`, and **push** (send) your `main` there:

```sh
git remote add origin ~/elixirbytes/handbook-origin.git
git push -u origin main
# To /home/sam/elixirbytes/handbook-origin.git
#  * [new branch]      main -> main
# branch 'main' set up to track 'origin/main'.
```

`-u` (short for `--set-upstream`) links your `main` to `origin/main`. From then on, plain `git push` and `git pull` on `main` know where to go.

## Assignment

1. Create the bare repo **`~/elixirbytes/handbook-origin.git`**, with `-b main` so its default branch matches yours (otherwise clones of it start on a branch that doesn't exist).
2. Add it as the remote **`origin`** of your handbook.
3. Push `main` with upstream tracking.
4. Run `ls ~/elixirbytes/handbook-origin.git`. It looks like the inside of your `.git` folder, because that's what it is.
5. Submit from `~/elixirbytes/handbook`.

Only `main` is pushed. Your other branches stay private until you push them.

The checks confirm `handbook-origin.git` is bare, `origin` points at it, its `main` is the same commit as yours, and your `main` tracks `origin/main`.
