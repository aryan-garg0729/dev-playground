# Repeating Substring

**Practice:** [Repeating Substring on CSES](https://cses.fi/problemset/task/2106)

## The task

Find the longest substring that appears at least twice, or print −1 if there's none.

## Before you code

- The largest LCP between neighbours in the suffix array.
- Or binary search on the length, with rolling hashes.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why can you binary search on the answer's length when using hashing?**
