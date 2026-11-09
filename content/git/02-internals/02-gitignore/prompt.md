# Ignore Files with .gitignore

Some files should **never** be committed: passwords, log files, your notes to yourself, big temporary files. If they show up in `git status` every time, sooner or later someone runs `git add .` and commits them.

A file named **`.gitignore`** lists what Git should pretend doesn't exist. One pattern per line:

```
answers.txt     # this exact file
*.log           # every file ending in .log
drafts/         # a whole folder
```

(In a real `.gitignore`, put comments on their own line, starting with `#`.)

Ignored files vanish from `git status`. Ask Git why a file is ignored with:

```sh
git check-ignore -v debug.log
# .gitignore:2:*.log    debug.log
```

The `.gitignore` file itself **is** committed, so everyone on the project ignores the same things.

## Assignment

In `~/git-course/project`:

1. Create two files that shouldn't be committed. In this course, `answers.txt` is where you write down things the checks ask for:

   ```sh
   echo "notes for myself" > answers.txt
   echo "some debug output" > debug.log
   ```

2. Create `.gitignore` that ignores `answers.txt` and every `.log` file:

   ```sh
   printf 'answers.txt\n*.log\n' > .gitignore
   ```

3. Run `git status`: only `.gitignore` shows up.
4. Commit it with the message **`A2: ignore answers and logs`**.
5. Run, then submit `gt.internals.gitignore`.
