from h import P, X

CSES = "https://cses.fi/problemset/task/"
CPA = "https://cp-algorithms.web.app/"
ATC = "https://atcoder.jp/contests/dp/tasks/dp_"

graphs = dict(slug="graph-traversal", id="dsa.graphs", title="Graph Traversal", lessons=[
P("dfs-and-bfs", "DFS and BFS", 1, ["graphs", "bfs", "dfs"], r"""
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
""", "**When should BFS mark a node as visited?**", [
    "When it's pushed onto the queue",
    "When it's popped from the queue",
    "Either; the result and the cost are the same",
    "Only the start node needs marking",
], "Mark on push. If you wait until pop, the same node can be pushed many times by different neighbours before it's processed. The distances stay right, but the queue can grow far beyond V, which turns into time or memory failures on dense graphs and big grids. Dijkstra is the opposite: there you check on pop, because a better distance can still arrive later."),

X("counting-rooms", "Counting Rooms", CSES + "1192", 1, ["graphs", "grid"], r"""
## The task

A map is a grid of floor (`.`) and wall (`#`) squares. Count the rooms: connected areas of floor, moving up, down, left and right.

## Before you code

- Loop over every cell. Each unvisited floor cell starts a new flood fill and a new room.
- Write the flood fill iteratively, for the reason in "DFS and BFS".
""", "The map can be 1000 × 1000. **In Python, what's the main risk of a recursive flood fill?**", [
    "A RecursionError, or a crash, when one room has about 10⁶ cells",
    "It counts the rooms wrong",
    "It visits some cells twice",
    "No risk; recursion is fine for grids",
], "A snake-shaped room can make the recursion as deep as the room is big, up to 10⁶ levels. Python's default limit is 1,000, and raising it with `sys.setrecursionlimit` can overflow the real C stack and crash the interpreter. An explicit stack (or BFS) is O(rows × cols) with no depth problem. C++ copes with deeper recursion, but 10⁶ levels can still overflow its stack."),

X("monsters", "Monsters", CSES + "1194", 4, ["graphs", "bfs", "multi-source"], r"""
## The task

A grid has walls, monsters and you. Every step, you and every monster move one square. Reach the border of the grid without a monster ever being able to catch you. Print a path, or say it's impossible.

## Before you code

- First, a multi-source BFS from **all** monsters gives each cell the earliest time any monster can be there.
- Then a BFS from your position, only stepping onto cells you reach strictly before every monster.
- Keep parent pointers so you can print the path.
""", "**Why is it enough to plan against the monsters' *earliest* arrival times, instead of simulating how they move?**", [
    "If a monster can reach a cell no later than you, it can wait there, so the cell is unsafe; otherwise it can never catch you on it",
    "Monsters always move towards you",
    "Monsters can't pass through each other",
    "Because a BFS from your position is enough on its own",
], "Assume the monsters play perfectly. Any cell a monster can reach at a time ≤ your arrival time is lost, because the monster can get there first and wait. A cell you reach strictly earlier is safe when you stand on it. That turns a chase into two BFS passes, O(rows × cols). Starting the monster BFS from every monster at once is what keeps it linear."),

P("topological-sort", "Topological Sort", 2, ["graphs", "topological-sort"], r"""
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
""", "Kahn's algorithm ends with `len(order) < n`. **What does that mean?**", [
    "The graph has a directed cycle",
    "The graph is disconnected",
    "Some nodes have no edges",
    "The queue was processed in the wrong order",
], "Nodes on a cycle (and everything downstream of one) never reach in-degree 0, because each waits for another node on the cycle, so they never enter the queue. A disconnected DAG still sorts completely: each part has its own in-degree-0 nodes, and isolated nodes have in-degree 0 from the start."),

X("course-schedule", "Course Schedule", CSES + "1679", 2, ["graphs", "topological-sort"], r"""
## The task

n courses have prerequisites ("course a before course b"). Print any valid order to take them, or say it's impossible.

## Before you code

- Kahn's algorithm, as written. The impossible case is the cycle check.
""", "**Which edge should \"course a must come before course b\" become?**", [
    "a → b, so a course's in-degree counts its prerequisites",
    "b → a",
    "Either direction gives a valid answer",
    "An undirected edge between a and b",
], "Edges point from what comes first to what depends on it, so a course with in-degree 0 has no pending prerequisites and can be taken now. Reversing every edge produces the reverse order, which lists each course *before* its prerequisites. You could fix that by printing the result backwards, but mixing up the direction is the classic bug."),

X("longest-flight-route", "Longest Flight Route", CSES + "1680", 3, ["graphs", "dag-dp"], r"""
## The task

In a directed flight network with no cycles, find the route from city 1 to city n that visits the most cities, and print it.

## Before you code

- DP in topological order: `best[v] = max(best[u] + 1)` over flights u → v, with routes starting only at city 1.
- Store each city's parent to rebuild the route.
""", "Some cities can't be reached from city 1. **How should `best` start?**", [
    "best[1] = 1 and every other city unreachable (−∞), so routes can't start anywhere else",
    "Every city starts at 1",
    "Every city starts at 0",
    "It doesn't matter in a DAG",
], "If every city started at 1, a long chain that never touches city 1 could look like the best way into n, giving a route that doesn't start at 1. Mark unreachable cities and never extend a route out of them. Print 'IMPOSSIBLE' if n stays unreachable. Running a BFS from city 1 first and keeping only the cities it reaches works too."),

P("cycle-detection", "Cycle Detection", 2, ["graphs", "dfs"], r"""
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
""", "In a directed graph, the DFS meets an edge u → v where v is **black** (finished). **What does that tell you?**", [
    "Nothing about a cycle: v's subtree is done and can't lead back to u",
    "There's a cycle through u and v",
    "The graph isn't a DAG",
    "u is unreachable from the start",
], "A black node has explored everything reachable from it without meeting u on the current path, so this edge (a cross or forward edge) closes no cycle. Only an edge to a grey node, which is still on the path, closes one. Treating every visited node as a cycle is the classic bug when undirected code is reused on a directed graph."),

X("round-trip", "Round Trip", CSES + "1669", 3, ["graphs", "dfs"], r"""
## The task

In an undirected road network, find any round trip (a cycle through at least 3 cities) and print it, or say none exists.

## Before you code

- DFS with parent tracking. When you see a visited neighbour that isn't your parent, walk the parent pointers back to it.
""", "**Why does \"visited and not my parent\" guarantee a cycle of at least 3 cities here?**", [
    "With no duplicate roads, the only way back to a visited neighbour without reusing the same road goes through at least one other city",
    "Because the graph is connected",
    "Because DFS visits cities in increasing order",
    "It doesn't; it can report a cycle of 2 cities",
], "The parent check stops you from treating the road you just used as a cycle. Without duplicate roads, any other visited neighbour was reached by a different route, so the cycle through it has at least 3 cities. With parallel edges, track the edge id instead of the parent node. To print it, walk parent pointers from the current city back to that neighbour, then close the loop."),

X("round-trip-ii", "Round Trip II", CSES + "1678", 3, ["graphs", "dfs"], r"""
## The task

Same as Round Trip, but the flights are one-way: find any directed cycle.

## Before you code

- Three colours, as in "Cycle Detection". When you reach a grey node, the cycle is on your DFS path.
""", "The DFS path is `1 → 4 → 2 → 5`, and city 5 has a flight to city 4 (grey). **Which cycle do you print?**", [
    "4 2 5 4", "1 4 2 5 1", "5 4", "2 5 2",
], "The cycle starts at the grey node you reached (4), follows the path down to the current node (5), then takes the flight back: 4 → 2 → 5 → 4. City 1 is on the path but not on the cycle, and including it would need a flight 5 → 1 that doesn't exist. Collect the nodes by walking parents from 5 back to 4, reverse them, and repeat the first at the end."),

X("planets-cycles", "Planets Cycles", CSES + "1751", 4, ["graphs", "functional-graph"], r"""
## The task

Every planet has exactly one teleporter to some planet. For every starting planet, count the teleports you make until you reach a planet you've already visited.

## Pattern recap: functional graphs

A graph where every node has exactly one outgoing edge is a **functional graph**. It's a set of cycles with trees hanging off them, all pointing in towards the cycles. For each node, the answer is "steps to reach its cycle + the cycle's length".

Find each cycle by walking until you meet a node from the current walk, then fill in the answers backwards along the walk.
""", "A planet is 3 teleports away from a cycle of length 5. **What's its answer?**", [
    "8", "3", "5", "7",
], "You make 3 teleports to reach the cycle, then 5 more take you around it and back to the first cycle planet, which you've seen: 3 + 5 = 8. Every planet on the cycle has the answer 5. Compute it with an iterative walk that stops either at a planet that's already solved (add its answer) or at one on the current walk (you've found a new cycle). Each planet is processed once: O(n)."),

P("strongly-connected-components", "Strongly Connected Components", 3, ["graphs", "scc"], r"""
In a directed graph, a **strongly connected component** (SCC) is a largest group of nodes where every node can reach every other. Collapse each SCC into one node and you get a DAG (the **condensation**), so problems on directed graphs with cycles become DAG problems.

**Kosaraju's algorithm** makes two passes:

1. DFS the graph and record the nodes in the order they finish.
2. Reverse every edge. Take the nodes in *decreasing* finish time; each DFS on the reversed graph from an unvisited node marks exactly one SCC.

**Tarjan's algorithm** does it in one DFS using low-link values. Both are O(V + E). Write them iteratively in Python.

Uses: "can every city reach every other?" (one SCC), 2-SAT, collapsing cycles before a DAG DP (most coins collected on a path).
""", "**The whole graph is strongly connected exactly when…**", [
    "Node 1 reaches every node, and node 1 also reaches every node in the reversed graph",
    "Node 1 reaches every node",
    "Every node has an outgoing edge",
    "The graph has a cycle",
], "Reaching every node from node 1 shows 1 → x for each x. Reaching every node from 1 in the reversed graph shows x → 1 in the original. Chain the two through node 1 and every node reaches every other. One direction isn't enough: a star with edges going out of node 1 passes the first check, but nothing can come back. That's two BFS runs, and it's the whole next problem."),

X("flight-routes-check", "Flight Routes Check", CSES + "1682", 3, ["graphs", "scc"], r"""
## The task

Check whether every city can reach every other city by one-way flights. If not, print a pair (a, b) with no route from a to b.

## Before you code

- The two-BFS check from "Strongly Connected Components". The BFS that fails tells you the pair.
""", "The BFS from city 1 on the **reversed** graph doesn't reach city 4. **Which pair do you print?**", [
    "4 1", "1 4", "4 4", "Either order is fine",
], "In the reversed graph, an edge x → y means y → x in the original. Not reaching 4 from 1 in the reversed graph means there's no original route from 4 to 1, so print `4 1`. If the *forward* BFS from 1 misses some city x instead, print `1 x`. Printing the pair backwards is the usual wrong answer."),

X("planets-and-kingdoms", "Planets and Kingdoms", CSES + "1683", 3, ["graphs", "scc"], r"""
## The task

Group planets into kingdoms: two planets share a kingdom when each can reach the other through teleporters. Print the number of kingdoms and a kingdom label for every planet.

## Before you code

- The kingdoms are exactly the SCCs. Use Kosaraju, written iteratively.
""", "In Kosaraju's second pass, **why take the nodes in decreasing finish time?**", [
    "The last node to finish is in an SCC nothing else points into; with edges reversed nothing leaves it, so its DFS can't leak into other SCCs",
    "To visit the nodes in numeric order",
    "So the second pass does less work",
    "Any order works in the second pass",
], "The node with the latest finish time lies in a 'source' SCC of the condensation: no other SCC has an edge into it. After reversing the edges nothing leaves that SCC, so a DFS from it marks exactly that SCC. Remove it and repeat, peeling off one SCC at a time. In an arbitrary order, a DFS on the reversed graph can run from one SCC into another and merge them."),
])

paths = dict(slug="shortest-paths", id="dsa.paths", title="Shortest Paths", lessons=[
P("dijkstra", "Dijkstra's Algorithm", 2, ["graphs", "dijkstra"], r"""
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
""", "**Why does the template skip an entry when `d > dist[u]`?**", [
    "It's an outdated entry: u was pushed again with a smaller distance, and that one was already processed",
    "To skip nodes with negative distances",
    "To stop the algorithm early",
    "Without it, Dijkstra returns wrong distances",
], "Python's heap has no 'decrease key', so improving a distance pushes a new entry and leaves the old one behind. When the old one pops, `d > dist[u]` shows it's stale. Without the check the distances are still right, but you re-scan u's edges for every stale entry, which can be very slow on bad inputs."),

X("shortest-routes-i", "Shortest Routes I", CSES + "1671", 2, ["graphs", "dijkstra"], r"""
## The task

In a directed graph with weighted flights, print the shortest distance from city 1 to every city.

## Before you code

- The template, with fast input reading.
- Distances can reach about 10¹⁴. In C++ or Java use 64-bit integers, and make "infinity" bigger than any real path.
""", "A C++ solution uses `int dist[]` and `INF = 1e9`. Weights go up to 10⁹, with up to 10⁵ cities. **What goes wrong?**", [
    "Real distances exceed both int and INF, so they overflow or look unreachable",
    "Nothing; 1e9 is large enough",
    "Dijkstra becomes O(n²)",
    "Only the output format breaks",
], "A path of 10⁵ flights of weight 10⁹ costs 10¹⁴. That overflows a 32-bit int, and even a correct distance above 10⁹ would compare as worse than 'infinity'. Use `long long` and an INF like `LLONG_MAX / 4` (or `1e18`), so `dist + w` can't overflow either."),

X("flight-discount", "Flight Discount", CSES + "1195", 3, ["graphs", "dijkstra", "state-graph"], r"""
## The task

Find the cheapest route from city 1 to city n when you may halve the price of one flight on the way (rounded down).

## Before you code

- State graph: `dist[city][used]` with used ∈ {0, 1}. From state 0 you can fly at full price (stay at 0) or half price (move to 1).
""", "An alternative without a state graph: run Dijkstra from city 1, and from city n on the reversed graph. **The answer is then the minimum, over every flight u → v with price w, of…**", [
    "d1[u] + w // 2 + dn[v]",
    "d1[u] + w + dn[v]",
    "d1[v] + w // 2 + dn[u]",
    "(d1[u] + w + dn[v]) // 2",
], "Use the coupon on flight u → v: reach u as cheaply as possible (d1[u]), pay half for that flight, then get from v to n as cheaply as possible (dn[v], from the reversed graph). Trying every flight is O(E). The third option mixes up the flight's direction; the last halves the whole trip instead of one flight. Both methods are O((V + E) log V), and the state graph is the one that extends to k coupons."),

X("flight-routes", "Flight Routes", CSES + "1196", 4, ["graphs", "dijkstra"], r"""
## The task

Print the prices of the k cheapest routes from city 1 to city n (k ≤ 10; routes may revisit cities).

## Before you code

- Dijkstra where each city may be popped up to k times. Each pop of city n is the next route price.
""", "**Why is it correct to stop expanding a city once it has been popped k times?**", [
    "A route through its (k+1)-th cheapest arrival is beaten by the k routes through its k cheaper arrivals, so it can't be among the k cheapest",
    "A route can visit each city at most k times",
    "Dijkstra guarantees each city has at most k edges",
    "It isn't correct; it only works when k = 1",
], "Suppose a route to n used the (k+1)-th cheapest way of reaching city c. Swapping in each of the k cheaper ways of reaching c gives k different routes that are all cheaper, so it isn't among the k cheapest. A pop counter per city replaces the single `dist` array, and the total is O(k · E log(k · E))."),

P("floyd-warshall", "Floyd–Warshall", 2, ["graphs", "all-pairs"], r"""
All-pairs shortest paths in O(n³), with three loops in the right order:

```cpp
// dist[i][j]: edge weight, INF if none, 0 on the diagonal.
// INF = 1e18 keeps INF + INF inside a long long.
void floyd(vector<vector<long long>>& dist, int n) {
    for (int k = 0; k < n; k++)
        for (int i = 0; i < n; i++)
            for (int j = 0; j < n; j++)
                dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j]);
}
```

```python
def floyd(dist, n):        # dist[i][j]: edge weight, INF if none, 0 on the diagonal
    for k in range(n):
        dk = dist[k]
        for i in range(n):
            di, dik = dist[i], dist[i][k]
            for j in range(n):
                if dik + dk[j] < di[j]:
                    di[j] = dik + dk[j]
```

After step k, `dist[i][j]` is the shortest path that uses only nodes 0..k as stopovers. It's a DP over "which middle nodes are allowed".

Use it when n is up to about 500 (in C++; Python struggles beyond about 150), when there are many pair queries, or when there are negative edges but no negative cycle. With several edges between two nodes, keep the cheapest.

Its cousin for "visit every node once" is bitmask DP: `dp[mask][v]` = best way to visit exactly the nodes in `mask`, ending at v.
""", "**What breaks if the `k` loop is innermost instead of outermost?**", [
    "The result is wrong: dist[i][k] and dist[k][j] may not be final yet when they're used",
    "Nothing; the loop order doesn't matter",
    "It becomes O(n⁴)",
    "It only breaks with negative edges",
], "The DP relies on step k using distances that already allow every middle node before k. With k innermost, `dist[i][k]` is used before it has been improved through other nodes, so some shortest paths are never found. Loop order is the most common Floyd–Warshall bug, and it passes small tests surprisingly often."),

X("shortest-routes-ii", "Shortest Routes II", CSES + "1672", 2, ["graphs", "all-pairs"], r"""
## The task

Up to 500 cities with two-way roads of given lengths, and many queries: the shortest distance from a to b, or −1 if there's no route.

## Before you code

- Floyd–Warshall once; then each query is a lookup.
- At n = 500 that's 1.25·10⁸ inner steps, too slow for plain Python. Use C++ for this one.
""", "The input has two roads between cities 2 and 3, of lengths 7 and 4. **What should `dist[2][3]` start at?**", [
    "4", "7", "11", "Whichever road comes last in the input",
], "The starting matrix holds the best *direct* road, so take the minimum over all roads between the same pair: `dist[a][b] = min(dist[a][b], c)`, and the same for `dist[b][a]`, since roads go both ways. Overwriting with the last road read is a bug that only fails on tests with duplicate roads."),

X("hamiltonian-flights", "Hamiltonian Flights", CSES + "1690", 4, ["graphs", "bitmask-dp"], r"""
## The task

Count the routes from city 1 to city n that visit every city exactly once (n ≤ 20), modulo 10⁹+7.

## Before you code

- `dp[mask][v]` = number of ways to visit exactly the cities in `mask`, starting at city 1 and ending at v.
- That's 2²⁰ · 20 ≈ 2·10⁷ states, so C++ is a must, and memory layout matters.
""", "**Which pruning is valid and roughly halves the work?**", [
    "Skip masks that don't contain city 1, and never extend a route out of city n before the mask is full",
    "Skip masks with an even number of cities",
    "Process the masks in decreasing order",
    "Leave city n out of the DP entirely",
], "Every route starts at city 1, so masks without it have zero ways; skipping them halves the states. City n must come last, so never extend a route out of n before every city is visited. Masks are processed in increasing order, because adding a city makes the mask bigger. The transition is `dp[mask | 1<<u][u] += dp[mask][v]` for each flight v → u with u not in mask."),

P("bellman-ford", "Bellman–Ford", 3, ["graphs", "negative-cycles"], r"""
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
""", "**Why walk the parent pointers n times from x before collecting the cycle?**", [
    "x might be on a path leading into the cycle; after n steps back you're guaranteed to be on the cycle",
    "To compute the cycle's length",
    "To check that the cycle is negative",
    "It isn't needed; x is always on the cycle",
], "The node relaxed in round n is affected by a negative cycle, but it can sit on a tail hanging off the cycle. A tail has fewer than n nodes, so n steps back can't still be on it, and you end on the cycle. From there, follow parents until you return to the same node, then reverse and print."),

X("cycle-finding", "Cycle Finding", CSES + "1197", 3, ["graphs", "negative-cycles"], r"""
## The task

Find any negative cycle in a directed graph and print it, or print NO.

## Before you code

- Bellman–Ford with every distance starting at 0, n rounds, then the parent walk from "Bellman–Ford".
""", "**Why start every `dist` at 0 instead of `dist[1] = 0` and the rest at infinity?**", [
    "A negative cycle that node 1 can't reach would be missed; zeros everywhere act like a source with an edge to every node",
    "To avoid overflow",
    "Because negative cycles always pass through node 0",
    "It's only faster; the answer is the same either way",
], "The problem asks for any negative cycle in the whole graph. With a single source, relaxations never reach a cycle in another part of the graph, so you'd print NO by mistake. Starting everything at 0 is the same as adding a virtual node with a 0-weight edge to every node, and it lets every cycle show up."),

X("high-score", "High Score", CSES + "1673", 4, ["graphs", "negative-cycles"], r"""
## The task

Find the maximum score of a path from room 1 to room n, where each tunnel adds its score (possibly negative). If the score can be made arbitrarily large, print −1.

## Before you code

- Negate the scores and run Bellman–Ford from room 1.
- Not every negative cycle matters: only one that's reachable from room 1 **and** can reach room n.
""", "In round n, room 7 improves again, but room n can't be reached from room 7. **What does that mean for the answer?**", [
    "Nothing: that cycle can't be used on a route that ends at n, so ignore it",
    "Print −1",
    "The answer is the score at room 7",
    "Run Bellman–Ford again from room 7",
], "An unbounded score needs a cycle you can enter from room 1 and leave towards room n. Run a reachability pass forward from room 1 and backward from room n (on the reversed graph), and print −1 only when a room still improving in the extra round passes both checks. Forgetting the 'can reach n' part is the classic wrong answer on this problem."),
])

dsu = dict(slug="dsu-mst", id="dsa.dsu", title="Disjoint Set Union & MST", lessons=[
P("disjoint-set-union", "Disjoint Set Union", 2, ["dsu", "mst"], r"""
A **disjoint set union** (DSU, or union–find) keeps a collection of groups and answers "are a and b in the same group?" and "merge the groups of a and b" in almost O(1) each.

```cpp
struct DSU {
    vector<int> parent, size;
    DSU(int n) : parent(n), size(n, 1) { iota(parent.begin(), parent.end(), 0); }

    int find(int x) {
        while (parent[x] != x) {
            parent[x] = parent[parent[x]];          // path halving
            x = parent[x];
        }
        return x;
    }

    bool unite(int a, int b) {                      // `union` is a keyword in C++
        a = find(a), b = find(b);
        if (a == b) return false;
        if (size[a] < size[b]) swap(a, b);
        parent[b] = a;                              // small tree under big tree
        size[a] += size[b];
        return true;
    }
};
```

```python
class DSU:
    def __init__(self, n):
        self.parent = list(range(n))
        self.size = [1] * n

    def find(self, x):
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]   # path halving
            x = self.parent[x]
        return x

    def union(self, a, b):
        a, b = self.find(a), self.find(b)
        if a == b:
            return False
        if self.size[a] < self.size[b]:
            a, b = b, a
        self.parent[b] = a              # small tree under big tree
        self.size[a] += self.size[b]
        return True
```

With union by size and path compression, each operation costs O(α(n)), which is at most 4 for any realistic n.

**Kruskal's minimum spanning tree (MST)**: sort the edges by weight, and take an edge whenever `union` returns True (it joins two different components). The taken edges form an MST.

Uses: components as edges arrive, redundant connection, accounts merge, grouping equivalent items, Kruskal.

A **rollback DSU** (union by size, *no* path compression) can undo its last union. It's the base of offline dynamic connectivity.

Reading: [DSU on cp-algorithms](https://cp-algorithms.web.app/data_structures/disjoint_set_union.html).
""", "Kruskal's algorithm runs out of edges having taken fewer than n − 1. **What does that mean?**", [
    "The graph is disconnected, so there's no spanning tree",
    "The MST was finished early",
    "The graph had a cycle",
    "Some edges had equal weights",
], "A spanning tree on n nodes always has exactly n − 1 edges. If Kruskal runs out of edges before taking that many, some components could never be joined: the graph is disconnected (CSES wants 'IMPOSSIBLE'). Cycles are normal; those edges are just skipped because `union` returns False. Equal weights can give several MSTs, all with the same total."),

X("road-reparation", "Road Reparation", CSES + "1675", 2, ["dsu", "mst"], r"""
## The task

Find the minimum total cost of repairing roads so that every city can reach every other, or say it's impossible.

## Before you code

- Kruskal from "Disjoint Set Union".
- The total can reach 2·10⁵ · 10⁹, so use 64-bit integers.
""", "Prim's algorithm also solves this. **When is Prim a better choice than Kruskal?**", [
    "On dense graphs: Prim's O(V²) array version beats sorting about V² edges",
    "Always; Kruskal is slower",
    "When the weights are negative",
    "When the graph is disconnected",
], "Kruskal is O(E log E), dominated by the sort, and the DSU work is almost free. Prim with a heap is O(E log V), about the same. Prim's real advantage is a dense graph: the O(V²) version with a plain array beats sorting V² edges. Both handle negative weights, since an MST only compares weights, and neither builds a spanning tree of a disconnected graph."),

X("road-construction", "Road Construction", CSES + "1676", 2, ["dsu"], r"""
## The task

Roads are built one at a time. After each road, print the number of components and the size of the largest one.

## Before you code

- A DSU with sizes. Keep `components` and `largest` as running values.
""", "**Why is a running maximum enough for \"largest component\", with no heap or multiset?**", [
    "Components only merge, so sizes only grow: the new largest is max(largest, size of the merged component)",
    "Because the largest component never changes",
    "Because a DSU keeps its components sorted",
    "It isn't enough; you need a multiset of sizes",
], "Building a road can only join two components into a bigger one; nothing ever shrinks. So after each successful union, `largest = max(largest, size[root])` and `components -= 1`. A road between cities that are already connected changes nothing. Removing roads would break this, which is the next problem."),

X("dynamic-connectivity", "Dynamic Connectivity", CSES + "2133", 5, ["dsu", "offline", "segment-tree"], r"""
## The task

Roads are both added and removed over time. After each event, print the number of components. (Your sheet marks it "very hard": it's a full offline technique.)

## Before you code

- Work offline. Each road is alive during some time intervals. Put each interval into a segment tree over time; it lands in O(log q) nodes.
- DFS the segment tree. Entering a node, union its roads with a rollback DSU; at a leaf, record the component count; leaving the node, roll back.
- Total: O(q log q log n).
""", "**Why must this DSU skip path compression?**", [
    "Path compression changes many parent pointers per find, so a union can't be undone by restoring one pointer",
    "Path compression makes find slower",
    "Path compression breaks union by size",
    "It doesn't need to skip it",
], "Rollback works by remembering, for each union, the one parent pointer and the one size it changed, and restoring them in reverse order. Path compression rewrites pointers all over the place during `find`, so there's no small record to undo. Union by size alone keeps trees O(log n) deep, which is good enough here."),
])

trees = dict(slug="trees", id="dsa.trees", title="Trees", lessons=[
P("tree-diameter-pattern", "Tree Diameter", 2, ["trees", "bfs"], r"""
The **diameter** of a tree is its longest path between any two nodes, counted in edges.

**Two BFS passes.** BFS from any node; the farthest node a is one end of a diameter. BFS from a; the farthest distance is the diameter.

It works because in a tree, the node farthest from *any* start is always the end of some diameter.

**DP alternative**: at each node, combine its two deepest child branches: `best = max(best, d1 + d2)`. It's also O(n), and it extends to weighted trees and constrained longest paths.

Facts you'll reuse:

- Every node's farthest node is at distance `max(dist_a[v], dist_b[v])`, where a and b are the two ends of a diameter.
- A tree with n nodes has n − 1 edges, so BFS is O(n).
- In Python, always traverse trees iteratively: a path-shaped tree can be 2·10⁵ levels deep.
""", "**Why doesn't the two-BFS trick find the longest simple path in a general graph?**", [
    "With cycles, the farthest node from a start needn't be an end of the longest path; that problem is NP-hard in general",
    "BFS doesn't work on graphs",
    "General graphs have no longest path",
    "It does; it works on every graph",
], "The proof relies on paths in a tree being unique: there's only one path between any two nodes. Cycles create alternatives and the argument falls apart. In fact, the longest simple path in a general graph is NP-hard. In trees it's linear, which is one reason trees get their own chapter."),

X("tree-diameter", "Tree Diameter", CSES + "1131", 1, ["trees", "bfs"], r"""
## The task

Find the diameter of a tree with n nodes.

## Before you code

- Two BFS passes, iterative.
""", "A path 1 – 2 – 3 – 4 – 5 has an extra node 6 attached to node 3. **What's the diameter?**", [
    "4", "3", "5", "6",
], "The longest path is 1 – 2 – 3 – 4 – 5, with 4 edges. Paths through 6 (like 1 – 2 – 3 – 6) have only 3. Diameters are counted in edges here; always check whether a problem wants edges or nodes."),

X("tree-distances-i", "Tree Distances I", CSES + "1132", 2, ["trees", "bfs"], r"""
## The task

For every node, print the distance to the node farthest from it.

## Before you code

- Find both diameter ends a and b, then use the fact from "Tree Diameter".
""", "**How many BFS passes does the diameter-ends method need in total?**", [
    "3", "2", "n", "1",
], "One BFS from any node finds a, one from a finds b (and gives `dist_a`), and one from b gives `dist_b`. The answer for v is `max(dist_a[v], dist_b[v])`. Three O(n) passes, instead of one BFS from every node at O(n²). Rerooting DP, later in this chapter, gets the same answer another way and generalizes further."),

P("binary-lifting", "Binary Lifting and LCA", 3, ["trees", "binary-lifting", "lca"], r"""
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
""", "**Why does the LCA loop stop *below* the common ancestor and then take one more step?**", [
    "Jumping only while the ancestors differ is the safe test: equal ancestors could be the LCA or anything above it",
    "To save one jump",
    "Because up[j][a] is undefined at the root",
    "It doesn't matter; jumping while they're equal works too",
], "If `up[j][a] == up[j][b]`, that node is *a* common ancestor, but maybe not the lowest one, so you can't jump there. If they differ, the LCA is strictly above both, so the jump is safe. After trying every j from large to small, a and b are the LCA's children on each side, and `up[0][a]` is the answer."),

X("company-queries-i", "Company Queries I", CSES + "1687", 2, ["trees", "binary-lifting"], r"""
## The task

A company's hierarchy is a tree rooted at the general director (employee 1). Queries: who is employee x's boss k levels up? Print −1 if there's no such boss.

## Before you code

- `kth_ancestor` from "Binary Lifting and LCA", with 0 meaning "past the top".
""", "**Why use 0 as the root's parent instead of the root itself?**", [
    "If the root pointed to itself, going too far would quietly return the root instead of signalling 'no ancestor'",
    "Employee 0 is the real director",
    "It makes the table smaller",
    "Python lists can't store self-references",
], "With `up[0][root] = root`, asking for 10 levels above an employee who is 3 levels deep returns the root, and you'd print 1 instead of −1. A sentinel 0 whose ancestors are all 0 makes 'past the top' visible: print −1 when the result is 0. Comparing k with `depth[x]` first works too."),

X("planets-queries-i", "Planets Queries I", CSES + "1750", 3, ["binary-lifting", "functional-graph"], r"""
## The task

Each planet teleports to one planet (a functional graph). Queries: where do you end up after k teleports from planet x, with k up to 10⁹?

## Before you code

- Binary lifting works on any functional graph, not only trees. With k ≤ 10⁹ you need 30 levels.
""", "With n = 2·10⁵ planets and 30 levels, **about how much memory does the jump table take as 32-bit integers?**", [
    "About 24 MB", "About 240 KB", "About 2.4 GB", "About 6 MB",
], "30 · 2·10⁵ = 6·10⁶ entries × 4 bytes ≈ 24 MB, which is fine in C++. In Python, a list of lists costs far more (an 8-byte pointer per entry plus the int objects), so use the `array` module, or answer the queries offline one bit level at a time. Memory is a real constraint for binary lifting, and interviewers like to ask about it."),

X("company-queries-ii", "Company Queries II", CSES + "1688", 3, ["trees", "lca"], r"""
## The task

The same company tree. Queries: the lowest common boss of employees a and b.

## Before you code

- The `lca` function. Compute depths with an iterative BFS from the root.
""", "a is an ancestor of b. **What does `lca` return, and where?**", [
    "a, right after lifting b to a's depth, because they're now equal",
    "a's parent",
    "The root",
    "b",
], "After the swap, the deeper node b is lifted `depth[b] − depth[a]` levels, and because a is its ancestor, it lands exactly on a. The `if a == b: return a` check catches that before the main loop. Forgetting that check returns a's parent instead, which is the most common LCA bug."),

X("distance-queries", "Distance Queries", CSES + "1135", 3, ["trees", "lca"], r"""
## The task

Queries on a tree: the number of edges on the path between nodes a and b.

## Before you code

- `depth[a] + depth[b] − 2 · depth[lca(a, b)]`.
""", "a has depth 5, b has depth 3, and their LCA has depth 2. **What's the distance?**", [
    "4", "6", "8", "2",
], "a climbs 5 − 2 = 3 edges to the LCA and b climbs 3 − 2 = 1, so the distance is 3 + 1 = 4 = 5 + 3 − 2 · 2. Subtracting the LCA's depth only once (6) counts the path from the root to the LCA once too often."),

P("tree-dp", "Tree DP and Rerooting", 3, ["trees", "dp"], r"""
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
""", "**Why does the bottom-up pass process `reversed(order)`?**", [
    "BFS order lists every parent before its children, so the reverse handles every child before its parent",
    "It's the order the nodes appear in the input",
    "It sorts the nodes by value",
    "Any order works for tree DP",
], "A node's DP needs its children's final values. In BFS order every node comes after its parent, so in the reversed list every node comes after all of its children. This replaces recursion, which matters in Python (the recursion limit) and in C++ on path-shaped trees with 2·10⁵ levels."),

X("tree-distances-ii", "Tree Distances II", CSES + "1133", 3, ["trees", "rerooting"], r"""
## The task

For every node, print the sum of its distances to all other nodes.

## Before you code

- The rerooting formula from "Tree DP and Rerooting".
""", "A tree has n = 5 nodes. Node u has `ans[u] = 8`, and its child v has a subtree of 2 nodes. **What's `ans[v]`?**", [
    "9", "7", "6", "8",
], "Moving from u to v brings the 2 nodes of v's subtree one step closer (−2) and pushes the other 5 − 2 = 3 nodes one step farther (+3): 8 − 2 + 3 = 9. The sums reach about n² ≈ 4·10¹⁰, so use 64-bit."),

X("tree-matching", "Tree Matching", CSES + "1130", 3, ["trees", "dp", "greedy"], r"""
## The task

Choose as many edges as possible so that no two chosen edges share a node: a maximum matching in a tree.

## Before you code

- DP: `dp[v][0]` = best in v's subtree with v unmatched; `dp[v][1]` = best with v matched to one of its children.
- Or a greedy rule. The question is about which one.
""", "**Which greedy rule gives a maximum matching in a tree?**", [
    "Process nodes bottom-up; match a node with its parent if both are still free",
    "Match the highest-degree nodes first",
    "Match edges in input order whenever both ends are free",
    "Match the root first, then work downwards",
], "A leaf has only one possible partner, its parent, and matching them never hurts (an exchange argument). Going bottom-up applies that rule to every node whose children are already decided. Input-order greedy only finds a *maximal* matching (one you can't extend), which can be smaller: on the path 1 – 2 – 3 – 4, taking edge 2 – 3 first leaves 1 match instead of 2. The DP gives the same answer and extends to weighted versions."),

P("euler-tour", "Euler Tour Flattening", 3, ["trees", "euler-tour"], r"""
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
""", "**How do you check in O(1) whether u is an ancestor of v?**", [
    "tin[u] <= tin[v] and tout[v] <= tout[u]",
    "tin[u] < tin[v]",
    "depth[u] < depth[v]",
    "tout[u] < tin[v]",
], "v's interval sits inside u's exactly when v is in u's subtree. Comparing only entry times or depths can't tell a cousin from a descendant. `tout[u] < tin[v]` means u's subtree ended before v started: they're in different branches. This O(1) test is handy inside LCA code and tree queries."),

X("subtree-queries", "Subtree Queries", CSES + "1137", 3, ["trees", "euler-tour", "fenwick"], r"""
## The task

A rooted tree has a value on each node. Either change a node's value, or report the sum of values in a node's subtree.

## Before you code

- Euler tour, plus a Fenwick tree over the `tin` positions.
""", "Node v's value changes from 4 to 9. **What's the Fenwick update?**", [
    "add(tin[v], 5)",
    "add(tin[v], 9)",
    "add 5 to the whole range [tin[v], tout[v]]",
    "add(v, 5)",
], "A Fenwick tree stores sums, so record the change, 9 − 4 = 5, at v's position `tin[v]`, and keep the current values in an array to compute the next change. Adding 9 would count the old value too. Indexing by v instead of `tin[v]` loses the 'subtree is a contiguous range' property. The query is `prefix(tout[v]) − prefix(tin[v] − 1)`, shifted to 1-based positions."),

X("path-queries", "Path Queries", CSES + "1138", 4, ["trees", "euler-tour", "fenwick"], r"""
## The task

Either change a node's value, or report the sum of values on the path from the root to a node.

## Before you code

- A value at u counts for every node in u's subtree. Turn the point update into a range update on `[tin[u], tout[u]]` and read a single point.
""", "Node u's value increases by d. **Which Fenwick-on-differences update keeps every root-path sum right?**", [
    "add(tin[u], d) and add(tout[u] + 1, -d)",
    "add(tin[u], d) only",
    "add(tout[u], d) and add(tin[u], -d)",
    "add(d) at every ancestor of u",
], "Root-to-v sums change for exactly the nodes in u's subtree, which is a contiguous range of positions. On a difference array a range add is two point updates: at the start, and just past the end. A point query (the prefix sum at `tin[v]`) then gives the full path sum. Both operations are O(log n)."),

X("distinct-colors", "Distinct Colors", CSES + "1139", 4, ["trees", "small-to-large"], r"""
## The task

Every node has a colour. For every node, count the distinct colours in its subtree.

## Pattern recap: small-to-large merging

Give each node a set of colours, and merge each child's set into its parent's by always pouring the **smaller** set into the larger one. A colour only moves when the set it's in at least doubles, so each colour moves O(log n) times: O(n log n) set operations in total.

Alternative: Euler tour, then offline "distinct values in a range" queries with a Fenwick tree.
""", "In small-to-large, the parent's set is smaller than one child's. **What should you do?**", [
    "Swap them, so the parent keeps the child's big set, then pour the small one in",
    "Copy the child's set into the parent's anyway",
    "Skip that child",
    "Sort both sets first",
], "Swapping is what makes it O(n log n): the big set is reused as it is, and only the small one's elements move. Copying a big child set into a small parent can cost O(n) per node, and O(n²) on a path-shaped tree. In Python, swap the references (`sets[v], sets[c] = sets[c], sets[v]`); in C++, `std::swap` on two sets is O(1)."),
])

dp = dict(slug="dynamic-programming", id="dsa.dp", title="Dynamic Programming", lessons=[
P("dp-states", "Designing a DP State", 2, ["dp"], r"""
Every DP answers three questions:

1. **State**: what do you need to know about a partial solution to finish it? (`dp[i][w]`: best value from the first i items with total weight w.)
2. **Transition**: how does a state follow from smaller ones? (take item i, or skip it)
3. **Base and order**: which states are trivial, and in what order do you fill in the rest?

Then check the cost: states × transitions per state. If that's too big, change the state, not the code.

**Swap state and value.** If `dp[i][w] = best value` has too many weights but the values are small, use `dp[i][v] = least weight that reaches value v`. The next problem is exactly this.

**Common shapes in interviews**

| Shape | Examples |
| --- | --- |
| prefix `dp[i]` | climbing stairs, house robber, word break, decode ways |
| two sequences `dp[i][j]` | LCS, edit distance, regex matching |
| interval `dp[l][r]` | burst balloons, matrix chain, merging stones |
| subset `dp[mask]` | travelling salesman, assignments with n ≤ 20 |
| tree `dp[v][state]` | house robber III, tree matching |
| digits `dp[pos][tight][...]` | count numbers ≤ N with a property |

Practise recovering the solution, not just its value: store the choice made at each state, or walk back through the table comparing values.

The rest of this chapter follows the AtCoder Educational DP Contest problems on your sheet ([contest page](https://atcoder.jp/contests/dp/tasks)).
""", "A knapsack has n = 100 items, capacity W = 10⁹ and every value ≤ 1000. **Which state is feasible?**", [
    "dp[i][v] = least weight that reaches total value exactly v",
    "dp[i][w] = best value with total weight ≤ w",
    "dp[w] = best value with weight exactly w",
    "dp[i] = best value from the first i items",
], "Capacity 10⁹ rules out any table indexed by weight. The total value is at most 100 · 1000 = 10⁵, so 'least weight for value v' has 10⁵ entries per item, about 10⁷ steps. The answer is the largest v with `dp[v] ≤ W`. A single-index `dp[i]` throws away what you need to know: how much weight is left."),

X("knapsack-2", "E - Knapsack 2", ATC + "e", 3, ["dp", "knapsack"], r"""
## The task

0/1 knapsack with a huge capacity and small values: the case from "Designing a DP State" (state swapping: least weight for a fixed value).

## Before you code

- One array indexed by value. For each item, loop v **downwards** so each item is used at most once.
""", "In the 1D loop `for v from V down to value[i]: dp[v] = min(dp[v], dp[v - value[i]] + weight[i])`, **why go downwards?**", [
    "So dp[v - value[i]] still holds the value from before this item, which stops the item being used twice",
    "It's faster",
    "To find the largest value first",
    "Upwards is also correct for 0/1 knapsack",
], "Going down, `dp[v - value[i]]` hasn't been updated for item i yet, so it describes a solution without item i. Going up, it may already include item i, and the item gets used again: that's the *unbounded* knapsack. Downwards for 0/1, upwards for unlimited copies."),

X("lcs", "F - LCS", ATC + "f", 3, ["dp", "strings"], r"""
## The task

Print one longest common subsequence of strings s and t (each up to 3000 characters): string DP plus solution reconstruction.

## Before you code

- Fill `dp[i][j]` = LCS length of `s[:i]` and `t[:j]`, then walk back from `dp[n][m]` to rebuild the string.
""", "Walking back at (i, j), you find `s[i-1] == t[j-1]`. **What do you do?**", [
    "Add that character and move to (i − 1, j − 1)",
    "Move to whichever of (i − 1, j) and (i, j − 1) is bigger",
    "Move to (i − 1, j)",
    "Stop: the walk is done",
], "When the last characters match, `dp[i][j] = dp[i-1][j-1] + 1`, so taking the match is always optimal. Otherwise, step towards the neighbour with the larger value (either one on a tie). The characters come out backwards, so reverse at the end. A 3001 × 3001 table is 9·10⁶ cells: fine in C++, slow with plain Python lists."),

X("longest-path", "G - Longest Path", ATC + "g", 2, ["dp", "graphs", "dag-dp"], r"""
## The task

Find the longest path, counted in edges, in a directed acyclic graph: DP on a DAG.

## Before you code

- Topological order (Kahn), then `dp[v] = max(dp[u] + 1)` over edges u → v.
- Or a memoized DFS on "longest path starting at v".
""", "**Why doesn't this approach work on a graph with cycles?**", [
    "There's no topological order, and longest *simple* path in a general graph is NP-hard",
    "Kahn's algorithm runs forever",
    "The dp values become negative",
    "It does; cycles only add a constant",
], "A DP needs an order where each state depends only on states already computed, and a DAG's topological order is exactly that. A cycle makes nodes depend on each other. Kahn's algorithm doesn't loop forever; it just stops with some nodes never processed. On graphs with cycles, collapse the SCCs first if revisiting nodes is allowed; otherwise the longest simple path is NP-hard."),

X("sushi", "J - Sushi", ATC + "j", 4, ["dp", "probability"], r"""
## The task

n dishes each start with 1 to 3 pieces of sushi. Each step, roll a die that picks a dish uniformly at random, and eat one piece from it if there are any left. Find the expected number of rolls until everything is eaten: expected-value DP.

## Before you code

- Dishes with the same number of pieces are interchangeable, so the state is just how many dishes have 1, 2 and 3 pieces: `E[a][b][c]`.
- A roll that picks an empty dish leads back to the same state. Solve that equation for E.
""", "From state (a, b, c) with z = n − a − b − c empty dishes, `E = 1 + (z/n)·E + (a/n)·E[a−1][b][c] + …`. **How do you handle the E on the right-hand side?**", [
    "Move it to the left and solve: E = (n + a·E[a−1][b][c] + b·E[a+1][b−1][c] + c·E[a][b+1][c−1]) / (a + b + c)",
    "Ignore it, since empty dishes change nothing",
    "Iterate until the values converge",
    "Replace it with 0",
], "The self-loop just means some rolls are wasted. Solve the linear equation once: E·(1 − z/n) = 1 + …, then multiply through by n and divide by a + b + c. Eating from a 2-piece dish turns it into a 1-piece dish (b − 1, a + 1), and a 3 becomes a 2. Loop c outermost, then b, then a, so every state you read is already computed. Ignoring the self-loop gives too small an answer, and iterating to convergence is slow and inexact."),

X("slimes", "N - Slimes", ATC + "n", 3, ["dp", "interval-dp"], r"""
## The task

n slimes sit in a row. Merging two neighbours costs the sum of their sizes. Merge everything into one slime at the least total cost: interval DP (and the base for Knuth's optimization).

## Before you code

- `dp[l][r]` = cheapest way to merge `a[l..r]` into one slime = min over splits k of `dp[l][k] + dp[k+1][r]`, plus `sum(l..r)` from prefix sums.
""", "**In which order must you fill `dp[l][r]`?**", [
    "By increasing interval length",
    "By increasing l, then increasing r",
    "By decreasing r",
    "Any order",
], "Every split uses two strictly shorter intervals, so all shorter intervals must be done first: loop the length from 2 to n, then l. (Increasing r with *decreasing* l also works, but increasing l with increasing r reads `dp[k+1][r]` before it's computed.) It's O(n³) = 6.4·10⁷ for n = 400. Knuth's optimization cuts it to O(n²) by limiting where the best split can be."),

X("matching", "O - Matching", ATC + "o", 4, ["dp", "bitmask-dp"], r"""
## The task

Count the perfect matchings between n men and n women (n ≤ 21), given who is compatible with whom, modulo 10⁹+7: bitmask DP over subsets.

## Before you code

- `dp[mask]` = ways to match the first `popcount(mask)` men with exactly the women in mask.
""", "**Why does the state need only the mask, and not also which man is next?**", [
    "Men are matched in a fixed order, so the number of set bits in mask already says which man is next",
    "Because every man is compatible with every woman",
    "Because the mask stores men, not women",
    "It does need it: the state must be dp[i][mask]",
], "Matching men in a fixed order means each matching is built in exactly one way, so nothing is double-counted. The popcount tells you which man is next, so a separate man index is redundant, and the state space is 2²¹ ≈ 2·10⁶. The transition tries each compatible woman not yet in the mask: O(2ⁿ · n)."),

X("independent-set", "P - Independent Set", ATC + "p", 2, ["dp", "trees"], r"""
## The task

Colour every node of a tree black or white so that no two adjacent nodes are both black. Count the colourings modulo 10⁹+7: tree DP.

## Before you code

- `dp[v][white]` and `dp[v][black]`, combined over v's children.
""", "**Which transition is right?**", [
    "white[v] = Π (white[c] + black[c]); black[v] = Π white[c]",
    "white[v] = Σ (white[c] + black[c]); black[v] = Σ white[c]",
    "white[v] = Π white[c]; black[v] = Π (white[c] + black[c])",
    "black[v] = white[v] for every v",
], "A white node lets each child be either colour, and the children's subtrees are independent, so their counts multiply. A black node forces every child to be white. Summing would count one child's choices at a time instead of the combinations. The answer is `white[root] + black[root]`, reduced mod 10⁹+7 after every multiplication."),

X("flowers", "Q - Flowers", ATC + "q", 4, ["dp", "fenwick"], r"""
## The task

Flowers stand in a row, each with a distinct height and a beauty. Keep a subsequence whose heights strictly increase, maximizing the total beauty: LIS optimized with a Fenwick or segment tree.

## Before you code

- `dp[h]` = best total ending with the flower of height h. For each flower: `dp[h] = beauty + max(dp[1..h−1])`.
- Keep `dp` in a Fenwick tree for prefix maximums.
""", "**Why can a Fenwick tree do prefix *max* here, when range max is normally a segment-tree job?**", [
    "Queries are always prefixes and each dp[h] only ever increases, so no node ever needs a value removed",
    "Fenwick trees support any range operation",
    "Because the heights are distinct",
    "It can't; this needs a segment tree",
], "A Fenwick tree can't subtract a maximum, but a prefix query never needs to: it combines whole nodes covering 1..h−1. It also needs updates that only raise values, because a node can't forget an old larger value. Here each height is written once, from 'nothing' to its value. O(n log n), and the totals need 64-bit."),

X("digit-sum", "S - Digit Sum", ATC + "s", 4, ["dp", "digit-dp"], r"""
## The task

Count the integers from 1 to K (K can have 10⁴ digits) whose digit sum is a multiple of D, modulo 10⁹+7: digit DP.

## Before you code

- Walk K's digits from the left with the state (position, digit sum mod D, tight). `tight` means the digits so far equal K's prefix, so the next digit can't exceed K's digit.
""", "The DP as described counts every number from 0 to K. **What's the final adjustment?**", [
    "Subtract 1: 0 is counted (its digit sum 0 is a multiple of D) but isn't in 1..K",
    "Add 1 for K itself",
    "Nothing",
    "Divide by D",
], "The all-zeros path (leading zeros allowed) represents 0, and a digit sum of 0 is divisible by every D, so the DP counts it. Subtract 1, then take the result mod 10⁹+7 again, adding the modulus first so it can't go negative. Leading zeros don't change a digit sum, so this problem needs no 'started' flag; problems about digit *patterns* usually do."),

X("frog-3", "Z - Frog 3", ATC + "z", 5, ["dp", "convex-hull-trick"], r"""
## The task

A frog jumps forward along stones of strictly increasing height. A jump from i to j costs (h_j − h_i)² + C. Minimize the total cost of reaching the last stone, with n up to 2·10⁵: the convex hull trick.

## Before you code

- `dp[j] = min over i < j of (dp[i] + (h_j − h_i)²) + C` is O(n²).
- Expand the square: for a fixed i, it's a straight line in h_j. The **convex hull trick** keeps those lines and finds the lowest one at h_j quickly.
""", "Expanding `dp[i] + (h_j − h_i)²`, **which line does stone i contribute, as a function of x = h_j?**", [
    "Slope −2·h_i and intercept dp[i] + h_i², with h_j² added afterwards",
    "Slope h_i and intercept dp[i]",
    "Slope 2·h_i and intercept dp[i] − h_i²",
    "Slope h_i² and intercept dp[i]",
], "(h_j − h_i)² = h_j² − 2·h_i·h_j + h_i². The h_j² term is the same for every i, so pull it out; what's left is the line y = (−2·h_i)·x + (dp[i] + h_i²). Because h increases, the slopes decrease and the queries increase, so a monotonic deque of lines gives O(n) total. A Li Chao tree works without those guarantees, in O(n log n)."),
])

strings = dict(slug="strings", id="dsa.strings", title="Strings", lessons=[
X("word-combinations", "Word Combinations", CSES + "1731", 3, ["strings", "trie", "dp"], r"""
## The task

Count the ways to split a string (up to 5000 characters) into words from a dictionary (up to 10⁵ words), modulo 10⁹+7.

## Pattern recap: trie

A **trie** is a tree of characters: each path from the root spells a prefix, and nodes mark where a word ends. Insert every word; then from any start position, walk the trie along the string to visit every dictionary word that starts there, stopping as soon as the trie has no matching child.

```cpp
vector<array<int, 26>> trie(1);   // child index per letter, 0 = none
vector<bool> is_end(1, false);

void insert(const string& w) {
    int node = 0;
    for (char ch : w) {
        int c = ch - 'a';
        if (!trie[node][c]) {
            trie[node][c] = trie.size();
            trie.push_back({});
            is_end.push_back(false);
        }
        node = trie[node][c];
    }
    is_end[node] = true;
}
```

```python
trie = [{}]              # each node: dict from character to child index
end = [False]


def insert(w):
    node = 0
    for ch in w:
        if ch not in trie[node]:
            trie[node][ch] = len(trie)
            trie.append({})
            end.append(False)
        node = trie[node][ch]
    end[node] = True
```

The counting part is the "word break" DP (LeetCode 139/140), with the trie making each position fast.
""", "**Why use a trie instead of checking every dictionary word at each position?**", [
    "One trie walk finds every matching word from a position and stops when no word continues; checking each word separately costs up to 10⁵ comparisons per position",
    "Tries use less memory than the word list",
    "Tries sort the words",
    "There's no difference in speed",
], "With 5000 positions and 10⁵ words, testing every word everywhere is far too slow. The trie walk from position i costs at most the length of the longest matching prefix, and usually much less. The DP itself is `dp[i] = Σ dp[i + len(w)]` over the words w that match at i, with `dp[n] = 1`."),

X("finding-periods", "Finding Periods", CSES + "1733", 3, ["strings", "hashing"], r"""
## The task

Find every period of a string: every length p such that the string is a prefix of its first p characters repeated over and over.

## Pattern recap: rolling hash

`h[i]` is the hash of the first i characters: `h[i+1] = (h[i] · B + code(s[i])) mod M`. Any substring's hash is `h[r] − h[l] · B^(r−l)` (mod M), in O(1) with precomputed powers. Compare substrings by comparing hashes. Use a random base and a large prime modulus (or two moduli) so collisions, and anti-hash tests, don't bite.

The Z-function or the prefix function solves this exactly, with no collision risk.
""", "**p is a period exactly when…**", [
    "s[p:] equals s[:n − p] (the string matches itself shifted by p)",
    "n is divisible by p",
    "s[:p] appears twice in s",
    "s[p] == s[0]",
], "Repeating the first p characters means `s[i] == s[i − p]` for every i ≥ p, which is exactly the shifted comparison. p doesn't have to divide n: 'abcab' has period 3 ('abc' + 'ab'). With hashing, each p is one O(1) comparison, O(n) in total. With the Z-function, p is a period when `z[p] == n − p`."),

P("kmp", "KMP and the Prefix Function", 3, ["strings", "kmp"], r"""
The **prefix function** `pi[i]` is the length of the longest proper prefix of `s[:i+1]` that is also a suffix of it (a **border**).

```cpp
vector<int> prefix_function(const string& s) {
    vector<int> pi(s.size(), 0);
    for (size_t i = 1; i < s.size(); i++) {
        int k = pi[i - 1];
        while (k && s[i] != s[k]) k = pi[k - 1];
        if (s[i] == s[k]) k++;
        pi[i] = k;
    }
    return pi;
}
```

```python
def prefix_function(s):
    pi = [0] * len(s)
    for i in range(1, len(s)):
        k = pi[i - 1]
        while k and s[i] != s[k]:
            k = pi[k - 1]
        if s[i] == s[k]:
            k += 1
        pi[i] = k
    return pi
```

It's O(n) by the same amortized argument as the monotonic stack: k rises by at most 1 per step, so it can't fall more than n times in total.

**Pattern matching**: compute it for `pattern + "#" + text`, with a separator that appears in neither. Every position where `pi == len(pattern)` ends a match.

**All borders** of s: `pi[n−1]`, then `pi[pi[n−1] − 1]`, and so on down to 0.

The **Z-function** (`z[i]` = the longest common prefix of s and `s[i:]`) solves the same problems. Learn whichever one you can write from memory.
""", "**What is the prefix function of `\"aabaaab\"`?**", [
    "[0, 1, 0, 1, 2, 2, 3]",
    "[0, 1, 0, 1, 2, 3, 3]",
    "[0, 0, 0, 1, 2, 2, 3]",
    "[0, 1, 0, 1, 2, 3, 4]",
], "At i = 5 the border 'aa' can't grow, because the character after it is 'b', not 'a'. So k falls back to `pi[1] = 1` and then grows to 2 ('aa'). At i = 6 the border 'aa' grows into 'aab', giving 3. The second option skips the fallback. And `pi[0]` is always 0, because a border must be shorter than the string."),

X("string-matching", "String Matching", CSES + "1753", 2, ["strings", "kmp"], r"""
## The task

Count how many times a pattern occurs in a text. Occurrences may overlap.

## Before you code

- The prefix function on `pattern + "#" + text`, or KMP run directly over the text.
""", "The text is `\"aaaa\"` and the pattern is `\"aa\"`. **What's the count?**", [
    "3", "2", "4", "1",
], "Occurrences start at positions 0, 1 and 2, and overlapping ones count. Python's `text.count('aa')` gives 2, because it skips past each match; that's the trap. KMP finds overlaps naturally: after a match it falls back to `pi[m − 1]` and keeps going."),

X("finding-borders", "Finding Borders", CSES + "1732", 2, ["strings", "kmp"], r"""
## The task

Print the lengths of every border of a string (prefixes that are also suffixes, not counting the whole string), in increasing order.

## Before you code

- Follow the chain from `pi[n − 1]`, as in "KMP and the Prefix Function".
""", "**What are the borders of `\"abacaba\"`?**", [
    "1 3", "3", "1 2 3", "1 3 7",
], "`pi[6] = 3` ('aba'), then `pi[2] = 1` ('a'), then `pi[0] = 0` ends the chain. Sorted, that's 1 3. The whole string isn't a border because a border must be proper, so 7 doesn't count. Every border of s is somewhere in that chain, which is why it lists them all in O(n)."),

P("palindromes", "Palindromes", 3, ["strings", "palindromes"], r"""
Three tools for palindrome problems:

**Expand around the centre**: O(n²), fine for n ≤ 5000, and the usual interview answer (longest palindromic substring, count palindromic substrings). Each of the 2n − 1 centres expands while both ends match.

**Manacher's algorithm**: O(n), the longest palindrome around *every* centre. It reuses the mirror image of each position inside the rightmost palindrome found so far, so it never re-checks a pair of characters already known to match.

```cpp
vector<int> manacher(const string& s) {
    string t = "#";
    for (char ch : s) {                 // gives even-length palindromes a centre
        t += ch;
        t += '#';
    }
    int n = t.size(), c = 0, r = 0;
    vector<int> p(n, 0);
    for (int i = 0; i < n; i++) {
        if (i < r) p[i] = min(r - i, p[2 * c - i]);
        while (i - p[i] - 1 >= 0 && i + p[i] + 1 < n && t[i - p[i] - 1] == t[i + p[i] + 1])
            p[i]++;
        if (i + p[i] > r) c = i, r = i + p[i];
    }
    return p;   // p[i] = radius in t = length of that palindrome in s
}
```

```python
def manacher(s):
    t = "#" + "#".join(s) + "#"      # gives even-length palindromes a centre
    n, p = len(t), [0] * len(t)
    c = r = 0
    for i in range(n):
        if i < r:
            p[i] = min(r - i, p[2 * c - i])
        while i - p[i] - 1 >= 0 and i + p[i] + 1 < n and t[i - p[i] - 1] == t[i + p[i] + 1]:
            p[i] += 1
        if i + p[i] > r:
            c, r = i, i + p[i]
    return p     # p[i] = radius in t = length of that palindrome in s
```

**Hashing**: compare the hash of a substring with the hash of its reverse. With a Fenwick or segment tree of hashes, this even survives updates.
""", "**Why does Manacher's algorithm put `#` between the characters?**", [
    "So even-length palindromes get a centre too, and one loop handles both kinds",
    "To mark the end of the string",
    "To make the string long enough for the mirror trick",
    "To avoid comparing equal characters",
], "'abba' has no middle character, but in '#a#b#b#a#' its centre is the `#` between the two b's. Every palindrome in t has odd length, and its radius `p[i]` equals the length of the original palindrome in s. The bounds checks in the `while` stop it from running off either end."),

X("longest-palindrome", "Longest Palindrome", CSES + "1111", 3, ["strings", "palindromes"], r"""
## The task

Print the longest palindromic substring of a string of up to 10⁶ characters.

## Before you code

- Expanding around every centre is O(n²) in the worst case (think `"aaaa…"`), so you need Manacher here.
""", "In the transformed string t, the best radius is `p[i] = 5` at `i = 9`. **Where does the palindrome start in s?**", [
    "At (i − p[i]) // 2 = 2",
    "At i − p[i] = 4",
    "At i // 2 = 4",
    "At (i − p[i]) // 2 + 1 = 3",
], "In t the palindrome covers positions i − p[i] to i + p[i], that is 4 to 14. The characters of s sit at odd positions (s[k] is at 2k + 1), so the first real character is at position 5, which is s[2]. The length is p[i] = 5, so the answer is `s[2:7]`. Off-by-one slips in this conversion cause most wrong answers with Manacher."),

X("palindrome-queries", "Palindrome Queries", CSES + "2420", 4, ["strings", "hashing", "fenwick"], r"""
## The task

A string supports two operations: change one character, or check whether `s[a..b]` is a palindrome.

## Before you code

- Two Fenwick trees of hashes, one over s and one over its reverse. A substring is a palindrome when its forward and backward hashes agree.
""", "Position k stores `code(s[k]) · B^k` in the forward tree. **How do you compare the forward hash of `s[a..b]` with the backward one?**", [
    "Scale both to the same power of B first (multiply the one with the smaller offset), since they're shifted by different powers",
    "Compare the two range sums directly",
    "Divide both sums by b − a",
    "Compare only the first and last characters",
], "The forward range sum is `Σ code(s[k]) · B^k` for k in [a, b], so it carries an extra factor of B^a. In the reversed string the same characters sit at positions n − 1 − b … n − 1 − a, with an extra factor of B^(n−1−b). Multiply each side by the power the other one has (or use modular inverses), then compare. Changing position k changes one term in each tree: O(log n) per operation."),

P("suffix-arrays", "Suffix Arrays", 4, ["strings", "suffix-array"], r"""
A **suffix array** lists a string's suffixes in sorted order, stored as their start indices. The **LCP array** holds the longest common prefix of each pair of neighbours in that order.

- **Build** by prefix doubling: sort by the first 2^k characters using the ranks for 2^(k−1). That's O(n log² n) with comparison sorting, or O(n log n) with radix sort. **Kasai's algorithm** then builds the LCP array in O(n).
- **Distinct substrings**: every substring is a prefix of some suffix, and neighbours share `lcp[i]` of their prefixes, so the count is `Σ (n − sa[i]) − Σ lcp[i]`.
- **Longest repeated substring**: the largest value in the LCP array.

Your sheet marks this optional, for a good reason: it rarely comes up in interviews. It's still a great test of writing careful O(n log n) code.
""", "**Why is the longest repeated substring equal to the largest LCP value?**", [
    "A repeated substring is a common prefix of two suffixes, and the longest common prefix of any two suffixes is always found between neighbours in sorted order",
    "Because the LCP array is sorted",
    "Because the suffix array lists substrings by length",
    "It isn't; you have to compare every pair of suffixes",
], "Suffixes that share a long prefix sort next to each other, or with suffixes in between that share at least as much. So the best pair among all suffixes is always some pair of neighbours, and only the n − 1 neighbour values matter. That's what lets suffix arrays answer substring questions in linear time after the build."),

X("distinct-substrings", "Distinct Substrings", CSES + "2105", 4, ["strings", "suffix-array"], r"""
## The task

Count the distinct substrings of a string of up to 10⁵ characters.

## Before you code

- The formula from "Suffix Arrays". A suffix automaton works too: sum `len(state) − len(link(state))` over all states.
""", "For `\"aaa\"`, `sa = [2, 1, 0]` and the LCP values between neighbours are `[1, 2]`. **How many distinct substrings are there?**", [
    "3", "6", "4", "5",
], "The suffixes 'a', 'aa' and 'aaa' have 1 + 2 + 3 = 6 prefixes between them. Neighbours share 1 and 2 of those, so 6 − 3 = 3: 'a', 'aa' and 'aaa'. Use 64-bit integers: the count can reach about n²/2 = 5·10⁹."),

X("repeating-substring", "Repeating Substring", CSES + "2106", 4, ["strings", "suffix-array", "hashing"], r"""
## The task

Find the longest substring that appears at least twice, or print −1 if there's none.

## Before you code

- The largest LCP between neighbours in the suffix array.
- Or binary search on the length, with rolling hashes.
""", "**Why can you binary search on the answer's length when using hashing?**", [
    "If some substring of length L repeats, its prefix of length L − 1 repeats too, so 'a repeat of length L exists' is monotone",
    "Because hashes come out sorted",
    "Because the answer is always a power of two",
    "You can't; you have to try every length",
], "Chop the last character off both copies of a repeated substring and you get a shorter repeat. So the lengths that have a repeat are exactly 0 up to the answer, and binary search applies. Checking one length L puts the hash of every length-L window into a set: O(n), so O(n log n) overall. It's the same 'binary search on the answer' idea as Factory Machines."),

P("string-dp", "String DP", 3, ["strings", "dp"], r"""
Two-string DP fills `dp[i][j]` for the prefixes `a[:i]` and `b[:j]`.

**Edit distance** (insert, delete or replace one character, each costing 1):

```cpp
int edit_distance(const string& a, const string& b) {
    int n = a.size(), m = b.size();
    vector<vector<int>> dp(n + 1, vector<int>(m + 1));
    for (int i = 0; i <= n; i++) dp[i][0] = i;
    for (int j = 0; j <= m; j++) dp[0][j] = j;
    for (int i = 1; i <= n; i++)
        for (int j = 1; j <= m; j++) {
            int same = a[i - 1] == b[j - 1];
            dp[i][j] = min({dp[i - 1][j] + 1,             // delete a[i-1]
                            dp[i][j - 1] + 1,             // insert b[j-1]
                            dp[i - 1][j - 1] + !same});   // replace, or keep
        }
    return dp[n][m];
}
```

```python
def edit_distance(a, b):
    n, m = len(a), len(b)
    dp = [[0] * (m + 1) for _ in range(n + 1)]
    for i in range(n + 1):
        dp[i][0] = i
    for j in range(m + 1):
        dp[0][j] = j
    for i in range(1, n + 1):
        for j in range(1, m + 1):
            same = a[i - 1] == b[j - 1]
            dp[i][j] = min(dp[i - 1][j] + 1,                 # delete a[i-1]
                           dp[i][j - 1] + 1,                 # insert b[j-1]
                           dp[i - 1][j - 1] + (not same))    # replace, or keep
    return dp[n][m]
```

Each row only reads the row above it, so two rows of memory are enough.

**DP over a pattern automaton**: "count the strings of length n that contain a pattern" walks the KMP automaton, whose state is how much of the pattern is matched so far. Counting the strings that *avoid* the pattern and subtracting from the total is usually easier.
""", "**What is the edit distance between `\"kitten\"` and `\"sitting\"`?**", [
    "3", "2", "4", "7",
], "Replace k → s, replace e → i, and insert g at the end: 3 operations. Two can't be enough: one operation must be an insert, because the lengths differ, and that leaves only one more for the two mismatched letters. The base row and column (`dp[i][0] = i`, `dp[0][j] = j`) are the cost of turning a prefix into the empty string; forgetting them is the most common bug."),

X("edit-distance", "Edit Distance", CSES + "1639", 2, ["strings", "dp"], r"""
## The task

Find the edit distance between two strings of up to 5000 characters each.

## Before you code

- The "String DP" template with two rows.
- 2.5·10⁷ cells is fine in C++; in Python it's borderline, so keep the inner loop tight.
""", "With two rows, `prev` and `cur`, **which value does the \"insert\" move read?**", [
    "cur[j − 1]", "prev[j − 1]", "prev[j]", "cur[j + 1]",
], "Insert comes from `dp[i][j − 1]`, which is in the row being filled: `cur[j − 1]`. Delete reads `prev[j]`, and replace/keep reads `prev[j − 1]`. Mixing these up gives answers that are only wrong on a few tests. Also remember to set `cur[0] = i` at the start of each row."),

X("required-substring", "Required Substring", CSES + "1112", 4, ["strings", "dp", "kmp"], r"""
## The task

Count the strings of length n over the letters A–Z that contain a given pattern of length m, modulo 10⁹+7.

## Before you code

- Count the strings that *don't* contain it: a DP over (length so far, KMP state = how many pattern characters are matched), never letting the state reach m.
- The answer is 26ⁿ minus that count.
""", "From KMP state k (k characters matched), you add the character c. **Which state do you move to?**", [
    "The length of the longest pattern prefix that's a suffix of pattern[:k] + c, found with the prefix function's fallback",
    "k + 1 if c matches, otherwise 0",
    "Always k + 1",
    "c's position in the alphabet",
], "A mismatch doesn't always reset to 0. With the pattern 'ABAB', state 3 ('ABA') plus 'A' gives 'ABAA', whose longest suffix that's also a pattern prefix is 'A', so the new state is 1. Precompute a table `go[k][c]` using the prefix function's fallback: m × 26 entries. The DP is then O(n · m · 26)."),

X("minimal-rotation", "Minimal Rotation", CSES + "1110", 4, ["strings"], r"""
## The task

Find the lexicographically smallest rotation of a string of up to 10⁶ characters.

## Pattern recap: two candidates at a time

Compare two candidate starts i and j character by character on `s + s`. When they differ at offset k, the losing start and the k starts right after it can all be thrown away.

```cpp
string min_rotation(const string& s) {
    int n = s.size(), i = 0, j = 1, k = 0;
    string t = s + s;
    while (i < n && j < n && k < n) {
        char a = t[i + k], b = t[j + k];
        if (a == b) {
            k++;
            continue;
        }
        if (a > b) i += k + 1;
        else j += k + 1;
        if (i == j) j++;
        k = 0;
    }
    return t.substr(min(i, j), n);
}
```

```python
def min_rotation(s):
    n = len(s)
    t = s + s
    i, j, k = 0, 1, 0
    while i < n and j < n and k < n:
        a, b = t[i + k], t[j + k]
        if a == b:
            k += 1
            continue
        if a > b:
            i += k + 1
        else:
            j += k + 1
        if i == j:
            j += 1
        k = 0
    start = min(i, j)
    return t[start:start + n]
```

That's O(n). Booth's algorithm and Lyndon factorization (Duval's algorithm) are the other classic approaches.
""", "Candidates i and j agree for k characters, then `t[i + k] > t[j + k]`. **Why can i jump all the way to i + k + 1?**", [
    "For each 0 ≤ d ≤ k, start i + d loses to start j + d, which meets the same deciding pair of characters",
    "Because those starts were already tested",
    "Because k is always 0",
    "It can't; only i + 1 is safe",
], "Starts i + d and j + d (for d ≤ k) line up the same matching stretch, shifted by d, and hit the same deciding pair `t[i + k] > t[j + k]`. So i + d is never better than j + d, and it can't be the unique answer. Skipping k + 1 starts at a time is what makes the whole scan O(n)."),
])

CHAPTERS = [graphs, paths, dsu, trees, dp, strings]
