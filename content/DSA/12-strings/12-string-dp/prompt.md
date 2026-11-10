# String DP

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

## Check yourself

**What is the edit distance between `"kitten"` and `"sitting"`?**
