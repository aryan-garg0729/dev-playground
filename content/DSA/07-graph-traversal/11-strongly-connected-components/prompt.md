# Strongly Connected Components

In a directed graph, a **strongly connected component** (SCC) is a largest group of nodes where every node can reach every other. Collapse each SCC into one node and you get a DAG (the **condensation**), so problems on directed graphs with cycles become DAG problems.

**Kosaraju's algorithm** makes two passes:

1. DFS the graph and record the nodes in the order they finish.
2. Reverse every edge. Take the nodes in *decreasing* finish time; each DFS on the reversed graph from an unvisited node marks exactly one SCC.

**Tarjan's algorithm** does it in one DFS using low-link values. Both are O(V + E). Write them iteratively in Python.

Uses: "can every city reach every other?" (one SCC), 2-SAT, collapsing cycles before a DAG DP (most coins collected on a path).

## Check yourself

**The whole graph is strongly connected exactly when…**
