# Factory Machines

**Practice:** [Factory Machines on CSES](https://cses.fi/problemset/task/1620)

## The task

n machines work at the same time; machine i needs `k[i]` seconds per product. Find the least time to make t products in total.

## Pattern recap: binary search on the answer

Use it when you can't compute the answer directly but you **can** check a guess quickly, and the check is monotone: once a guess works, every bigger guess works too.

```cpp
// smallest x in [lo, hi] with ok(x), assuming ok(hi) is true
long long first_true(long long lo, long long hi, function<bool(long long)> ok) {
    while (lo < hi) {
        long long mid = lo + (hi - lo) / 2;
        if (ok(mid)) hi = mid;
        else lo = mid + 1;
    }
    return lo;
}
```

```python
def first_true(lo, hi, ok):
    # smallest x in [lo, hi] with ok(x), assuming ok(hi) is True
    while lo < hi:
        mid = (lo + hi) // 2
        if ok(mid):
            hi = mid
        else:
            lo = mid + 1
    return lo
```

Signals: "minimum time / maximum size such that…", "minimize the largest…" (split array largest sum, ship packages within D days, Koko eating bananas).

## Before you code

- Write `ok(T)`: how many products do the machines finish in T seconds?
- That count can be huge. With fixed-width integers, stop adding once it reaches t.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Which `hi` is guaranteed to satisfy `ok(hi)` for every valid input?**
