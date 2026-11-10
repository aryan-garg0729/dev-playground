# Maximum Xor Subarray

**Practice:** [Maximum Xor Subarray on CSES](https://cses.fi/problemset/task/1655)

## The task

Find the largest XOR of a contiguous subarray.

## Before you code

- Turn subarrays into pairs of prefix XORs (`p[r] ^ p[l]`), then maximize over pairs with a bitwise trie.
- Insert `p[0] = 0` before anything else.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why must the trie start with the prefix value 0 inserted?**
