# A Setting for One Repo

Let's try the levels for real. `--show-origin` tells you which file a setting came from:

```sh
git config --show-origin --get user.email
# file:/Users/you/.gitconfig    you@example.com
```

## Assignment

In `~/git-course/project`:

1. Set a **local** email for this repo only, ending in `@work.example` (for example `you@work.example`):

   ```sh
   git config --local user.email "you@work.example"
   ```

2. Check where it comes from now: `git config --show-origin --get user.email`.
3. Make a commit with it. Copy and paste:

   ```sh
   echo 'I also like cooking.' >> about.md
   git add about.md
   git commit -m "B2: mention cooking"
   ```

4. Compare the authors: `git --no-pager log --format='%h %ae %s'`.
5. Run, then submit `gt.internals.local-config`.

Other projects on your computer still use your global email.
