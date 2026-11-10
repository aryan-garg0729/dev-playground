# S - Digit Sum

**Practice:** [S - Digit Sum on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_s)

## The task

Count the integers from 1 to K (K can have 10⁴ digits) whose digit sum is a multiple of D, modulo 10⁹+7: digit DP.

## Before you code

- Walk K's digits from the left with the state (position, digit sum mod D, tight). `tight` means the digits so far equal K's prefix, so the next digit can't exceed K's digit.

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

The DP as described counts every number from 0 to K. **What's the final adjustment?**
