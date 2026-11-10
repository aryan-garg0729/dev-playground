# E - Knapsack 2

**Practice:** [E - Knapsack 2 on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_e)

## The task

0/1 knapsack with a huge capacity and small values: the case from "Designing a DP State" (state swapping: least weight for a fixed value).

## Before you code

- One array indexed by value. For each item, loop v **downwards** so each item is used at most once.

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

In the 1D loop `for v from V down to value[i]: dp[v] = min(dp[v], dp[v - value[i]] + weight[i])`, **why go downwards?**
