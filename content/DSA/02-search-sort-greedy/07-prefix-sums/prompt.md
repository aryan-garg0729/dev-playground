# Prefix Sums

Need many range sums? Precompute once:

```cpp
vector<long long> pre(a.size() + 1, 0);
for (size_t i = 0; i < a.size(); i++) pre[i + 1] = pre[i] + a[i];
// sum of a[l..r-1] (l inclusive, r exclusive)
long long total = pre[r] - pre[l];
```

```python
pre = [0]
for x in a:
    pre.append(pre[-1] + x)
# sum of a[l:r] (l inclusive, r exclusive)
total = pre[r] - pre[l]
```

## Counting subarrays with a given sum

`a[l:r]` sums to x exactly when `pre[r] - pre[l] == x`. Walk r from left to right and count the earlier prefixes equal to `pre[r] - x`:

```cpp
long long count_sum(const vector<long long>& a, long long x) {
    map<long long, long long> seen;   // map, not unordered_map: safe from anti-hash tests
    seen[0] = 1;
    long long cur = 0, ans = 0;
    for (long long v : a) {
        cur += v;
        auto it = seen.find(cur - x);
        if (it != seen.end()) ans += it->second;
        seen[cur]++;
    }
    return ans;
}
```

```python
from collections import defaultdict


def count_sum(a, x):
    seen = defaultdict(int)
    seen[0] = 1
    cur = ans = 0
    for v in a:
        cur += v
        ans += seen[cur - x]
        seen[cur] += 1
    return ans
```

O(n), and it works with negative numbers. Variants: "subarray sum divisible by k" (store `cur % k`), "longest subarray with sum x" (store the first index of each prefix), "equal number of 0s and 1s" (map 0 to −1 first).

## Two pointers or a hash map?

A sliding window only works when every value is non-negative, because then growing the window can only raise the sum. With negatives, use the hash map.

## In 2D

`pre[i][j]` is the sum of the top-left i × j rectangle. Any rectangle is four lookups: `pre[r2][c2] - pre[r1][c2] - pre[r2][c1] + pre[r1][c1]`.

## Check yourself

The array has negative numbers. You need to count the subarrays that sum to x in O(n) or O(n log n). **What works?**
