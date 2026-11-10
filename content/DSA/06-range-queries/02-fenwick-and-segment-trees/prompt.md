# Fenwick and Segment Trees

When the array **changes** between queries, prefix sums and sparse tables break. Two structures give O(log n) updates and queries.

**Fenwick tree (BIT)**: short and fast; handles prefix sums (and other invertible operations).

```cpp
struct BIT {
    vector<long long> t;                   // 1-indexed
    BIT(int n) : t(n + 1, 0) {}

    void add(int i, long long d) {
        for (; i < (int)t.size(); i += i & -i) t[i] += d;
    }

    long long prefix(int i) {              // sum of positions 1..i
        long long s = 0;
        for (; i > 0; i -= i & -i) s += t[i];
        return s;
    }
};
```

```python
class BIT:
    def __init__(self, n):
        self.t = [0] * (n + 1)          # 1-indexed

    def add(self, i, d):
        while i < len(self.t):
            self.t[i] += d
            i += i & -i

    def prefix(self, i):                # sum of positions 1..i
        s = 0
        while i > 0:
            s += self.t[i]
            i -= i & -i
        return s
```

A range sum is `prefix(r) - prefix(l - 1)`. For **range update, point query**, store differences: `add(l, d)` and `add(r + 1, -d)`, then the value at i is `prefix(i)`.

**Segment tree**: each node stores the answer for its segment. It handles any associative operation (min, max, sum, gcd), finds "the first index with value ≥ x" by descending from the root, and supports **range updates with lazy propagation**: a pending update stored at a node and pushed to its children only when needed.

| Need | Use |
| --- | --- |
| point update, range sum | Fenwick |
| range add, point query | Fenwick on differences |
| range add, range sum | two Fenwicks, or a lazy segment tree |
| point update, range min/max | segment tree |
| first position with value ≥ x | segment tree descent |
| 2D grid, point update, rectangle sum | 2D Fenwick |

Reading: [segment trees](https://cp-algorithms.web.app/data_structures/segment_tree.html) and [Fenwick trees](https://cp-algorithms.web.app/data_structures/fenwick.html) on cp-algorithms.

## Check yourself

**Why is a Fenwick tree a poor fit for "point update, range *minimum*"?**
