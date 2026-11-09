# Clean History Friday

Priya, Friday afternoon:

> **HB-12: Glossary**
>
> Dex keeps saying "MTTR" in incident calls and Mo keeps nodding. We need a glossary. I'm also changing CONTRIBUTING on main right now, so you'll be behind me by the time you're done. When it lands I want **one** commit, on top of my change, in a straight line. No wip, no merge bubble, no leftover branch.

## Assignment

Everything you need is from this chapter.

1. From `main`, create **`docs/glossary`** and write the glossary in three commits whose messages start with **`wip`**:

   ```sh
   printf '# Glossary\n\n- SLO: the reliability target we promise.\n' > glossary.md
   git add glossary.md && git commit -m "wip slo"
   echo "- MTTR: mean time to recovery." >> glossary.md && git commit -am "wip mttr"
   echo "- Canary: a small first deploy that warns us early." >> glossary.md && git commit -am "wip canary"
   ```

2. Meanwhile, Priya commits on `main`. Switch to `main` and make her change:

   ```sh
   echo "- Explain acronyms the first time you use them." >> CONTRIBUTING.md
   git commit -am "F3: ask for acronyms to be explained"
   ```

3. Back on `docs/glossary`: squash the three wip commits into one, **`G3: add the glossary`**.
4. Rebase it onto `main`, fast-forward `main`, and delete `docs/glossary`.
5. Stay on `main` and submit.

Each check is a hit on the boss:

- exactly one `G3`, and no `wip` commit on any branch
- `G3`'s parent is `F3`, and `G3` is `main`'s tip
- no merge commit since `F3`
- the glossary has SLO, MTTR and Canary
- `docs/glossary` is gone

Defeat it to unlock **Undoing things**.
