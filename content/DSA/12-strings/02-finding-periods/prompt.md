# Finding Periods

**Practice:** [Finding Periods on CSES](https://cses.fi/problemset/task/1733)

## The task

Find every period of a string: every length p such that the string is a prefix of its first p characters repeated over and over.

## Pattern recap: rolling hash

`h[i]` is the hash of the first i characters: `h[i+1] = (h[i] · B + code(s[i])) mod M`. Any substring's hash is `h[r] − h[l] · B^(r−l)` (mod M), in O(1) with precomputed powers. Compare substrings by comparing hashes. Use a random base and a large prime modulus (or two moduli) so collisions, and anti-hash tests, don't bite.

The Z-function or the prefix function solves this exactly, with no collision risk.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**p is a period exactly when…**
