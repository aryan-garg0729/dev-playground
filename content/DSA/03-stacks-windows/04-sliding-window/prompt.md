# Sliding Window

A **window** is a contiguous range `a[l:r]` that slides to the right. Three tools cover most window problems.

**1. Two pointers (variable size).** Grow r; while the window breaks the rule, shrink l. It works when shrinking a bad window can make it good again, and growing a good one can only make it worse.

```cpp
int longest_with_at_most_k_distinct(const vector<int>& a, int k) {
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
}
```

```python
def longest_with_at_most_k_distinct(a, k):
    count, l, best = {}, 0, 0
    for r, x in enumerate(a):
        count[x] = count.get(x, 0) + 1
        while len(count) > k:
            count[a[l]] -= 1
            if count[a[l]] == 0:
                del count[a[l]]
            l += 1
        best = max(best, r - l + 1)
    return best
```

**2. Monotonic deque (fixed size, min or max).** Keep indices whose values increase from front to back. Drop the front when it leaves the window; pop the back when a smaller value arrives.

```cpp
vector<int> window_mins(const vector<int>& a, int k) {
    deque<int> dq;
    vector<int> out;
    for (int i = 0; i < (int)a.size(); i++) {
        while (!dq.empty() && a[dq.back()] >= a[i]) dq.pop_back();
        dq.push_back(i);
        if (dq.front() <= i - k) dq.pop_front();
        if (i >= k - 1) out.push_back(a[dq.front()]);
    }
    return out;
}
```

```python
from collections import deque


def window_mins(a, k):
    dq, out = deque(), []
    for i, x in enumerate(a):
        while dq and a[dq[-1]] >= x:
            dq.pop()
        dq.append(i)
        if dq[0] <= i - k:
            dq.popleft()
        if i >= k - 1:
            out.append(a[dq[0]])
    return out
```

**3. Two heaps (median or k-th).** The lower half in a max-heap, the upper half in a min-heap, with lazy deletion for the element that leaves the window.

Counting trick: "subarrays with *exactly* k distinct values" = atMost(k) − atMost(k − 1).

## Check yourself

**Which problem can't be solved with a plain two-pointer window?**
