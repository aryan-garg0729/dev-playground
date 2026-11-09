# A Shared Copy

So far the project lives only on your computer. To work with other people, you need a shared copy that everyone sends commits to and gets commits from. That's a **remote**.

On a real team the remote lives on GitHub, GitLab or a company server. Underneath, it's just another Git repository, usually a **bare** one: only the `.git` part, with no working files, because nobody edits files on the server. We'll make one in a folder next to your project, so nothing in this chapter needs the internet (GitHub comes later in the chapter).

```sh
git init --bare -b main ~/git-course/project-origin.git
```

Tell your repo about it, under the usual name `origin`, then **push** (send) your `main` there:

```sh
git remote add origin ~/git-course/project-origin.git
git push -u origin main
#  * [new branch]      main -> main
# branch 'main' set up to track 'origin/main'.
```

`-u` links your `main` to `origin/main`. From then on, plain `git push` and `git pull` know where to go.

## Assignment

1. Create the bare repository **`~/git-course/project-origin.git`** with `-b main`, as above.
2. In `~/git-course/project`, add it as the remote **`origin`**.
3. Push `main` with `-u`.
4. Peek inside: `ls ~/git-course/project-origin.git`. It looks like your `.git` folder, because that's what it is.
5. Run, then submit `gt.sharing.publish` from `~/git-course/project`.

Only `main` was pushed. Your other branches stay private until you push them.
