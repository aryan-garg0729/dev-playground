# Heaps and Ordered Sets

Two structures cover most "keep a changing collection" problems.

| Need | Structure | Python | C++ |
| --- | --- | --- | --- |
| Smallest (or largest) item, many times | heap | `heapq` | `priority_queue` |
| Smallest item ≥ x, or largest ≤ x, with deletes | ordered multiset | none built in | `multiset` + `lower_bound` |
| k-th smallest, running median | two heaps | `heapq` twice | `priority_queue` twice |

## Lazy deletion

A heap can't delete from the middle. Instead, record the item as deleted and skip it when it reaches the top:

```cpp
// In C++ a multiset usually beats this: s.erase(s.find(x)) deletes one copy.
priority_queue<int, vector<int>, greater<int>> heap;   // min-heap
map<int, int> dead;

void remove(int x) { dead[x]++; }

int top() {
    while (dead[heap.top()]) {
        dead[heap.top()]--;
        heap.pop();
    }
    return heap.top();
}
```

```python
import heapq
from collections import Counter

heap, dead = [], Counter()


def remove(x):
    dead[x] += 1


def top():
    while dead[heap[0]]:
        dead[heap[0]] -= 1
        heapq.heappop(heap)
    return heap[0]
```

## No multiset in Python?

Python's standard library has no ordered multiset, and CSES doesn't install `sortedcontainers`. Your options:

- `bisect.insort` on a list: O(n) per insert, but a fast memory move that often passes for n ≈ 2·10⁵.
- Work offline: read every query first, then process them in an order that only needs sorting (see Traffic Lights).
- A "next free slot" DSU to skip deleted items (chapter "Disjoint Set Union & MST").

If you practise CSES in Python, expect time-limit trouble on the hardest problems. C++ is the safer choice there; LeetCode-style interviews are fine in Python.

## Check yourself

You must repeatedly find the largest value ≤ x **and remove it**, over 2·10⁵ operations. **Which structure makes every operation O(log n)?**
