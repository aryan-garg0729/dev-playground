# Dynamic Connectivity

**Practice:** [Dynamic Connectivity on CSES](https://cses.fi/problemset/task/2133)

## The task

Roads are both added and removed over time. After each event, print the number of components. (Your sheet marks it "very hard": it's a full offline technique.)

## Before you code

- Work offline. Each road is alive during some time intervals. Put each interval into a segment tree over time; it lands in O(log q) nodes.
- DFS the segment tree. Entering a node, union its roads with a rollback DSU; at a leaf, record the component count; leaving the node, roll back.
- Total: O(q log q log n).

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why must this DSU skip path compression?**
