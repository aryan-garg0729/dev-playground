# Maximum Subarray Sum

**Practice:** [Maximum Subarray Sum on CSES](https://cses.fi/problemset/task/1643)

## The task

Find the largest sum of a non-empty contiguous subarray.

## Pattern recap: Kadane's algorithm

For each position, keep the best sum of a subarray that **ends there**. Either extend the previous one or start fresh:

```cpp
long long best = a[0], cur = a[0];
for (size_t i = 1; i < a.size(); i++) {
    cur = max(a[i], cur + a[i]);      // a is a vector<long long>
    best = max(best, cur);
}
```

```python
best = cur = a[0]
for x in a[1:]:
    cur = max(x, cur + x)
    best = max(best, cur)
```

It's a one-variable DP. Variants: maximum product subarray (track both min and max), maximum circular subarray (total minus the minimum subarray).

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

A common version starts with `best = cur = 0` and uses `cur = max(0, cur + x)`. **What does it return for `[-3, -1, -2]`?**
