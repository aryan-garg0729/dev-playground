# Tell Git Who You Are

Every commit records **who** made it. So before your first commit, tell Git your name and email:

```sh
git config --global user.name "Your Name"
git config --global user.email "you@example.com"
```

`--global` means "for every project on this computer", so you only do this once. The settings are saved in a file in your home folder, `~/.gitconfig`.

One more setting: the first **branch** of a new project (you'll learn about branches soon) is traditionally called `main`. Older Git versions call it `master`. Let's all agree on `main`:

```sh
git config --global init.defaultBranch main
```

See what you've set with:

```sh
git config --global --list
```

## Assignment

1. Set your name and email (use your real ones if you'll ever share your work).
2. Set the default branch name to `main`.
3. Run, then submit `gt.basics.config`.
