# Tidy History

Boss fight! Clean up a messy branch, catch it up with `main` and land it as one straight line.

## Assignment

In `~/git-course/project`:

1. From `main`, create **`add-links`** and write a links page in three messy commits. Copy and paste:

   ```sh
   printf '# Links\n\n- https://git-scm.com\n' > links.md
   git add links.md && git commit -m "wip git"
   echo '- https://docs.github.com' >> links.md && git commit -am "wip github"
   echo '- https://learngitbranching.js.org' >> links.md && git commit -am "wip practice"
   ```

2. Meanwhile, `main` moves on. Switch to `main` and commit:

   ```sh
   echo 'I live by the sea.' >> about.md
   git commit -am "F4: say where I live"
   ```

3. Back on `add-links`: squash the three wip commits into one, **`G4: add the links page`**.
4. Rebase it onto `main`, bring `main` up to it with a fast-forward, and delete `add-links`.
5. Stay on `main`, then submit `gt.rebasing.tidy-history-boss`.
