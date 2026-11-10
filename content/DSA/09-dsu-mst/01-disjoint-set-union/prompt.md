# Disjoint Set Union

A **disjoint set union** (DSU, or union–find) keeps a collection of groups and answers "are a and b in the same group?" and "merge the groups of a and b" in almost O(1) each.

```cpp
struct DSU {
    vector<int> parent, size;
    DSU(int n) : parent(n), size(n, 1) { iota(parent.begin(), parent.end(), 0); }

    int find(int x) {
        while (parent[x] != x) {
            parent[x] = parent[parent[x]];          // path halving
            x = parent[x];
        }
        return x;
    }

    bool unite(int a, int b) {                      // `union` is a keyword in C++
        a = find(a), b = find(b);
        if (a == b) return false;
        if (size[a] < size[b]) swap(a, b);
        parent[b] = a;                              // small tree under big tree
        size[a] += size[b];
        return true;
    }
};
```

```python
class DSU:
    def __init__(self, n):
        self.parent = list(range(n))
        self.size = [1] * n

    def find(self, x):
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]   # path halving
            x = self.parent[x]
        return x

    def union(self, a, b):
        a, b = self.find(a), self.find(b)
        if a == b:
            return False
        if self.size[a] < self.size[b]:
            a, b = b, a
        self.parent[b] = a              # small tree under big tree
        self.size[a] += self.size[b]
        return True
```

With union by size and path compression, each operation costs O(α(n)), which is at most 4 for any realistic n.

**Kruskal's minimum spanning tree (MST)**: sort the edges by weight, and take an edge whenever `union` returns True (it joins two different components). The taken edges form an MST.

Uses: components as edges arrive, redundant connection, accounts merge, grouping equivalent items, Kruskal.

A **rollback DSU** (union by size, *no* path compression) can undo its last union. It's the base of offline dynamic connectivity.

Reading: [DSU on cp-algorithms](https://cp-algorithms.web.app/data_structures/disjoint_set_union.html).

## Check yourself

Kruskal's algorithm runs out of edges having taken fewer than n − 1. **What does that mean?**
