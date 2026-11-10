# Path Queries

**Practice:** [Path Queries on CSES](https://cses.fi/problemset/task/1138)

## The task

Either change a node's value, or report the sum of values on the path from the root to a node.

## Before you code

- A value at u counts for every node in u's subtree. Turn the point update into a range update on `[tin[u], tout[u]]` and read a single point.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Node u's value increases by d. **Which Fenwick-on-differences update keeps every root-path sum right?**
