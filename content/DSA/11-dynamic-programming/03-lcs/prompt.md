# F - LCS

**Practice:** [F - LCS on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_f)

## The task

Print one longest common subsequence of strings s and t (each up to 3000 characters): string DP plus solution reconstruction.

## Before you code

- Fill `dp[i][j]` = LCS length of `s[:i]` and `t[:j]`, then walk back from `dp[n][m]` to rebuild the string.

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Walking back at (i, j), you find `s[i-1] == t[j-1]`. **What do you do?**
