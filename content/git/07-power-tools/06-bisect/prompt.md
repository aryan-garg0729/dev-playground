# Who Broke the Pager Number?

A customer says nobody answered the pager last night. The number in the handbook is wrong: someone swapped two digits. There are fifteen commits since it was last right, and `git blame` is no help because the line has been rewritten.

You could check them one by one. Or you could let Git **binary search**: mark one commit bad and an older one good, and Git checks out the middle. You say good or bad, it halves the range again. Fifteen commits take about four steps; a thousand take about ten.

```sh
git bisect start <bad> <good>
git bisect good        # or: git bisect bad   (repeat until Git names the culprit)
git bisect reset       # end the session and go back where you were
```

Better still, give it a script that exits **0 for good** and **1 for bad**, and Git runs the whole search itself:

```sh
git bisect run ../bisect-test.sh
# 1f2e3d4c... is the first bad commit
```

The script lives *outside* the repo on purpose: bisect rewrites your working tree at every step, and a script inside the repo could vanish under it.

## Assignment

Set up the incident from `~/elixirbytes/handbook` on a clean `main`. The culprit's position is random, so your answer won't match anyone else's:

```sh
cat > ../bisect-test.sh <<'SH'
#!/bin/sh
# exit 0 = good (the right number is there), 1 = bad
grep -qF '+1-555-0199' oncall.md
SH
chmod +x ../bisect-test.sh

bad=$(( $(od -An -N1 -tu1 /dev/urandom) % 13 + 2 ))
git switch -c audit/history
for i in $(seq -w 1 15); do
  if [ "$i" -eq "$bad" ]; then
    sed 's/+1-555-0199/+1-555-0919/' oncall.md > o.tmp && mv o.tmp oncall.md
  else
    echo "- audit note $i" >> CHANGELOG.md
  fi
  git add -A && git commit -qm "churn $i: handbook edit"
done
git switch main
```

1. `audit/history`'s tip is bad; `main` (where it started) is good. Start a bisect between them.
2. Let `git bisect run ../bisect-test.sh` find the first bad commit. Don't hunt through the log by hand.
3. Add its **full** hash to `answers.txt`:

   ```
   bisect_bad=<the first bad commit>
   ```

4. End the bisect session. You should be back on `main`. Submit.

The checks confirm you're out of bisect mode and on `main`, find the commit that really introduced `+1-555-0919`, and compare it with your answer.
