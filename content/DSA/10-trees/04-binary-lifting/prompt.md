# Binary Lifting and LCA

"Go up k levels from node v" in O(log n): store the 2^j-th ancestor of every node.

```cpp
const int LOG = 18;                              // 2^18 > 2·10^5
vector<vector<int>> up(LOG, vector<int>(n + 1)); // up[0][v] = parent; the root's parent is 0
// fill: up[j][v] = up[j - 1][up[j - 1][v]]

int kth_ancestor(int v, int k) {
    for (int j = 0; j < LOG && k && v; j++, k >>= 1)
        if (k & 1) v = up[j][v];
    return k ? 0 : v;                            // 0 means "no such ancestor"
}
```

```python
LOG = 18                                  # 2^18 > 2·10^5
up = [[0] * (n + 1) for _ in range(LOG)]  # up[0][v] = parent; the root's parent is 0
# fill: up[j][v] = up[j - 1][up[j - 1][v]]


def kth_ancestor(v, k):
    j = 0
    while k and v:
        if k & 1:
            v = up[j][v]
        k >>= 1
        j += 1
    return v                              # 0 means "no such ancestor"
```

**Lowest common ancestor (LCA)**: lift the deeper node to the other's depth, then lift both together, from the biggest jump down, while their ancestors differ. One more step up is the LCA.

```cpp
int lca(int a, int b) {
    if (depth[a] < depth[b]) swap(a, b);
    a = kth_ancestor(a, depth[a] - depth[b]);
    if (a == b) return a;
    for (int j = LOG - 1; j >= 0; j--)
        if (up[j][a] != up[j][b]) {
            a = up[j][a];
            b = up[j][b];
        }
    return up[0][a];
}
```

```python
def lca(a, b):
    if depth[a] < depth[b]:
        a, b = b, a
    a = kth_ancestor(a, depth[a] - depth[b])
    if a == b:
        return a
    for j in reversed(range(LOG)):
        if up[j][a] != up[j][b]:
            a, b = up[j][a], up[j][b]
    return up[0][a]
```

**Distance** between a and b: `depth[a] + depth[b] - 2 * depth[lca(a, b)]`.

Binary lifting works on any structure where each node has one "next" node: trees (parent) and functional graphs (teleporters) alike.

## Check yourself

**Why does the LCA loop stop *below* the common ancestor and then take one more step?**
