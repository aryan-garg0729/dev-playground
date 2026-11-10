# Distance Queries

**Practice:** [Distance Queries on CSES](https://cses.fi/problemset/task/1135)

## The task

Queries on a tree: the number of edges on the path between nodes a and b.

## Before you code

- `depth[a] + depth[b] − 2 · depth[lca(a, b)]`.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

a has depth 5, b has depth 3, and their LCA has depth 2. **What's the distance?**
