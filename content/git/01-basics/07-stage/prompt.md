# Stage a File

Our project will be a tiny website made of text pages. Let's create the home page and watch it move through the areas.

A new file starts **untracked**: Git sees it but isn't saving it.

```sh
git status
# Untracked files:
#         index.md
```

`git add` puts it in the staging area:

```sh
git add index.md
git status
# Changes to be committed:
#         new file:   index.md
```

`git status --short` gives the same news in one line per file: `??` means untracked, `A` means added (staged).

## Assignment

In `~/git-course/project`:

1. Create the home page. Copy and paste:

   ```sh
   printf '# My Website\n\nWelcome to my website.\n' > index.md
   ```

2. Run `git status`: it's untracked.
3. Stage it with `git add index.md`, and run `git status` again.
4. **Don't commit yet.** Run, then submit `gt.basics.stage`.
