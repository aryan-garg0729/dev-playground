# Hotel Queries

**Practice:** [Hotel Queries on CSES](https://cses.fi/problemset/task/1143)

## The task

Hotels stand in a row, each with some free rooms. Groups arrive in order; each needs r rooms in a single hotel and takes the **first** (leftmost) hotel with enough free rooms. Print that hotel, or 0 if none fits.

## Before you code

- A max segment tree over the free rooms.
- Descend from the root: go left whenever the left child's maximum is big enough.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why descend the tree instead of binary searching on the index with a range-max query at each step?**
