# The First Snapshot

Time to move the handbook in. A **commit** is a snapshot of the staging area plus a little metadata: who, when, and a message saying why.

```sh
git add index.md
git commit -m "start the handbook"
# [main (root-commit) 3f2a1bc] start the handbook
#  1 file changed, 1 insertion(+)
```

`3f2a1bc` is the start of the commit's **hash**, its unique id. Yours will be different, because the hash covers the content, your name and the exact time. (How the hash is computed is out of scope. Just know it's a fingerprint of the commit.)

## Labels

Because hashes differ for every learner, the checks find your commits by **label**: each commit message in this course starts with a letter and the chapter number, like `A1:`, `B1:`, `A2:`. The rest of the message is yours to write, but the label must be exact.

## Assignment

Copy the handbook out of the old doc. Paste this into your terminal in `~/elixirbytes/handbook` to create the two pages:

```sh
cat > index.md <<'MD'
# ElixirBytes Handbook

How the potion shop stays open.

- [On-call](oncall.md)
MD

cat > oncall.md <<'MD'
# On-call

Primary pager: +1-555-0142
Escalation window: 5 minutes
Handover: Mondays, 10:00

## Rotation

One week each, Monday to Monday.
MD
```

Later chapters edit these exact lines, so don't reword them.

1. Run `git status`. Both files are untracked.
2. Stage both with `git add` and run `git status` again.
3. Commit them with a message starting **`A1:`**, for example `git commit -m "A1: start the handbook"`.
4. Run `git status` once more. It should say `nothing to commit, working tree clean`.
5. Submit.

The checks look for a commit labelled `A1:` containing both files, made with your `@elixirbytes.shop` email, and a clean working tree.
