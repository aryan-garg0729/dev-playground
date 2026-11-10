# Bellman–Ford

Bellman–Ford handles negative edge weights and detects negative cycles, in O(V · E).

```cpp
struct Edge { int u, v; long long w; };

vector<long long> bellman_ford(int n, const vector<Edge>& edges, int src) {
    const long long INF = LLONG_MAX / 4;
    vector<long long> dist(n, INF);
    dist[src] = 0;
    for (int round = 0; round < n - 1; round++) {
        bool changed = false;
        for (const auto& e : edges)
            if (dist[e.u] != INF && dist[e.u] + e.w < dist[e.v]) {
                dist[e.v] = dist[e.u] + e.w;
                changed = true;
            }
        if (!changed) break;
    }
    return dist;
}
```

```python
def bellman_ford(n, edges, src):
    INF = float("inf")
    dist = [INF] * n
    dist[src] = 0
    for _ in range(n - 1):
        changed = False
        for u, v, w in edges:
            if dist[u] != INF and dist[u] + w < dist[v]:
                dist[v] = dist[u] + w
                changed = True
        if not changed:
            break
    return dist
```

A shortest path has at most n − 1 edges, so n − 1 rounds are enough. If an **n-th** round still improves something, there's a negative cycle reachable from the source.

To print a negative cycle: remember a node x relaxed in round n, along with the parent pointers. Walk the parents n times from x (that puts you on the cycle), then keep following parents until you come back to the same node.

To find **any** negative cycle, not just one reachable from a source, start every `dist` at 0, as if a virtual source had an edge to every node.

Longest path = shortest path with negated weights. Then "unbounded" means a negative cycle that lies on some path to the target.

## Check yourself

**Why walk the parent pointers n times from x before collecting the cycle?**
