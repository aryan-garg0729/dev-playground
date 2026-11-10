# Topological Sort

A **topological order** lists the nodes of a directed graph so every edge goes from an earlier node to a later one. It exists exactly when the graph has no directed cycle, that is, when it's a **DAG** (directed acyclic graph).

**Kahn's algorithm** (BFS on in-degrees):

```cpp
vector<int> topo(int n, const vector<vector<int>>& adj) {
    vector<int> indeg(n, 0), order;
    for (int u = 0; u < n; u++)
        for (int v : adj[u]) indeg[v]++;
    queue<int> q;
    for (int u = 0; u < n; u++)
        if (indeg[u] == 0) q.push(u);
    while (!q.empty()) {
        int u = q.front();
        q.pop();
        order.push_back(u);
        for (int v : adj[u])
            if (--indeg[v] == 0) q.push(v);
    }
    if ((int)order.size() < n) return {};   // empty: there's a cycle
    return order;
}
```

```python
from collections import deque


def topo(n, adj):
    indeg = [0] * n
    for u in range(n):
        for v in adj[u]:
            indeg[v] += 1
    q = deque(u for u in range(n) if indeg[u] == 0)
    order = []
    while q:
        u = q.popleft()
        order.append(u)
        for v in adj[u]:
            indeg[v] -= 1
            if indeg[v] == 0:
                q.append(v)
    return order if len(order) == n else None   # None: there's a cycle
```

Signals: prerequisites, build order, alien dictionary, course schedule.

**DP on a DAG**: process nodes in topological order, and every edge's source is final before you use it. Longest path, number of paths and earliest finish times all follow.

## Check yourself

Kahn's algorithm ends with `len(order) < n`. **What does that mean?**
