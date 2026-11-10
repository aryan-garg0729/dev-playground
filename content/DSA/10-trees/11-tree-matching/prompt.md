# Tree Matching

**Practice:** [Tree Matching on CSES](https://cses.fi/problemset/task/1130)

## The task

Choose as many edges as possible so that no two chosen edges share a node: a maximum matching in a tree.

## Before you code

- DP: `dp[v][0]` = best in v's subtree with v unmatched; `dp[v][1]` = best with v matched to one of its children.
- Or a greedy rule. The question is about which one.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Which greedy rule gives a maximum matching in a tree?**
