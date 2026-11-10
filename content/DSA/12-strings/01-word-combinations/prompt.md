# Word Combinations

**Practice:** [Word Combinations on CSES](https://cses.fi/problemset/task/1731)

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

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why use a trie instead of checking every dictionary word at each position?**
