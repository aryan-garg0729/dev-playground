# Put It on GitHub

**GitHub** is a website that hosts Git repositories. It's where most open-source projects live, and where most companies keep their code. Under the hood it's exactly what you just built: a remote you push to and pull from, plus a website around it for browsing code, reviewing changes and discussing them.

A repo can have **several remotes**. You'll keep `origin` (your local practice remote) and add GitHub as a second one, called `github`.

## Assignment

1. Create a free account at [github.com](https://github.com) if you don't have one.
2. Let Git sign in to GitHub. The easiest way is GitHub's own CLI, `gh`:

   ```sh
   brew install gh        # macOS; on Linux see cli.github.com
   gh auth login          # choose GitHub.com, HTTPS, and log in with a browser
   ```

3. On GitHub, create a **new, empty repository** called `my-website` (no README, no license: it must be empty). It can be public or private.
4. In `~/git-course/project`, add it as the remote **`github`** and push `main`:

   ```sh
   git remote add github https://github.com/YOUR_USERNAME/my-website.git
   git push github main
   ```

5. Refresh the repository page on GitHub. Your pages and your whole history are there.
6. Run, then submit `gt.sharing.github`.

The last check asks GitHub which commit its `main` is at, and compares it with yours.
