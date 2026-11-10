# Range Updates and Sums

**Practice:** [Range Updates and Sums on CSES](https://cses.fi/problemset/task/1735)

**Read:** [Segment Tree on cp-algorithms](https://cp-algorithms.web.app/data_structures/segment_tree.html)

## The task

An array supports three operations: add x to every value in a range, set every value in a range to x, and report the sum of a range.

## Before you code

- A lazy segment tree with two kinds of pending tag: "set to v" and "add d".
- Decide how two tags combine before you write any code.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

A node already has a pending "add 5" and now receives "set to 2". **What's its pending tag afterwards?**
