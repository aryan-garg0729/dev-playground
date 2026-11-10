# Dijkstra's Algorithm

For non-negative edge weights, Dijkstra's algorithm finds shortest paths from one source in O((V + E) log V).

```cpp
vector<long long> dijkstra(int src, const vector<vector<pair<int, int>>>& adj) {
    const long long INF = LLONG_MAX / 4;
    vector<long long> dist(adj.size(), INF);
    priority_queue<pair<long long, int>, vector<pair<long long, int>>, greater<>> pq;
    dist[src] = 0;
    pq.push({0, src});
    while (!pq.empty()) {
        auto [d, u] = pq.top();
        pq.pop();
        if (d > dist[u]) continue;                   // stale entry, already improved
        for (auto [v, w] : adj[u])
            if (d + w < dist[v]) {
                dist[v] = d + w;
                pq.push({dist[v], v});
            }
    }
    return dist;
}
```

```python
import heapq


def dijkstra(src, adj, n):
    INF = float("inf")
    dist = [INF] * n
    dist[src] = 0
    pq = [(0, src)]
    while pq:
        d, u = heapq.heappop(pq)
        if d > dist[u]:
            continue                      # stale entry, already improved
        for v, w in adj[u]:
            if d + w < dist[v]:
                dist[v] = d + w
                heapq.heappush(pq, (dist[v], v))
    return dist
```

**State graphs.** When the cost depends on more than the node (a coupon used or not, fuel left, stops so far), make that extra information part of the node: `dist[node][used]`. Then run the same algorithm on the bigger graph.

**k shortest routes**: let each node be popped up to k times instead of once. The j-th time node n pops, you have the j-th shortest route.

Negative edges break Dijkstra: a node popped as "final" could still improve later through a negative edge. Use Bellman–Ford for those.

## Check yourself

**Why does the template skip an entry when `d > dist[u]`?**
