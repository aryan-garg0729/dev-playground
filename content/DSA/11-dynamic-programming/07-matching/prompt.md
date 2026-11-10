# O - Matching

**Practice:** [O - Matching on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_o)

## The task

Count the perfect matchings between n men and n women (n ≤ 21), given who is compatible with whom, modulo 10⁹+7: bitmask DP over subsets.

## Before you code

- `dp[mask]` = ways to match the first `popcount(mask)` men with exactly the women in mask.

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why does the state need only the mask, and not also which man is next?**
