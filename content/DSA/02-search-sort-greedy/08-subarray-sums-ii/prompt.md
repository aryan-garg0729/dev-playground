# Subarray Sums II

**Practice:** [Subarray Sums II on CSES](https://cses.fi/problemset/task/1661)

## The task

Count the subarrays whose sum is exactly x. Values can be negative.

## Before you code

- This is the template from "Prefix Sums".
- In C++, `unordered_map` with the default hash can be slowed to a crawl by crafted tests (anti-hash tests). Use `map`, or a randomized custom hash.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

In the template, **what does `seen[0] = 1` account for?**
