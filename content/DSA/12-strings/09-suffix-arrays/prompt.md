# Suffix Arrays

A **suffix array** lists a string's suffixes in sorted order, stored as their start indices. The **LCP array** holds the longest common prefix of each pair of neighbours in that order.

- **Build** by prefix doubling: sort by the first 2^k characters using the ranks for 2^(k−1). That's O(n log² n) with comparison sorting, or O(n log n) with radix sort. **Kasai's algorithm** then builds the LCP array in O(n).
- **Distinct substrings**: every substring is a prefix of some suffix, and neighbours share `lcp[i]` of their prefixes, so the count is `Σ (n − sa[i]) − Σ lcp[i]`.
- **Longest repeated substring**: the largest value in the LCP array.

Your sheet marks this optional, for a good reason: it rarely comes up in interviews. It's still a great test of writing careful O(n log n) code.

## Check yourself

**Why is the longest repeated substring equal to the largest LCP value?**
