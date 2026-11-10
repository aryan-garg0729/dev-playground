# G - Longest Path

**Practice:** [G - Longest Path on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_g)

## The task

Find the longest path, counted in edges, in a directed acyclic graph: DP on a DAG.

## Before you code

- Topological order (Kahn), then `dp[v] = max(dp[u] + 1)` over edges u → v.
- Or a memoized DFS on "longest path starting at v".

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why doesn't this approach work on a graph with cycles?**
