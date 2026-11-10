# Cycle Detection

Two different problems hide under "find a cycle".

**Undirected graph.** In a DFS, reaching an already visited node that isn't the one you came from means a cycle. Track the parent (by edge id if there can be parallel edges).

**Directed graph.** "Visited" isn't enough: reaching a node that finished in another branch isn't a cycle. Use three colours:

- white: not visited yet
- grey: on the current DFS path
- black: finished

An edge into a **grey** node is a back edge, and that's a cycle.

```cpp
vector<int> color;   // 0 white, 1 grey, 2 black

// Recursive is fine in C++ for typical limits; very deep graphs may still need an explicit stack.
bool dfs(int u, const vector<vector<int>>& adj) {
    color[u] = 1;
    for (int v : adj[u]) {
        if (color[v] == 1) return true;              // back edge: a cycle
        if (color[v] == 0 && dfs(v, adj)) return true;
    }
    color[u] = 2;
    return false;
}

bool has_cycle(int n, const vector<vector<int>>& adj) {
    color.assign(n, 0);
    for (int s = 0; s < n; s++)
        if (color[s] == 0 && dfs(s, adj)) return true;
    return false;
}
```

```python
def has_cycle(n, adj):
    color = [0] * n                      # 0 white, 1 grey, 2 black
    for s in range(n):
        if color[s]:
            continue
        color[s] = 1
        stack = [(s, iter(adj[s]))]
        while stack:
            u, it = stack[-1]
            v = next(it, None)
            if v is None:
                color[u] = 2
                stack.pop()
            elif color[v] == 1:
                return True
            elif color[v] == 0:
                color[v] = 1
                stack.append((v, iter(adj[v])))
    return False
```

To print the cycle, keep a parent array and walk back from u to v.

Other tools: Kahn's algorithm also detects directed cycles (nodes left over), and a DSU detects undirected cycles as edges are added.

## Check yourself

In a directed graph, the DFS meets an edge u → v where v is **black** (finished). **What does that tell you?**
