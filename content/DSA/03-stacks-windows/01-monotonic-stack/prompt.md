# Monotonic Stack

"For every element, find the nearest one to its left that's smaller" sounds like O(n²). A **monotonic stack** does it in O(n): a stack whose values stay sorted because you pop everything that can never be an answer again.

```cpp
vector<int> prev_smaller(const vector<int>& a) {
    vector<int> res, st;              // st holds indices; their values increase
    for (int i = 0; i < (int)a.size(); i++) {
        while (!st.empty() && a[st.back()] >= a[i])
            st.pop_back();            // blocked by a[i] from now on
        res.push_back(st.empty() ? -1 : st.back());
        st.push_back(i);
    }
    return res;
}
```

```python
def prev_smaller(a):
    res, st = [], []          # st holds indices; their values increase
    for i, x in enumerate(a):
        while st and a[st[-1]] >= x:
            st.pop()          # a[st[-1]] is blocked by x from now on
        res.append(st[-1] if st else -1)
        st.append(i)
    return res
```

Why popping is safe: once x arrives, any value ≥ x to its left is both farther away and no smaller, so it can never be the nearest smaller value for anything to the right.

Variants:

- next greater / next smaller: the same loop, or run it right to left
- daily temperatures, stock span, remove k digits, sum of subarray minimums
- largest rectangle in a histogram, maximal rectangle of 1s in a grid

## Check yourself

The loop has a `while` inside a `for`. **What's the total running time?**
