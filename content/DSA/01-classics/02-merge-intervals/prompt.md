# Merge Intervals

**Practice:** [Merge Intervals on LeetCode](https://leetcode.com/problems/merge-intervals/)

## The task

Given a list of intervals `[start, end]`, merge every overlapping group and return the merged intervals. Intervals that only touch (`[1, 4]` and `[4, 5]`) count as overlapping.

## Pattern recap: sort, then sweep

Sort by start. Then each interval either extends the last merged one or starts a new one:

```cpp
vector<vector<int>> merge(vector<vector<int>>& intervals) {
    sort(intervals.begin(), intervals.end());
    vector<vector<int>> out;
    for (auto& iv : intervals) {
        if (!out.empty() && iv[0] <= out.back()[1])
            out.back()[1] = max(out.back()[1], iv[1]);
        else
            out.push_back(iv);
    }
    return out;
}
```

```python
def merge(intervals):
    intervals.sort(key=lambda iv: iv[0])
    out = []
    for start, end in intervals:
        if out and start <= out[-1][1]:
            out[-1][1] = max(out[-1][1], end)
        else:
            out.append([start, end])
    return out
```

Sorting costs O(n log n) and the sweep is O(n). The same "sort by start, keep the running end" sweep solves *insert interval*, *meeting rooms* and *minimum arrows to burst balloons* with small changes.

## Before you code

- Decide what "overlap" means for touching endpoints. Here, touching intervals merge.
- Ask yourself what happens when one interval sits entirely inside another.

## Check your approach

Solve it on LeetCode first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

A buggy version replaces `max(out[-1][1], end)` with plain `end`. **What does it return for `[[1, 10], [2, 3], [4, 5]]`?**
