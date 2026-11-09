# Dex Gets a Copy

`git clone` copies a whole repository, every commit of it, and sets `origin` to point back where it came from:

```sh
cd ~/elixirbytes
git clone handbook-origin.git dex-handbook
# Cloning into 'dex-handbook'...
```

A clone is a full repo, not a link. Dex can commit, branch and read history with no server at all, and only talks to the remote when he pushes or pulls.

## Assignment

You'll play Dex in his own clone for a few lessons. Real teammates have their own laptops; here it's a second folder.

1. Clone `~/elixirbytes/handbook-origin.git` into **`~/elixirbytes/dex-handbook`**.
2. Inside `dex-handbook`, give it Dex's identity as **local** config: name `Dex`, email **`dex@elixirbytes.shop`**.
3. As Dex, add the contacts page and commit it as **`A5:`**:

   ```sh
   cat > contacts.md <<'MD'
   # Escalation contacts

   - Platform: Priya (priya@elixirbytes.shop)
   - Brewing floor: ext. 4410
   MD
   ```

4. Push it. Dex's `main` already tracks `origin/main` because `clone` set that up, so plain `git push` works.
5. Go back to `~/elixirbytes/handbook` (**don't pull yet**) and submit.

The checks confirm Dex's clone has his email and came from the shared remote, and that the remote's `main` now ends with `A5` by `dex@elixirbytes.shop`.
