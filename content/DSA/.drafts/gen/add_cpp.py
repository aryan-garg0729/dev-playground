"""One-off: insert a ```cpp twin above every ```python block in data_a.py / data_b.py (in order)."""
import os
import re

A = [
r"""struct Node {
    int key, val;
    Node *prev = nullptr, *next = nullptr;
    Node(int k = 0, int v = 0) : key(k), val(v) {}
};

void unlink(Node* node) {
    node->prev->next = node->next;
    node->next->prev = node->prev;
}

void push_front(Node* head, Node* node) {
    node->prev = head;
    node->next = head->next;
    head->next->prev = node;
    head->next = node;
}""",
r"""vector<vector<int>> merge(vector<vector<int>>& intervals) {
    sort(intervals.begin(), intervals.end());
    vector<vector<int>> out;
    for (auto& iv : intervals) {
        if (!out.empty() && iv[0] <= out.back()[1])
            out.back()[1] = max(out.back()[1], iv[1]);
        else
            out.push_back(iv);
    }
    return out;
}""",
r"""// sorts a[lo, hi) and returns its inversion count
long long sort_count(vector<int>& a, int lo, int hi, vector<int>& tmp) {
    if (hi - lo <= 1) return 0;
    int mid = (lo + hi) / 2;
    long long inv = sort_count(a, lo, mid, tmp) + sort_count(a, mid, hi, tmp);
    int i = lo, j = mid, k = lo;
    while (i < mid && j < hi) {
        if (a[i] <= a[j]) {
            tmp[k++] = a[i++];
        } else {
            tmp[k++] = a[j++];
            inv += mid - i;          // every left element still waiting
        }
    }
    while (i < mid) tmp[k++] = a[i++];
    while (j < hi) tmp[k++] = a[j++];
    copy(tmp.begin() + lo, tmp.begin() + hi, a.begin() + lo);
    return inv;
}""",
r"""// smallest x in [lo, hi] with ok(x), assuming ok(hi) is true
long long first_true(long long lo, long long hi, function<bool(long long)> ok) {
    while (lo < hi) {
        long long mid = lo + (hi - lo) / 2;
        if (ok(mid)) hi = mid;
        else lo = mid + 1;
    }
    return lo;
}""",
r"""vector<long long> pre(a.size() + 1, 0);
for (size_t i = 0; i < a.size(); i++) pre[i + 1] = pre[i] + a[i];
// sum of a[l..r-1] (l inclusive, r exclusive)
long long total = pre[r] - pre[l];""",
r"""long long count_sum(const vector<long long>& a, long long x) {
    map<long long, long long> seen;   // map, not unordered_map: safe from anti-hash tests
    seen[0] = 1;
    long long cur = 0, ans = 0;
    for (long long v : a) {
        cur += v;
        auto it = seen.find(cur - x);
        if (it != seen.end()) ans += it->second;
        seen[cur]++;
    }
    return ans;
}""",
r"""long long best = a[0], cur = a[0];
for (size_t i = 1; i < a.size(); i++) {
    cur = max(a[i], cur + a[i]);      // a is a vector<long long>
    best = max(best, cur);
}""",
r"""// In C++ a multiset usually beats this: s.erase(s.find(x)) deletes one copy.
priority_queue<int, vector<int>, greater<int>> heap;   // min-heap
map<int, int> dead;

void remove(int x) { dead[x]++; }

int top() {
    while (dead[heap.top()]) {
        dead[heap.top()]--;
        heap.pop();
    }
    return heap.top();
}""",
r"""vector<int> prev_smaller(const vector<int>& a) {
    vector<int> res, st;              // st holds indices; their values increase
    for (int i = 0; i < (int)a.size(); i++) {
        while (!st.empty() && a[st.back()] >= a[i])
            st.pop_back();            // blocked by a[i] from now on
        res.push_back(st.empty() ? -1 : st.back());
        st.push_back(i);
    }
    return res;
}""",
r"""int longest_with_at_most_k_distinct(const vector<int>& a, int k) {
    unordered_map<int, int> count;
    int l = 0, best = 0;
    for (int r = 0; r < (int)a.size(); r++) {
        count[a[r]]++;
        while ((int)count.size() > k) {
            if (--count[a[l]] == 0) count.erase(a[l]);
            l++;
        }
        best = max(best, r - l + 1);
    }
    return best;
}""",
r"""vector<int> window_mins(const vector<int>& a, int k) {
    deque<int> dq;
    vector<int> out;
    for (int i = 0; i < (int)a.size(); i++) {
        while (!dq.empty() && a[dq.back()] >= a[i]) dq.pop_back();
        dq.push_back(i);
        if (dq.front() <= i - k) dq.pop_front();
        if (i >= k - 1) out.push_back(a[dq.front()]);
    }
    return out;
}""",
r"""int single(const vector<int>& nums) {   // every value appears twice except one
    int acc = 0;
    for (int x : nums) acc ^= x;
    return acc;
}""",
r"""long long power(long long a, long long e, long long m) {
    long long r = 1;
    a %= m;
    while (e) {
        if (e & 1) r = r * a % m;
        a = a * a % m;
        e >>= 1;
    }
    return r;
}""",
r"""const long long MOD = 1e9 + 7;
const int N = 1e6;
vector<long long> fact(N + 1), inv(N + 1);

void init() {
    fact[0] = 1;
    for (int i = 1; i <= N; i++) fact[i] = fact[i - 1] * i % MOD;
    inv[N] = power(fact[N], MOD - 2, MOD);
    for (int i = N; i > 0; i--) inv[i - 1] = inv[i] * i % MOD;
}

long long C(int a, int b) {
    if (b < 0 || b > a) return 0;
    return fact[a] * inv[b] % MOD * inv[a - b] % MOD;
}""",
r"""vector<bool> sieve(int n) {
    vector<bool> is_p(n + 1, true);
    is_p[0] = is_p[1] = false;
    for (long long i = 2; i * i <= n; i++)
        if (is_p[i])
            for (long long j = i * i; j <= n; j += i) is_p[j] = false;
    return is_p;
}""",
r"""vector<vector<int>> build(const vector<int>& a) {
    int n = a.size();
    vector<vector<int>> st = {a};
    for (int j = 1; (1 << j) <= n; j++) {
        int half = 1 << (j - 1);
        vector<int> row(n - (1 << j) + 1);
        for (int i = 0; i < (int)row.size(); i++)
            row[i] = min(st[j - 1][i], st[j - 1][i + half]);
        st.push_back(row);
    }
    return st;
}

int query(const vector<vector<int>>& st, int l, int r) {   // inclusive l..r
    int j = __lg(r - l + 1);
    return min(st[j][l], st[j][r - (1 << j) + 1]);
}""",
r"""struct BIT {
    vector<long long> t;                   // 1-indexed
    BIT(int n) : t(n + 1, 0) {}

    void add(int i, long long d) {
        for (; i < (int)t.size(); i += i & -i) t[i] += d;
    }

    long long prefix(int i) {              // sum of positions 1..i
        long long s = 0;
        for (; i > 0; i -= i & -i) s += t[i];
        return s;
    }
};""",
]

B = [
r"""vector<int> bfs(int start, const vector<vector<int>>& adj) {
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
}""",
r"""vector<int> topo(int n, const vector<vector<int>>& adj) {
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
}""",
r"""vector<int> color;   // 0 white, 1 grey, 2 black

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
}""",
r"""vector<long long> dijkstra(int src, const vector<vector<pair<int, int>>>& adj) {
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
}""",
r"""// dist[i][j]: edge weight, INF if none, 0 on the diagonal.
// INF = 1e18 keeps INF + INF inside a long long.
void floyd(vector<vector<long long>>& dist, int n) {
    for (int k = 0; k < n; k++)
        for (int i = 0; i < n; i++)
            for (int j = 0; j < n; j++)
                dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j]);
}""",
r"""struct Edge { int u, v; long long w; };

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
}""",
r"""struct DSU {
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
};""",
r"""const int LOG = 18;                              // 2^18 > 2·10^5
vector<vector<int>> up(LOG, vector<int>(n + 1)); // up[0][v] = parent; the root's parent is 0
// fill: up[j][v] = up[j - 1][up[j - 1][v]]

int kth_ancestor(int v, int k) {
    for (int j = 0; j < LOG && k && v; j++, k >>= 1)
        if (k & 1) v = up[j][v];
    return k ? 0 : v;                            // 0 means "no such ancestor"
}""",
r"""int lca(int a, int b) {
    if (depth[a] < depth[b]) swap(a, b);
    a = kth_ancestor(a, depth[a] - depth[b]);
    if (a == b) return a;
    for (int j = LOG - 1; j >= 0; j--)
        if (up[j][a] != up[j][b]) {
            a = up[j][a];
            b = up[j][b];
        }
    return up[0][a];
}""",
r"""vector<int> order = {root}, parent(n + 1, 0);
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
}""",
r"""vector<int> tin, tout;   // sized n + 1
int timer = 0;

void dfs(int v, int p) {
    tin[v] = timer++;
    for (int c : adj[v])
        if (c != p) dfs(c, v);
    tout[v] = timer - 1;
}""",
r"""vector<array<int, 26>> trie(1);   // child index per letter, 0 = none
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
}""",
r"""vector<int> prefix_function(const string& s) {
    vector<int> pi(s.size(), 0);
    for (size_t i = 1; i < s.size(); i++) {
        int k = pi[i - 1];
        while (k && s[i] != s[k]) k = pi[k - 1];
        if (s[i] == s[k]) k++;
        pi[i] = k;
    }
    return pi;
}""",
r"""vector<int> manacher(const string& s) {
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
}""",
r"""int edit_distance(const string& a, const string& b) {
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
}""",
r"""string min_rotation(const string& s) {
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
}""",
]

here = os.path.dirname(os.path.abspath(__file__))
for name, snippets in (("data_a.py", A), ("data_b.py", B)):
    path = os.path.join(here, name)
    src = open(path).read()
    assert "```cpp" not in src, f"{name} already has C++ snippets"
    parts = re.split(r"(?m)^```python\n", src)
    assert len(parts) - 1 == len(snippets), (name, len(parts) - 1, len(snippets))
    out = parts[0]
    for code, rest in zip(snippets, parts[1:]):
        out += "```cpp\n" + code + "\n```\n\n```python\n" + rest
    open(path, "w").write(out)
    print(name, len(snippets), "C++ snippets added")
