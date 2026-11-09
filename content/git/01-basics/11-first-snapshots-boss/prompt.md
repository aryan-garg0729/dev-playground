# First Snapshots

Boss time! Make two separate, tidy commits using only what this chapter taught: create, stage exactly what belongs together, commit.

## Assignment

In `~/git-course/project`:

1. Create the contact page **and** add the link to the home page. Copy and paste:

   ```sh
   cat > contact.md <<'EOF'
   # Contact

   Email: hello@example.com
   Office hours: Mondays, 10:00

   ## Social

   - Blog: example.com/blog
   EOF
   echo '- [Contact](contact.md)' >> index.md
   ```

2. Make **two** commits, one after the other:
   - **`C1: add the contact page`** with **only** `contact.md`
   - **`D1: link the contact page`** with **only** the change to `index.md`
3. Finish with nothing left uncommitted.
4. Submit `gt.basics.first-snapshots-boss`.

Hint: `git add` takes file names, so you can stage one file at a time.
