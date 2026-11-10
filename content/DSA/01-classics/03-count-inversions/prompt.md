# Count Inversions

**Practice:** [Count Inversions on SPOJ](https://www.spoj.com/problems/INVCNT/)

## The task

Count the pairs `i < j` with `a[i] > a[j]` (an **inversion**: a pair that's out of order). With n around 2·10⁵, the O(n²) double loop is too slow.

## Pattern recap: count while you merge sort

Merge sort already compares elements across the two halves. When the merge takes an element from the **right** half, every element still waiting in the left half is bigger than it, and each of those pairs is an inversion:

```cpp
// sorts a[lo, hi) and returns its inversion count
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
}
```

```python
def sort_count(a):
    if len(a) <= 1:
        return a, 0
    mid = len(a) // 2
    left, x = sort_count(a[:mid])
    right, y = sort_count(a[mid:])
    merged, inv, i, j = [], x + y, 0, 0
    while i < len(left) and j < len(right):
        if left[i] <= right[j]:
            merged.append(left[i])
            i += 1
        else:
            merged.append(right[j])
            j += 1
            inv += len(left) - i
    merged += left[i:] + right[j:]
    return merged, inv
```

That's O(n log n). The other standard answer is a Fenwick tree (chapter "Range Queries"): walk left to right and, for each value, count how many earlier values are bigger.

## Before you code

- The count can reach about n²/2 ≈ 2·10¹⁰. In C++ or Java, that needs a 64-bit integer.
- Equal values are **not** inversions. Check which side the merge takes on a tie.

## Check your approach

Solve it on SPOJ first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

In the middle of a merge, `left = [3, 5, 8]` with `i` pointing at `5` (the `3` is already merged), and `right[j]` is `4`. The merge takes the `4`. **How many inversions does that step add?**
