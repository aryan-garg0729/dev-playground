# N - Slimes

**Practice:** [N - Slimes on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_n)

## The task

n slimes sit in a row. Merging two neighbours costs the sum of their sizes. Merge everything into one slime at the least total cost: interval DP (and the base for Knuth's optimization).

## Before you code

- `dp[l][r]` = cheapest way to merge `a[l..r]` into one slime = min over splits k of `dp[l][k] + dp[k+1][r]`, plus `sum(l..r)` from prefix sums.

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**In which order must you fill `dp[l][r]`?**
