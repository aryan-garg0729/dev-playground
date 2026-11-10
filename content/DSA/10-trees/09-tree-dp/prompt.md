# Tree DP and Rerooting

**Tree DP**: root the tree, then compute each node's answer from its children's answers, bottom-up. Processing nodes in reverse BFS order puts every child before its parent, with no recursion:

```cpp
vector<int> order = {root}, parent(n + 1, 0);
parent[root] = -1;
for (size_t i = 0; i < order.size(); i++) {         // the vector grows as we go: BFS order
    int u = order[i];
    for (int v : adj[u])
        if (v != parent[u]) {
            parent[v] = u;
            order.push_back(v);
        }
}
for (int i = (int)order.size() - 1; i >= 0; i--) {  // children before parents
    int u = order[i];
    // ...
}
```

```python
order, parent = [root], [0] * (n + 1)
parent[root] = -1
for u in order:                  # the list grows while we loop: BFS order
    for v in adj[u]:
        if v != parent[u]:
            parent[v] = u
            order.append(v)
for u in reversed(order):        # children before parents
    ...
```

Typical states: `size[v]`, `dp[v][taken]` (matching, independent set, house robber III), `down[v]` = the deepest path below v.

**Rerooting**: you want each node's answer as if *it* were the root. One pass bottom-up gives the answers inside each subtree; one pass top-down pushes the "rest of the tree" contribution from parent to child. Two O(n) passes instead of n separate DPs.

Example: the sum of distances from v to every node. Bottom-up, compute `size[v]` and `sub[v]` (the sum of distances from v into its subtree). Then `ans[root] = sub[root]`, and moving the root from u to its child v brings `size[v]` nodes one step closer and pushes the other `n − size[v]` one step farther:

```text
ans[v] = ans[u] - size[v] + (n - size[v])
```

## Check yourself

**Why does the bottom-up pass process `reversed(order)`?**
