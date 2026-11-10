# Floyd–Warshall

All-pairs shortest paths in O(n³), with three loops in the right order:

```cpp
// dist[i][j]: edge weight, INF if none, 0 on the diagonal.
// INF = 1e18 keeps INF + INF inside a long long.
void floyd(vector<vector<long long>>& dist, int n) {
    for (int k = 0; k < n; k++)
        for (int i = 0; i < n; i++)
            for (int j = 0; j < n; j++)
                dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j]);
}
```

```python
def floyd(dist, n):        # dist[i][j]: edge weight, INF if none, 0 on the diagonal
    for k in range(n):
        dk = dist[k]
        for i in range(n):
            di, dik = dist[i], dist[i][k]
            for j in range(n):
                if dik + dk[j] < di[j]:
                    di[j] = dik + dk[j]
```

After step k, `dist[i][j]` is the shortest path that uses only nodes 0..k as stopovers. It's a DP over "which middle nodes are allowed".

Use it when n is up to about 500 (in C++; Python struggles beyond about 150), when there are many pair queries, or when there are negative edges but no negative cycle. With several edges between two nodes, keep the cheapest.

Its cousin for "visit every node once" is bitmask DP: `dp[mask][v]` = best way to visit exactly the nodes in `mask`, ending at v.

## Check yourself

**What breaks if the `k` loop is innermost instead of outermost?**
