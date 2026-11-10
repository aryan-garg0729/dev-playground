# Distinct Colors

**Practice:** [Distinct Colors on CSES](https://cses.fi/problemset/task/1139)

## The task

Every node has a colour. For every node, count the distinct colours in its subtree.

## Pattern recap: small-to-large merging

Give each node a set of colours, and merge each child's set into its parent's by always pouring the **smaller** set into the larger one. A colour only moves when the set it's in at least doubles, so each colour moves O(log n) times: O(n log n) set operations in total.

Alternative: Euler tour, then offline "distinct values in a range" queries with a Fenwick tree.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

In small-to-large, the parent's set is smaller than one child's. **What should you do?**
