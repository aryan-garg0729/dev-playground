# Find the Bad Commit

Visitors say your email address on the contact page is wrong: someone swapped two letters. There are fifteen commits since it was last right. Which one broke it?

You could check them one by one. Or let Git do a **binary search**: you mark one commit bad and an older one good, and Git checks out the one in the middle. You say good or bad, and it halves the range again. Fifteen commits take about four steps; a thousand take about ten.

```sh
git bisect start BAD GOOD
git bisect good        # or: git bisect bad   (repeat until Git names the culprit)
git bisect reset       # end the session and go back to where you were
```

Better still, give it a script that exits **0 for good** and **1 for bad**, and Git runs the whole search by itself:

```sh
git bisect run ../bisect-test.sh
# 1f2e3d4c... is the first bad commit
```

The script lives **outside** the repo on purpose: bisect changes your files at every step, and a script inside the repo could change too.

## Assignment

Set it up from `~/git-course/project` on a clean `main`. The bad commit is in a random spot, so your answer won't match anyone else's. Copy and paste:

```sh
cat > ../bisect-test.sh <<'EOF'
#!/bin/sh
# exit 0 = good (the right email is there), 1 = bad
grep -qF 'hi@example.com' contact.md
EOF
chmod +x ../bisect-test.sh

bad=$(( $(od -An -N1 -tu1 /dev/urandom) % 13 + 2 ))
git switch -c audit/history
for i in $(seq -w 1 15); do
  if [ "$i" -eq "$bad" ]; then
    sed 's/hi@example.com/hi@exmaple.com/' contact.md > c.tmp && mv c.tmp contact.md
  else
    echo "- note $i" >> notes.md
  fi
  git add -A && git commit -qm "change $i"
done
git switch main
```

1. `audit/history`'s tip is bad; `main`, where it started, is good. Start a bisect between them.
2. Let `git bisect run ../bisect-test.sh` find the first bad commit. Don't hunt through the log by hand.
3. Add its **full** id to `answers.txt` as `bisect_bad=COMMIT_ID`.
4. End the bisect session; you should be back on `main`. Run, then submit `gt.toolbox.bisect`.
