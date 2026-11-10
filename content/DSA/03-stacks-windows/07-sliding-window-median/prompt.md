# Sliding Window Median

**Practice:** [Sliding Window Median on CSES](https://cses.fi/problemset/task/1076)

## The task

Print the median of each window of k elements. For even k, use the lower of the two middle values.

## Before you code

- Two heaps with lazy deletion, or an ordered multiset (C++: keep an iterator at the median as elements come and go).
- Keep the sizes balanced: the max-heap `low` holds ⌈k/2⌉ live elements, and its top is the answer.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

With two heaps (max-heap `low`, min-heap `high`), the element leaving the window is marked for lazy deletion. **What must you update right away, even though the element is still physically inside a heap?**
