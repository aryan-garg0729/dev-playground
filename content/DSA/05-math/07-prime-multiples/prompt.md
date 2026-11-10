# Prime Multiples

**Practice:** [Prime Multiples on CSES](https://cses.fi/problemset/task/2185)

## The task

Given k distinct primes (k ≤ 20) and n up to 10¹⁸, count the numbers from 1 to n that are divisible by at least one of the primes.

## Pattern recap: inclusion–exclusion

|A ∪ B ∪ C| = |A| + |B| + |C| − |A ∩ B| − |A ∩ C| − |B ∩ C| + |A ∩ B ∩ C|.

Over every non-empty subset of the primes, add `n // product` for odd-sized subsets and subtract it for even-sized ones. That's 2^k subsets, about 10⁶ for k = 20.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Products of primes can exceed 10¹⁸. **What's the right fix?**
