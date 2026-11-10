# Euler Tour Flattening

Number the nodes in DFS order, and record when each node is entered (`tin`) and left (`tout`). Then **v's subtree is the contiguous range `[tin[v], tout[v]]`**, and subtree problems become array range problems.

```cpp
vector<int> tin, tout;   // sized n + 1
int timer = 0;

void dfs(int v, int p) {
    tin[v] = timer++;
    for (int c : adj[v])
        if (c != p) dfs(c, v);
    tout[v] = timer - 1;
}
```

```python
tin, tout, timer = [0] * (n + 1), [0] * (n + 1), 0
stack = [(root, 0, False)]
while stack:
    v, p, done = stack.pop()
    if done:
        tout[v] = timer - 1
        continue
    tin[v] = timer
    timer += 1
    stack.append((v, p, True))
    for c in adj[v]:
        if c != p:
            stack.append((c, v, False))
```

- **Subtree sum with point updates**: a Fenwick tree indexed by `tin`.
- **Root-to-node path sum with point updates**: a value at u counts for every node in u's subtree, so add it as a range update on `[tin[u], tout[u]]` and read the single point `tin[v]` (Fenwick on differences).
- **Is u an ancestor of v?** `tin[u] <= tin[v] <= tout[u]`.

## Check yourself

**How do you check in O(1) whether u is an ancestor of v?**
