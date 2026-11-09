# Commit as an ElixirBytes Dev

Priya has one rule for the handbook: every change must be traceable to a person with an `@elixirbytes.shop` address. Your other projects should keep whatever identity they use today, so a global change is out.

That's what the **local** level is for. Run these inside the repo:

```sh
git config --local user.name "Sam Rivera"
git config --local user.email "sam@elixirbytes.shop"
```

They write to `.git/config`, a plain text file you can read:

```sh
cat .git/config
# [core]
#         repositoryformatversion = 0
#         ...
# [user]
#         name = Sam Rivera
#         email = sam@elixirbytes.shop
```

Ask Git which value it will actually use, and where it found it:

```sh
git config --show-origin --get user.email
# file:.git/config        sam@elixirbytes.shop
```

## Assignment

1. In the handbook repo, set a **local** `user.name` (your name) and a **local** `user.email` ending in **`@elixirbytes.shop`** (use your first name, like `sam@elixirbytes.shop`; it doesn't need to be a real mailbox).
2. Confirm with `git config --show-origin --get user.email` that the value comes from `.git/config`.
3. Submit from `~/elixirbytes/handbook`.

The checks read the local name and email, and confirm the email Git will use comes from `.git/config`. Don't use `--global` here: that would change every repo on your machine.
