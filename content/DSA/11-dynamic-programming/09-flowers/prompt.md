# Q - Flowers

**Practice:** [Q - Flowers on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_q)

## The task

Flowers stand in a row, each with a distinct height and a beauty. Keep a subsequence whose heights strictly increase, maximizing the total beauty: LIS optimized with a Fenwick or segment tree.

## Before you code

- `dp[h]` = best total ending with the flower of height h. For each flower: `dp[h] = beauty + max(dp[1..h−1])`.
- Keep `dp` in a Fenwick tree for prefix maximums.

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why can a Fenwick tree do prefix *max* here, when range max is normally a segment-tree job?**
