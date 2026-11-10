# Longest Palindrome

**Practice:** [Longest Palindrome on CSES](https://cses.fi/problemset/task/1111)

## The task

Print the longest palindromic substring of a string of up to 10⁶ characters.

## Before you code

- Expanding around every centre is O(n²) in the worst case (think `"aaaa…"`), so you need Manacher here.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

In the transformed string t, the best radius is `p[i] = 5` at `i = 9`. **Where does the palindrome start in s?**
