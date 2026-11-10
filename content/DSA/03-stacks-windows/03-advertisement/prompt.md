# Advertisement

**Practice:** [Advertisement on CSES](https://cses.fi/problemset/task/1142)

## The task

A fence is made of boards with given heights. Find the largest rectangle that fits inside it: the "largest rectangle in a histogram" problem (also LeetCode 84).

## Before you code

- For each board, the widest rectangle of exactly its height reaches from just after its previous smaller board to just before its next smaller board.
- Two monotonic-stack passes (or one careful pass) give both neighbours.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Heights are `[2, 1, 5, 6, 2, 3]`. **What's the largest area?**
