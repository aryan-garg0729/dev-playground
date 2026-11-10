# Distinct Values Subarrays II

**Practice:** [Distinct Values Subarrays II on CSES](https://cses.fi/problemset/task/2428)

## The task

Count the subarrays that contain at most k distinct values.

## Before you code

- Two pointers with a count map, like the window template.
- For each right end r, which starts give a valid subarray?

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

`[l, r]` is the longest valid window ending at r. **How many valid subarrays end at r?**
