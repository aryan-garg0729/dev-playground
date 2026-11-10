# Tree Diameter

The **diameter** of a tree is its longest path between any two nodes, counted in edges.

**Two BFS passes.** BFS from any node; the farthest node a is one end of a diameter. BFS from a; the farthest distance is the diameter.

It works because in a tree, the node farthest from *any* start is always the end of some diameter.

**DP alternative**: at each node, combine its two deepest child branches: `best = max(best, d1 + d2)`. It's also O(n), and it extends to weighted trees and constrained longest paths.

Facts you'll reuse:

- Every node's farthest node is at distance `max(dist_a[v], dist_b[v])`, where a and b are the two ends of a diameter.
- A tree with n nodes has n − 1 edges, so BFS is O(n).
- In Python, always traverse trees iteratively: a path-shaped tree can be 2·10⁵ levels deep.

## Check yourself

**Why doesn't the two-BFS trick find the longest simple path in a general graph?**
