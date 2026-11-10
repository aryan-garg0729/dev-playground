# Company Queries I

**Practice:** [Company Queries I on CSES](https://cses.fi/problemset/task/1687)

## The task

A company's hierarchy is a tree rooted at the general director (employee 1). Queries: who is employee x's boss k levels up? Print −1 if there's no such boss.

## Before you code

- `kth_ancestor` from "Binary Lifting and LCA", with 0 meaning "past the top".

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why use 0 as the root's parent instead of the root itself?**
