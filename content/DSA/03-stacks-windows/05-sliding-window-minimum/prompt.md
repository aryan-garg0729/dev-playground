# Sliding Window Minimum

**Practice:** [Sliding Window Minimum on CSES](https://cses.fi/problemset/task/3221)

## The task

Find the minimum of every window of k consecutive elements. The input can be very large, so read the statement for how it's given and what exactly to print.

## Before you code

- The monotonic deque from "Sliding Window".
- A heap with lazy deletion also works, but it's O(n log n). On a big input you want the O(n) deque and fast input reading.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

In `window_mins`, **why does the deque pop from the back while `a[dq[-1]] >= x`, before appending x?**
