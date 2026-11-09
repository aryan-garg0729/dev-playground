# Build Up History

One commit is a start. History gets useful as commits pile up. Before you commit, it's worth checking exactly what changed:

```sh
git diff            # changes you haven't staged yet
git diff --staged   # changes that are staged
```

`+` lines were added, `-` lines were removed:

```diff
 # My Website

 Welcome to my website.
+- [About](about.md)
```

To stage **every** change in the folder at once, use `git add .` (the dot means "this folder"). Handy, but run `git status` first so nothing sneaks in.

## Assignment

In `~/git-course/project`:

1. Add an about page and link it from the home page. Copy and paste:

   ```sh
   printf '# About\n\nI am learning Git.\n' > about.md
   echo '- [About](about.md)' >> index.md
   ```

2. Look at the change with `git diff`. (The new file doesn't show yet: `git diff` only shows changes to files Git already tracks.)
3. Stage both files with `git add .`, check with `git diff --staged`, then commit with the message **`B1: add the about page`**.
4. Run, then submit `gt.basics.history`.
