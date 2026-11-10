# Road Construction

**Practice:** [Road Construction on CSES](https://cses.fi/problemset/task/1676)

## The task

Roads are built one at a time. After each road, print the number of components and the size of the largest one.

## Before you code

- A DSU with sizes. Keep `components` and `largest` as running values.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why is a running maximum enough for "largest component", with no heap or multiset?**
