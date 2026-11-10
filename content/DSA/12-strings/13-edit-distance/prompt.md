# Edit Distance

**Practice:** [Edit Distance on CSES](https://cses.fi/problemset/task/1639)

## The task

Find the edit distance between two strings of up to 5000 characters each.

## Before you code

- The "String DP" template with two rows.
- 2.5·10⁷ cells is fine in C++; in Python it's borderline, so keep the inner loop tight.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

With two rows, `prev` and `cur`, **which value does the "insert" move read?**
