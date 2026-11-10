# Forest Queries II

**Practice:** [Forest Queries II on CSES](https://cses.fi/problemset/task/1739)

## The task

An n × n grid has trees and empty squares. Either flip one square's state, or count the trees inside a rectangle.

## Before you code

- A 2D Fenwick tree: the 1D loops, nested. O(log² n) per operation.
- A rectangle query uses the same four-term inclusion–exclusion as 2D prefix sums.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**What does `add(x, y, d)` look like in a 2D Fenwick tree?**
