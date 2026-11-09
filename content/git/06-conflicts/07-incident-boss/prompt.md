# Incident Review Template

Priya:

> **HB-30: Incident reviews**
>
> After last month we're doing proper incident reviews. Build the template on `feat/incident-review`. Also note in on-call which rota tool owns the pager. Heads up: the pager migration finishes today, so the number on main is about to change under you.
>
> When it lands I want one clean commit, rebased on the new number, and then a **real merge commit** so the history shows where the feature entered main, like a pull request would.

## Assignment

Everything you need is from this chapter and the rebase chapter. The order matters: squash **first**, rebase **second**. Three wip commits rebased over a conflict means resolving it up to three times; one squashed commit means once.

1. From `main`, create **`feat/incident-review`** and make three wip commits:

   ```sh
   printf '# Incident review\n\nWrite one within 48 hours of every SEV1 or SEV2.\n' > incidents.md
   git add incidents.md && git commit -m "wip template"
   printf '\n## Timeline\n\n- Detected:\n- Mitigated:\n- Resolved:\n' >> incidents.md
   git commit -am "wip timeline"
   sed 's/^Primary pager: .*/Primary pager: +1-555-0142, Pagerly rota platform/' oncall.md > o.tmp && mv o.tmp oncall.md
   git commit -am "wip pager note"
   ```

2. Meanwhile, the migration lands on `main`:

   ```sh
   git switch main
   sed 's/^Primary pager: .*/Primary pager: +1-555-0199/' oncall.md > o.tmp && mv o.tmp oncall.md
   git commit -am "H6: move the pager to the new number"
   ```

3. On the branch, squash the three wip commits into **`I6: add the incident review template`**.
4. Rebase it onto `main`. Resolve the one conflict so the line is exactly:

   ```
   Primary pager: +1-555-0199, Pagerly rota platform
   ```

5. Merge it into `main` as a real merge commit with a message starting **`J6:`**. After the rebase, `main` *could* fast-forward; `git merge --no-ff` makes Git create the merge commit anyway.
6. Submit.

Each check is a hit on the boss:

- exactly one `I6`, and no `wip` anywhere
- `I6`'s only parent is `H6`
- `J6` merges `H6` and `I6` at the tip of `main`
- the pager line has the new number and your rota note, with no markers
- `incidents.md` has both sections

Defeat it to unlock **Power tools**.
