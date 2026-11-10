# Palindromes

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

## Check yourself

**Why does Manacher's algorithm put `#` between the characters?**
