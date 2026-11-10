# Minimal Rotation

**Practice:** [Minimal Rotation on CSES](https://cses.fi/problemset/task/1110)

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

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Candidates i and j agree for k characters, then `t[i + k] > t[j + k]`. **Why can i jump all the way to i + k + 1?**
