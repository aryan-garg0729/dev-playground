# KMP and the Prefix Function

The **prefix function** `pi[i]` is the length of the longest proper prefix of `s[:i+1]` that is also a suffix of it (a **border**).

```cpp
vector<int> prefix_function(const string& s) {
    vector<int> pi(s.size(), 0);
    for (size_t i = 1; i < s.size(); i++) {
        int k = pi[i - 1];
        while (k && s[i] != s[k]) k = pi[k - 1];
        if (s[i] == s[k]) k++;
        pi[i] = k;
    }
    return pi;
}
```

```python
def prefix_function(s):
    pi = [0] * len(s)
    for i in range(1, len(s)):
        k = pi[i - 1]
        while k and s[i] != s[k]:
            k = pi[k - 1]
        if s[i] == s[k]:
            k += 1
        pi[i] = k
    return pi
```

It's O(n) by the same amortized argument as the monotonic stack: k rises by at most 1 per step, so it can't fall more than n times in total.

**Pattern matching**: compute it for `pattern + "#" + text`, with a separator that appears in neither. Every position where `pi == len(pattern)` ends a match.

**All borders** of s: `pi[n−1]`, then `pi[pi[n−1] − 1]`, and so on down to 0.

The **Z-function** (`z[i]` = the longest common prefix of s and `s[i:]`) solves the same problems. Learn whichever one you can write from memory.

## Check yourself

**What is the prefix function of `"aabaaab"`?**
