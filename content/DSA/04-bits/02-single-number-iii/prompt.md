# Single Number III

**Practice:** [Single Number III on LeetCode](https://leetcode.com/problems/single-number-iii/description/)

## The task

Every number appears exactly twice except two, which appear once. Find those two in O(n) time and O(1) extra space.

## Before you code

- XOR everything first. What do you know about the result?
- How could one of its bits split the numbers into two groups?

## Check your approach

Solve it on LeetCode first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

After XORing everything you have `s = a ^ b` with `a ≠ b`. **Why is `s & -s` useful?**
