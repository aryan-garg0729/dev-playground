# Z - Frog 3

**Practice:** [Z - Frog 3 on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_z)

## The task

A frog jumps forward along stones of strictly increasing height. A jump from i to j costs (h_j − h_i)² + C. Minimize the total cost of reaching the last stone, with n up to 2·10⁵: the convex hull trick.

## Before you code

- `dp[j] = min over i < j of (dp[i] + (h_j − h_i)²) + C` is O(n²).
- Expand the square: for a fixed i, it's a straight line in h_j. The **convex hull trick** keeps those lines and finds the lowest one at h_j quickly.

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Expanding `dp[i] + (h_j − h_i)²`, **which line does stone i contribute, as a function of x = h_j?**
