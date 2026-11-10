# Palindrome Queries

**Practice:** [Palindrome Queries on CSES](https://cses.fi/problemset/task/2420)

## The task

A string supports two operations: change one character, or check whether `s[a..b]` is a palindrome.

## Before you code

- Two Fenwick trees of hashes, one over s and one over its reverse. A substring is a palindrome when its forward and backward hashes agree.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Position k stores `code(s[k]) · B^k` in the forward tree. **How do you compare the forward hash of `s[a..b]` with the backward one?**
