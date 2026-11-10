# Static Range Minimum Queries

**Practice:** [Static Range Minimum Queries on CSES](https://cses.fi/problemset/task/1647)

**Read:** [Sparse Table on cp-algorithms](https://cp-algorithms.web.app/data_structures/sparse-table.html)

## The task

An array never changes. Answer many queries "minimum of `a[l..r]`".

## Pattern recap: sparse table

`st[j][i]` is the minimum of the 2^j elements starting at i. Build it in O(n log n), then answer any query with two overlapping blocks in O(1), because min doesn't mind counting an element twice.

```cpp
vector<vector<int>> build(const vector<int>& a) {
    int n = a.size();
    vector<vector<int>> st = {a};
    for (int j = 1; (1 << j) <= n; j++) {
        int half = 1 << (j - 1);
        vector<int> row(n - (1 << j) + 1);
        for (int i = 0; i < (int)row.size(); i++)
            row[i] = min(st[j - 1][i], st[j - 1][i + half]);
        st.push_back(row);
    }
    return st;
}

int query(const vector<vector<int>>& st, int l, int r) {   // inclusive l..r
    int j = 31 - __builtin_clz(r - l + 1);   // floor(log2(length))
    return min(st[j][l], st[j][r - (1 << j) + 1]);
}
```

```python
def build(a):
    st = [a[:]]
    j = 1
    while (1 << j) <= len(a):
        prev, half = st[-1], 1 << (j - 1)
        st.append([min(prev[i], prev[i + half]) for i in range(len(a) - (1 << j) + 1)])
        j += 1
    return st


def query(st, l, r):            # inclusive l..r
    j = (r - l + 1).bit_length() - 1
    return min(st[j][l], st[j][r - (1 << j) + 1])
```

Also used for range max, range gcd, and LCA through an Euler tour.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why can't a sparse table answer range *sum* queries in O(1) the same way?**
