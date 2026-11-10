# Greedy by Sorting

Many greedy problems are "sort by the right key, then take what you can in one pass". The hard part is choosing the key and convincing yourself (and the interviewer) that it's right.

## The exchange argument

Take any optimal answer that disagrees with your greedy choice. Show that swapping in the greedy choice makes it no worse. Then greedy is optimal too.

Example: to fit the most non-overlapping intervals, take the interval that **ends first**. Any optimal answer starts with some interval that ends at or after yours. Swap it for yours and nothing else breaks, because yours frees the timeline sooner.

## Common keys

| Goal | Sort by |
| --- | --- |
| Most non-overlapping intervals | end time |
| Least total (or average) finishing time | duration, shortest first |
| Merge or cover intervals | start time |
| Pair up items (boats, assignments) | sort, then two pointers from both ends |

## When greedy fails

If a small counterexample breaks your rule, switch to DP. Coin change with coins `{1, 3, 4}` and target 6 is the classic one: greedy takes `4 + 1 + 1`, but `3 + 3` uses fewer coins.

## Check yourself

**You want the largest set of non-overlapping intervals. Which greedy rule is correct?**
