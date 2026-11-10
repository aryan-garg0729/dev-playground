# Distinct Substrings

**Practice:** [Distinct Substrings on CSES](https://cses.fi/problemset/task/2105)

## The task

Count the distinct substrings of a string of up to 10⁵ characters.

## Before you code

- The formula from "Suffix Arrays". A suffix automaton works too: sum `len(state) − len(link(state))` over all states.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

For `"aaa"`, `sa = [2, 1, 0]` and the LCP values between neighbours are `[1, 2]`. **How many distinct substrings are there?**
