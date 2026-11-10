# DFS and BFS

Most graph interview problems are a DFS or BFS on an adjacency list or a grid.

```cpp
vector<int> bfs(int start, const vector<vector<int>>& adj) {
    vector<int> dist(adj.size(), -1);
    dist[start] = 0;
    queue<int> q;
    q.push(start);
    while (!q.empty()) {
        int u = q.front();
        q.pop();
        for (int v : adj[u])
            if (dist[v] == -1) {
                dist[v] = dist[u] + 1;
                q.push(v);
            }
    }
    return dist;
}
```

```python
from collections import deque


def bfs(start, adj):
    dist = {start: 0}
    q = deque([start])
    while q:
        u = q.popleft()
        for v in adj[u]:
            if v not in dist:
                dist[v] = dist[u] + 1
                q.append(v)
    return dist
```

- **BFS** gives shortest paths when every edge costs the same (unweighted graphs, grids).
- **DFS** explores everything reachable: connected components, flood fill, cycle detection, topological order.
- **Multi-source BFS**: put *every* source in the queue at distance 0. One O(V + E) pass gives each cell its distance to the nearest source (rotting oranges, walls and gates, nearest exit).
- **0-1 BFS**: edges cost 0 or 1. Use a deque and push 0-cost edges to the front.

Grid neighbours: `for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1))`.

Python note: recursive DFS hits the recursion limit at about 1,000 levels, and a 1000 × 1000 grid can need 10⁶. Use an explicit stack, or BFS.

## Check yourself

**When should BFS mark a node as visited?**
