from h import P, X

CSES = "https://cses.fi/problemset/task/"
CPA = "https://cp-algorithms.web.app/"

classics = dict(slug="classics", id="dsa.classics", title="Interview Classics", lessons=[
X("lru-cache", "LRU Cache", "https://leetcode.com/problems/lru-cache/", 3, ["design", "hashing", "linked-list"], r"""
## The task

Design a cache with a fixed capacity that supports `get(key)` and `put(key, value)` in O(1) average time. When it's full, `put` evicts the **least recently used** key. A `get` counts as a use.

## Pattern recap: hash map + doubly linked list

- The list keeps keys in order of use: most recent at the head, least recent at the tail.
- The map points from each key to its list node, so you can reach any node in O(1).
- On every use, unlink the node and move it to the head. When the cache overflows, drop the tail.
- Sentinel (dummy) head and tail nodes remove every "is the list empty?" special case.

```cpp
struct Node {
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
}
```

```python
class Node:
    __slots__ = ("key", "val", "prev", "next")

    def __init__(self, key=0, val=0):
        self.key, self.val = key, val
        self.prev = self.next = None


def unlink(node):
    node.prev.next = node.next
    node.next.prev = node.prev


def push_front(head, node):
    node.prev, node.next = head, head.next
    head.next.prev = node
    head.next = node
```

In Python, `collections.OrderedDict` with `move_to_end(key)` and `popitem(last=False)` does the same job. Interviewers often ask for the hand-built version anyway, so practise both.

## Before you code

- Write `get` and `put` in terms of the two helpers and the code stays short.
- `put` on an existing key updates the value **and** counts as a use.
- Store the key inside the node: when you evict the tail, you need it to delete the map entry.
""", "**Why does the list need to be *doubly* linked?**", [
    "So you can unlink any node in O(1) when the map hands you a pointer to it",
    "So you can walk backwards from the tail to find the least recently used key",
    "Because a singly linked list can't hold both a key and a value",
    "It doesn't: a singly linked list plus the map is also O(1) per operation",
], "To move a node to the front you must unlink it, and unlinking means rewiring its predecessor's `next`. A singly linked node doesn't know its predecessor, so you would have to walk from the head to find it: O(n). The tail is reachable in O(1) through the sentinel either way, so walking backwards isn't the reason. With `prev` pointers, `unlink` is two assignments, which is what makes both `get` and `put` O(1)."),

X("merge-intervals", "Merge Intervals", "https://leetcode.com/problems/merge-intervals/", 2, ["sorting", "intervals"], r"""
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
""", "A buggy version replaces `max(out[-1][1], end)` with plain `end`. **What does it return for `[[1, 10], [2, 3], [4, 5]]`?**", [
    "[[1, 3], [4, 5]]",
    "[[1, 10]]",
    "[[1, 5]]",
    "[[1, 10], [4, 5]]",
], "`[1, 10]` starts the output. `[2, 3]` overlaps (2 ≤ 10), but the bug sets the end to 3, shrinking the interval to `[1, 3]`. Then `[4, 5]` no longer overlaps (4 > 3), so it's appended: `[[1, 3], [4, 5]]`. The right answer is `[[1, 10]]`, because both later intervals sit inside `[1, 10]`. Contained intervals are the case `max` exists for."),

X("count-inversions", "Count Inversions", "https://www.spoj.com/problems/INVCNT/", 3, ["divide-and-conquer", "sorting"], r"""
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
""", "In the middle of a merge, `left = [3, 5, 8]` with `i` pointing at `5` (the `3` is already merged), and `right[j]` is `4`. The merge takes the `4`. **How many inversions does that step add?**", [
    "2", "1", "3", "0",
], "`4` is smaller than every element still waiting in `left`, which is `5` and `8`: `len(left) - i = 3 - 1 = 2`. The `3` was merged earlier because it's smaller than `4`, so it doesn't form an inversion with it. Adding `len(left)` (3) instead of `len(left) - i` is the classic slip."),
])

greedy = dict(slug="search-sort-greedy", id="dsa.greedy", title="Searching, Sorting & Greedy", lessons=[
X("factory-machines", "Factory Machines", CSES + "1620", 2, ["binary-search"], r"""
## The task

n machines work at the same time; machine i needs `k[i]` seconds per product. Find the least time to make t products in total.

## Pattern recap: binary search on the answer

Use it when you can't compute the answer directly but you **can** check a guess quickly, and the check is monotone: once a guess works, every bigger guess works too.

```cpp
// smallest x in [lo, hi] with ok(x), assuming ok(hi) is true
long long first_true(long long lo, long long hi, function<bool(long long)> ok) {
    while (lo < hi) {
        long long mid = lo + (hi - lo) / 2;
        if (ok(mid)) hi = mid;
        else lo = mid + 1;
    }
    return lo;
}
```

```python
def first_true(lo, hi, ok):
    # smallest x in [lo, hi] with ok(x), assuming ok(hi) is True
    while lo < hi:
        mid = (lo + hi) // 2
        if ok(mid):
            hi = mid
        else:
            lo = mid + 1
    return lo
```

Signals: "minimum time / maximum size such that…", "minimize the largest…" (split array largest sum, ship packages within D days, Koko eating bananas).

## Before you code

- Write `ok(T)`: how many products do the machines finish in T seconds?
- That count can be huge. With fixed-width integers, stop adding once it reaches t.
""", "**Which `hi` is guaranteed to satisfy `ok(hi)` for every valid input?**", [
    "min(k) * t", "t", "max(k)", "sum(k)",
], "The fastest machine alone makes t products in `min(k) * t` seconds, so that much time is always enough. `t` fails when every machine needs more than 1 second per product, and `max(k)` or `sum(k)` fail when t is large. With k and t up to 10⁹, `hi` can be 10¹⁸. That still fits in a 64-bit integer, but the running sum inside `ok` can overflow unless you stop early."),

P("greedy-by-sorting", "Greedy by Sorting", 2, ["greedy", "sorting"], r"""
Many greedy problems are "sort by the right key, then take what you can in one pass". The hard part is choosing the key and convincing yourself (and the interviewer) that it's right.

## The exchange argument

Take any optimal answer that disagrees with your greedy choice. Show that swapping in the greedy choice makes it no worse. Then greedy is optimal too.

Example: to fit the most non-overlapping intervals, take the interval that **ends first**. Any optimal answer starts with some interval that ends at or after yours. Swap it for yours and nothing else breaks, because yours frees the timeline sooner.

## Common keys

| Goal | Sort by |
| --- | --- |
| Most non-overlapping intervals | end time |
| Least total (or average) finishing time | duration, shortest first |
| Merge or cover intervals | start time |
| Pair up items (boats, assignments) | sort, then two pointers from both ends |

## When greedy fails

If a small counterexample breaks your rule, switch to DP. Coin change with coins `{1, 3, 4}` and target 6 is the classic one: greedy takes `4 + 1 + 1`, but `3 + 3` uses fewer coins.
""", "**You want the largest set of non-overlapping intervals. Which greedy rule is correct?**", [
    "Pick the interval that ends first",
    "Pick the shortest interval first",
    "Pick the interval that starts first",
    "Pick the interval that overlaps the fewest others first",
], "Ending first leaves the most room for everything after it, which is what the exchange argument proves. Shortest-first fails on `[1, 5], [4, 7], [6, 10]`: it takes `[4, 7]`, which blocks both others, while `[1, 5], [6, 10]` gives 2. Earliest start fails when one long interval starts first. Fewest-overlaps-first sounds clever but has known counterexamples too."),

X("movie-festival", "Movie Festival", CSES + "1629", 1, ["greedy", "intervals"], r"""
## The task

n movies have start and end times. You watch whole movies, one at a time. Maximize the number of movies you watch.

## Before you code

- Which key from "Greedy by Sorting" applies?
- Decide what happens when one movie ends exactly when the next starts. The CSES sample answers it.
""", "Movies are sorted by end time, and `last_end` is when the last chosen movie ends. **When do you take the next movie?**", [
    "start >= last_end", "start > last_end", "end >= last_end", "start <= last_end",
], "A movie that starts exactly when the previous one ends fits: the sample takes 3–5 and then 5–8. With `>` you would skip those movies and undercount. The other two options test the wrong thing, since after sorting every later end is already ≥ `last_end`. Total O(n log n) for the sort."),

X("tasks-and-deadlines", "Tasks and Deadlines", CSES + "1630", 2, ["greedy", "sorting"], r"""
## The task

n tasks each have a duration a and a deadline d. You do them one after another, starting at time 0. Finishing a task at time f earns `d − f` (which can be negative). You must do every task. Maximize the total reward.

## Before you code

- Write the total reward as a formula. Which part depends on the order, and which doesn't?
""", "**Which order maximizes the total reward?**", [
    "By duration, shortest first",
    "By deadline, earliest first",
    "By d − a, largest first",
    "By duration, longest first",
], "Total reward = (sum of deadlines) − (sum of finish times). The sum of deadlines is the same in every order, so only the finish times matter, and running the shortest tasks first minimizes their sum: each duration counts once for its own task and once for every task after it. Sorting by deadline is the tempting trap; deadlines don't affect the best order at all here. Use 64-bit integers: the total is large and can be negative."),

X("reading-books", "Reading Books", CSES + "1631", 2, ["greedy", "math"], r"""
## The task

Two readers must each read all n books. Book i takes `t[i]` minutes for either of them, and nobody can read a book while the other one has it. Minimize the time until both are done.

## Before you code

- Try a few small cases by hand. When does one long book force someone to wait?
""", "Books take `[2, 8, 3]` minutes. **What's the minimum total time?**", [
    "16", "13", "26", "11",
], "The answer is `max(sum, 2 * max)`. Usually both readers stay busy the whole time and finish together after `sum` minutes, by reading in a staggered order. But if one book is longer than all the others combined, it's the bottleneck: the readers take turns with it, so it alone takes `2 * max`. Here the 8-minute book is longer than 2 + 3 = 5, so the answer is 2 · 8 = 16, not 13."),

X("towers", "Towers", CSES + "1073", 3, ["greedy", "binary-search"], r"""
## The task

Cubes arrive in a fixed order. Each cube goes on top of an existing tower whose top cube is **strictly bigger**, or starts a new tower. Minimize the number of towers.

## Before you code

- Keep the tower tops in a sorted list. Which tower should a cube go on, so the others stay as useful as possible?
- This is the patience-sorting idea behind the O(n log n) longest increasing subsequence.
""", "The sorted tower tops are `[3, 5, 5, 9]` and the next cube has size 5. **Where does it go?**", [
    "On the tower with top 9 (bisect_right)",
    "On a tower with top 5 (bisect_left)",
    "It starts a new tower",
    "On the tower with top 3",
], "A cube needs a strictly bigger top, so the 5s are out and 3 is too small. The smallest top bigger than 5 is 9, which `bisect_right(tops, 5)` finds (index 3), and replacing it keeps the list sorted: `[3, 5, 5, 5]`. `bisect_left` would land on a 5 and break the strict rule. Each cube costs one binary search, O(n log n) total. Using the smallest bigger top keeps the large tops free for large cubes."),

P("prefix-sums", "Prefix Sums", 1, ["prefix-sum", "hashing"], r"""
Need many range sums? Precompute once:

```cpp
vector<long long> pre(a.size() + 1, 0);
for (size_t i = 0; i < a.size(); i++) pre[i + 1] = pre[i] + a[i];
// sum of a[l..r-1] (l inclusive, r exclusive)
long long total = pre[r] - pre[l];
```

```python
pre = [0]
for x in a:
    pre.append(pre[-1] + x)
# sum of a[l:r] (l inclusive, r exclusive)
total = pre[r] - pre[l]
```

## Counting subarrays with a given sum

`a[l:r]` sums to x exactly when `pre[r] - pre[l] == x`. Walk r from left to right and count the earlier prefixes equal to `pre[r] - x`:

```cpp
long long count_sum(const vector<long long>& a, long long x) {
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
}
```

```python
from collections import defaultdict


def count_sum(a, x):
    seen = defaultdict(int)
    seen[0] = 1
    cur = ans = 0
    for v in a:
        cur += v
        ans += seen[cur - x]
        seen[cur] += 1
    return ans
```

O(n), and it works with negative numbers. Variants: "subarray sum divisible by k" (store `cur % k`), "longest subarray with sum x" (store the first index of each prefix), "equal number of 0s and 1s" (map 0 to −1 first).

## Two pointers or a hash map?

A sliding window only works when every value is non-negative, because then growing the window can only raise the sum. With negatives, use the hash map.

## In 2D

`pre[i][j]` is the sum of the top-left i × j rectangle. Any rectangle is four lookups: `pre[r2][c2] - pre[r1][c2] - pre[r2][c1] + pre[r1][c1]`.
""", "The array has negative numbers. You need to count the subarrays that sum to x in O(n) or O(n log n). **What works?**", [
    "Prefix sums plus a hash map of prefix counts",
    "Two pointers (a sliding window)",
    "Sort the array, then two pointers",
    "Binary search over the prefix sums",
], "Two pointers rely on the sum only growing as the window grows, which negatives break. Binary search over prefix sums needs them sorted, which again only holds without negatives. Sorting the array destroys which elements are next to each other, so subarrays lose their meaning. The hash map counts, for each right end, how many earlier prefixes equal `cur - x`, and it doesn't care about signs."),

X("subarray-sums-ii", "Subarray Sums II", CSES + "1661", 2, ["prefix-sum", "hashing"], r"""
## The task

Count the subarrays whose sum is exactly x. Values can be negative.

## Before you code

- This is the template from "Prefix Sums".
- In C++, `unordered_map` with the default hash can be slowed to a crawl by crafted tests (anti-hash tests). Use `map`, or a randomized custom hash.
""", "In the template, **what does `seen[0] = 1` account for?**", [
    "Subarrays that start at index 0",
    "An empty input array",
    "Zeros inside the array",
    "Nothing; it only avoids a KeyError",
], "A subarray that starts at index 0 has sum `cur` itself, and it matches when `cur - x == 0`: the empty prefix before the first element. Without `seen[0] = 1`, those subarrays are never counted. For `a = [5]` and `x = 5` you'd get 0 instead of 1. With `defaultdict` a missing key already reads as 0, so the line isn't about KeyErrors."),

X("maximum-subarray-sum", "Maximum Subarray Sum", CSES + "1643", 1, ["dp", "kadane"], r"""
## The task

Find the largest sum of a non-empty contiguous subarray.

## Pattern recap: Kadane's algorithm

For each position, keep the best sum of a subarray that **ends there**. Either extend the previous one or start fresh:

```cpp
long long best = a[0], cur = a[0];
for (size_t i = 1; i < a.size(); i++) {
    cur = max(a[i], cur + a[i]);      // a is a vector<long long>
    best = max(best, cur);
}
```

```python
best = cur = a[0]
for x in a[1:]:
    cur = max(x, cur + x)
    best = max(best, cur)
```

It's a one-variable DP. Variants: maximum product subarray (track both min and max), maximum circular subarray (total minus the minimum subarray).
""", "A common version starts with `best = cur = 0` and uses `cur = max(0, cur + x)`. **What does it return for `[-3, -1, -2]`?**", [
    "0", "-1", "-6", "-3",
], "`cur` is clamped to 0 at every step and `best` never drops below its starting 0, so the result is 0: the sum of an *empty* subarray. The task wants a non-empty one, so the right answer is −1. Start from `a[0]` (or negative infinity) and use `max(x, cur + x)`. Any test where every value is negative catches this."),

P("heaps-and-ordered-sets", "Heaps and Ordered Sets", 2, ["heap", "ordered-set"], r"""
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
""", "You must repeatedly find the largest value ≤ x **and remove it**, over 2·10⁵ operations. **Which structure makes every operation O(log n)?**", [
    "A balanced BST / ordered multiset",
    "A min-heap",
    "A hash map from value to count",
    "A sorted Python list with bisect",
], "An ordered multiset (C++ `multiset`, Java `TreeMap`) answers 'largest ≤ x' with `upper_bound` and one step back, and deletes in O(log n). A heap only knows its top element. A hash map has no order at all. A sorted list finds the value in O(log n), but deleting it shifts every later element, which is O(n). That often passes because the shift is fast, but it isn't O(log n)."),

X("concert-tickets", "Concert Tickets", CSES + "1091", 3, ["ordered-set", "greedy"], r"""
## The task

n tickets have prices. m customers arrive in order, each with a maximum price. Each customer gets the most expensive ticket they can afford, and that ticket is gone. Print what each customer pays, or −1.

## Before you code

- An ordered multiset solves it directly.
- In Python, try sorting the prices and using a DSU that jumps from a sold ticket to the nearest unsold one on its left.
""", "Prices are `[3, 5, 5, 7, 8]`, and three customers arrive, each willing to pay up to 6. **What gets printed?**", [
    "5 5 3", "5 5 -1", "5 3 -1", "7 7 7",
], "Each customer takes the largest unsold price ≤ 6: first a 5, then the other 5, then the 3. Only a fourth such customer would get −1. With a C++ multiset: `it = upper_bound(x)`; if `it` is `begin()`, print −1; otherwise step back once, print it and erase *that iterator*. `erase(value)` would delete every copy of the price."),

X("traffic-lights", "Traffic Lights", CSES + "1163", 3, ["ordered-set", "offline"], r"""
## The task

A street runs from 0 to x. Traffic lights are added one at a time at given positions. After each addition, print the length of the longest stretch of street with no light.

## Before you code

- Online approach: a set of light positions plus a multiset of gap lengths. Each new light splits one gap into two.
- There's also an offline trick that needs neither structure. The question is about it.
""", "Offline trick: add all the lights, then remove them in reverse order and record the answers backwards. **Why is that easier?**", [
    "Removing a light merges two gaps, so the longest gap can only grow and one running max is enough",
    "Reversing the order sorts the positions",
    "The last light added is always in the middle of the street",
    "It isn't easier; it's the same work as the online version",
], "Adding lights only shrinks gaps, so the longest one can disappear and you need a multiset to find the next longest. Removing lights only merges gaps, so the answer never goes down: track `max(best, merged_gap)`. You still need each light's neighbours, and a doubly linked list over the sorted positions gives them in O(1) per removal. Total: one sort, then O(n), which is easy in Python."),
])

stack_window = dict(slug="stacks-windows", id="dsa.stack-window", title="Stacks & Sliding Windows", lessons=[
P("monotonic-stack", "Monotonic Stack", 2, ["stack"], r"""
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
""", "The loop has a `while` inside a `for`. **What's the total running time?**", [
    "O(n), because each index is pushed once and popped at most once",
    "O(n²) in the worst case",
    "O(n log n)",
    "O(n · h), where h is the largest stack height",
], "The inner loop's total work over the whole run is bounded by the number of pushes, which is n. One `while` can pop many items, but those items are gone for good. This amortized argument is what interviewers want to hear when they point at the nested loop."),

X("nearest-smaller-values", "Nearest Smaller Values", CSES + "1645", 2, ["stack"], r"""
## The task

For each position, find the nearest position to its left with a **strictly** smaller value. Print positions 1-based, or 0 if there's none.

## Before you code

- It's the template from "Monotonic Stack", with 1-based output.
""", "**Which pop condition is right for \"strictly smaller\"?**", [
    "a[st[-1]] >= x", "a[st[-1]] > x", "a[st[-1]] < x", "a[st[-1]] <= x",
], "An equal value isn't strictly smaller, so it can't be x's answer, and nothing to the right needs it either, because x is just as small and closer. Pop it. With `>` and `a = [2, 2]`, the second 2 would report position 1. The last two options keep a decreasing stack, which finds nearest *greater* values instead."),

X("advertisement", "Advertisement", CSES + "1142", 3, ["stack"], r"""
## The task

A fence is made of boards with given heights. Find the largest rectangle that fits inside it: the "largest rectangle in a histogram" problem (also LeetCode 84).

## Before you code

- For each board, the widest rectangle of exactly its height reaches from just after its previous smaller board to just before its next smaller board.
- Two monotonic-stack passes (or one careful pass) give both neighbours.
""", "Heights are `[2, 1, 5, 6, 2, 3]`. **What's the largest area?**", [
    "10", "12", "6", "8",
], "Height 5 spans the boards of height 5 and 6 (its smaller neighbours are the 1 on the left and the 2 on the right): 5 × 2 = 10. The second 2 spans the four boards 5, 6, 2, 3 for 8, and height 1 spans all six for 6. 12 would need a width-2 rectangle of height 6, but the 6's neighbours are shorter. Area = `h[i] * (right[i] - left[i] - 1)`, where left and right are the indices of the previous and next smaller boards. Use 64-bit."),

P("sliding-window", "Sliding Window", 2, ["two-pointers", "deque"], r"""
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
""", "**Which problem can't be solved with a plain two-pointer window?**", [
    "Shortest subarray with sum ≥ s, where values may be negative",
    "Longest substring without repeating characters",
    "Shortest subarray with sum ≥ s, where all values are positive",
    "Longest subarray with at most k zeros",
], "With negatives, shrinking the window can *raise* the sum and growing it can lower it, so the 'shrink while bad' rule is wrong. That version (LeetCode 862) uses prefix sums with a monotonic deque instead. The other three keep the monotone property: adding elements can only add repeats, sum or zeros."),

X("sliding-window-minimum", "Sliding Window Minimum", CSES + "3221", 2, ["deque"], r"""
## The task

Find the minimum of every window of k consecutive elements. The input can be very large, so read the statement for how it's given and what exactly to print.

## Before you code

- The monotonic deque from "Sliding Window".
- A heap with lazy deletion also works, but it's O(n log n). On a big input you want the O(n) deque and fast input reading.
""", "In `window_mins`, **why does the deque pop from the back while `a[dq[-1]] >= x`, before appending x?**", [
    "Those values are no smaller than x and leave the window before it, so they can never be a window's minimum again",
    "To keep the deque at most k long",
    "To keep the deque sorted by index",
    "It's only an optimization; the answers would be the same without it",
], "While x is in the window, no earlier value ≥ x can be the minimum, and those values leave the window *before* x does, so they're useless from now on. Popping keeps the deque's values increasing, so the front is always the minimum. The length is controlled by the `popleft` check, and indices are sorted anyway because you only append. Without the pops the front isn't the minimum, so the answers would be wrong."),

X("distinct-values-subarrays-ii", "Distinct Values Subarrays II", CSES + "2428", 3, ["two-pointers"], r"""
## The task

Count the subarrays that contain at most k distinct values.

## Before you code

- Two pointers with a count map, like the window template.
- For each right end r, which starts give a valid subarray?
""", "`[l, r]` is the longest valid window ending at r. **How many valid subarrays end at r?**", [
    "r - l + 1", "1", "(r - l + 1) * (r - l + 2) / 2", "k",
], "Every shorter subarray ending at r (starting at l + 1, l + 2, …, r) has at most as many distinct values, so it's valid too: r − l + 1 subarrays in all. Adding the triangle number at every r counts the same subarrays many times. The total can reach about n²/2, so use 64-bit. Both pointers only move forward: O(n)."),

X("sliding-window-median", "Sliding Window Median", CSES + "1076", 4, ["heap", "ordered-set"], r"""
## The task

Print the median of each window of k elements. For even k, use the lower of the two middle values.

## Before you code

- Two heaps with lazy deletion, or an ordered multiset (C++: keep an iterator at the median as elements come and go).
- Keep the sizes balanced: the max-heap `low` holds ⌈k/2⌉ live elements, and its top is the answer.
""", "With two heaps (max-heap `low`, min-heap `high`), the element leaving the window is marked for lazy deletion. **What must you update right away, even though the element is still physically inside a heap?**", [
    "The live-size counter of the heap it belongs to, so balancing uses the real sizes",
    "Nothing; the clean-up happens when it reaches the top",
    "The other heap's top",
    "Both heaps, rebuilt with heapify",
], "Balancing compares sizes. If a dead element still counts towards `low`, the heaps look balanced when they aren't, and the wrong element ends up as the median. Keep a live-size counter per heap: decrement it when you mark an element dead, and work out which heap it's in by comparing it with `low`'s top. Pop dead elements whenever they surface at a top. Each step is O(log n) amortized, while a heapify rebuild would cost O(k) per window."),
])

bits = dict(slug="bits", id="dsa.bits", title="Bit Manipulation", lessons=[
P("bit-tricks", "Bit Tricks", 2, ["bits"], r"""
The facts behind most interview bit problems:

| Fact | Why it helps |
| --- | --- |
| `x ^ x == 0`, `x ^ 0 == x` | pairs cancel out |
| XOR is commutative and associative | order doesn't matter |
| `x & -x` | the lowest set bit of x |
| `x & (x - 1)` | x with its lowest set bit cleared |
| `(x >> i) & 1` | bit i of x |
| prefix XOR `p[r] ^ p[l]` | XOR of `a[l:r]`, just like prefix sums |

```cpp
int single(const vector<int>& nums) {   // every value appears twice except one
    int acc = 0;
    for (int x : nums) acc ^= x;
    return acc;
}
```

```python
def single(nums):          # every value appears twice except one
    acc = 0
    for x in nums:
        acc ^= x
    return acc
```

**Bitwise trie.** To maximize `x ^ y` over a set of y values, store them in a binary trie from the highest bit down. For each query, walk to the child with the *opposite* bit whenever it exists. That's O(bits) per query.

Python ints never overflow and behave as if they had infinitely many sign bits, so add a mask such as `& 0xFFFFFFFF` when a problem expects 32-bit wraparound.

Precedence trap: in C, C++ and Java, `==` binds tighter than `&`, so `x & 1 == 0` means `x & (1 == 0)`. Always add parentheses.
""", "**For a positive x, what does `(x & (x - 1)) == 0` test?**", [
    "x is a power of two", "x is even", "x is odd", "x has every bit set",
], "`x - 1` flips the lowest set bit and every zero below it, so `x & (x - 1)` clears exactly the lowest set bit. The result is 0 only when that was the *only* set bit, which means x is a power of two. Evenness is `(x & 1) == 0`. 'Every bit set' is `(x & (x + 1)) == 0`."),

X("single-number-iii", "Single Number III", "https://leetcode.com/problems/single-number-iii/description/", 3, ["bits", "xor"], r"""
## The task

Every number appears exactly twice except two, which appear once. Find those two in O(n) time and O(1) extra space.

## Before you code

- XOR everything first. What do you know about the result?
- How could one of its bits split the numbers into two groups?
""", "After XORing everything you have `s = a ^ b` with `a ≠ b`. **Why is `s & -s` useful?**", [
    "It's a bit where a and b differ, so it splits all numbers into two groups with one of a, b in each",
    "It equals a",
    "It's a bit that a and b have in common",
    "It counts how many numbers there are",
], "`s` is non-zero because a ≠ b, and each of its set bits is a position where a and b differ. `s & -s` isolates the lowest one. XOR the numbers that have that bit into one accumulator and the rest into another. Each pair lands in the same group and cancels, so one accumulator ends as a and the other as b. In C++ with `int`, `-s` overflows when `s` is `INT_MIN`, so use an unsigned or 64-bit type."),

X("maximum-xor-subarray", "Maximum Xor Subarray", CSES + "1655", 4, ["bits", "trie"], r"""
## The task

Find the largest XOR of a contiguous subarray.

## Before you code

- Turn subarrays into pairs of prefix XORs (`p[r] ^ p[l]`), then maximize over pairs with a bitwise trie.
- Insert `p[0] = 0` before anything else.
""", "**Why must the trie start with the prefix value 0 inserted?**", [
    "So subarrays that start at index 0 are counted",
    "So the trie isn't empty when you query, otherwise it crashes",
    "Because 0 XOR anything is 0",
    "It's only needed when the array has negative values",
], "The XOR of `a[0:r]` is `p[r] ^ p[0]`, with `p[0] = 0`. Without the 0, those subarrays, including the whole array, are missed: for `[7]` you'd report nothing (or 0) instead of 7. It's the same idea as `seen[0] = 1` in prefix-sum counting. Each insert and query walks about 30 bits, so the total is O(30 · n)."),
])

math = dict(slug="math", id="dsa.math", title="Math & Counting", lessons=[
X("exponentiation-ii", "Exponentiation II", CSES + "1712", 3, ["math", "modular"], r"""
## The task

Compute a^(b^c) modulo 10⁹+7 for many queries.

## Pattern recap: fast power and Fermat

Binary exponentiation squares the base and halves the exponent:

```cpp
long long power(long long a, long long e, long long m) {
    long long r = 1;
    a %= m;
    while (e) {
        if (e & 1) r = r * a % m;
        a = a * a % m;
        e >>= 1;
    }
    return r;
}
```

```python
def power(a, e, m):
    r = 1
    a %= m
    while e:
        if e & 1:
            r = r * a % m
        a = a * a % m
        e >>= 1
    return r
```

That's O(log e), and Python's built-in `pow(a, e, m)` does exactly this.

**Fermat's little theorem**: for a prime p and an a that p doesn't divide, `a^(p−1) ≡ 1 (mod p)`. So exponents can be reduced mod p − 1, and `a^(p−2)` is the modular inverse of a.
""", "**How do you reduce the exponent `b^c` before using it?**", [
    "Compute b^c mod (10⁹+6)",
    "Compute b^c mod (10⁹+7)",
    "Compute b^c exactly, then reduce it",
    "Reduce b mod (10⁹+7), then raise it to c",
], "By Fermat, `a^(p−1) ≡ 1 (mod p)` for p = 10⁹+7 whenever p doesn't divide a, so exponents repeat with period p − 1 = 10⁹+6. Reducing mod p is the classic wrong answer. `b^c` itself is astronomically large. The one exception is when a is a multiple of p: then the answer is 0 for any positive exponent, so handle it if the limits allow it."),

X("binomial-coefficients", "Binomial Coefficients", CSES + "1079", 3, ["math", "combinatorics"], r"""
## The task

Answer many queries of C(a, b) modulo 10⁹+7, with a and b up to 10⁶.

## Pattern recap: factorials and inverse factorials

Precompute once, answer each query in O(1):

```cpp
const long long MOD = 1e9 + 7;
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
}
```

```python
MOD = 10**9 + 7
N = 10**6
fact = [1] * (N + 1)
for i in range(1, N + 1):
    fact[i] = fact[i - 1] * i % MOD
inv = [1] * (N + 1)
inv[N] = pow(fact[N], MOD - 2, MOD)
for i in range(N, 0, -1):
    inv[i - 1] = inv[i] * i % MOD


def C(a, b):
    if b < 0 or b > a:
        return 0
    return fact[a] * inv[b] % MOD * inv[a - b] % MOD
```

Division doesn't exist mod p; multiplying by an inverse is how you divide.
""", "The code raises a value to the power `MOD - 2` only once. **Why is `inv[i - 1] = inv[i] * i` correct?**", [
    "Because 1/(i−1)! = i · (1/i!)",
    "Because every inverse factorial is the same",
    "Because MOD is odd",
    "It's an approximation that's only exact for small i",
], "`(i−1)! = i! / i`, so `1/(i−1)! = i · (1/i!)`. Starting from the single inverse computed with Fermat, each step down costs one multiplication: O(N + log MOD) in total instead of O(N log MOD). It relies on MOD being prime, so every non-zero value has an inverse."),

X("next-prime", "Next Prime", CSES + "3396", 3, ["math", "primes"], r"""
## The task

For each given number n, print the smallest prime greater than n. The numbers are large, so check the limits in the statement.

## Pattern recap: primality tools

- **Trial division** up to √n: O(√n) per number.
- **Sieve of Eratosthenes**: every prime up to N in O(N log log N).

```cpp
vector<bool> sieve(int n) {
    vector<bool> is_p(n + 1, true);
    is_p[0] = is_p[1] = false;
    for (long long i = 2; i * i <= n; i++)
        if (is_p[i])
            for (long long j = i * i; j <= n; j += i) is_p[j] = false;
    return is_p;
}
```

```python
def sieve(n):
    is_p = bytearray([1]) * (n + 1)
    is_p[0:2] = b"\x00\x00"
    for i in range(2, int(n ** 0.5) + 1):
        if is_p[i]:
            is_p[i * i :: i] = bytes(len(range(i * i, n + 1, i)))
    return is_p
```

- **Segmented sieve**: primes inside a window `[L, R]` with a large L, sieved using the primes up to √R.
- **Miller–Rabin**: a fast probabilistic test; with a fixed set of bases it's exact for every 64-bit integer.

Gaps between consecutive primes are small: below 10¹⁸ they never exceed about 1,500.
""", "**Why does it matter that the gaps between primes are small?**", [
    "You only test a handful of candidates, so each query costs a few primality tests",
    "You can skip testing and add a fixed offset instead",
    "It proves every odd number in the gap is prime",
    "It means the sieve can stop at n",
], "The next prime is never far away, so try n + 1, n + 2, … and each candidate costs one primality test (skipping even numbers halves the work). With trial division each test is O(√n), up to 10⁶ steps for n near 10¹². With Miller–Rabin it's a few modular powers, which is what keeps many queries fast."),

P("burnside", "Counting Under Symmetry", 3, ["combinatorics"], r"""
How many necklaces of n beads in m colours are there, if rotations of the same necklace count once? Counting all colourings and dividing by n is wrong, because some colourings look the same under several rotations.

**Burnside's lemma**: the number of distinct objects is the *average*, over all symmetries, of the number of colourings that the symmetry leaves unchanged.

For necklace rotations: rotating by k positions splits the beads into `gcd(k, n)` cycles, and a colouring stays the same exactly when each cycle is a single colour. So:

```text
distinct = (1 / n) · Σ_{k=0}^{n-1} m^gcd(k, n)
```

Modulo a prime, the `1 / n` is a multiplication by `pow(n, MOD - 2, MOD)`.

Recipe for any symmetry problem:

1. List the symmetries (n rotations; for a square grid, 4 rotations).
2. For each one, count the cycles it makes on the cells.
3. Add up `colours ^ cycles` and divide by the number of symmetries.

Reading: [Burnside's lemma on cp-algorithms](https://cp-algorithms.web.app/combinatorics/burnside.html) and the [video from your sheet](https://youtu.be/D0d9bYZ_qDY).
""", "A necklace has 4 beads in 2 colours, and rotations count as the same necklace. **How many distinct necklaces are there?**", [
    "6", "4", "8", "16",
], "Rotations by 0, 1, 2 and 3 positions make 4, 1, 2 and 1 cycles, so they fix 2⁴ = 16, 2¹ = 2, 2² = 4 and 2¹ = 2 colourings. The average is (16 + 2 + 4 + 2) / 4 = 6: 0000, 0001, 0011, 0101, 0111, 1111. The naive 16 / 4 = 4 undercounts, because 0000 and 0101 look the same under more than one rotation."),

X("counting-necklaces", "Counting Necklaces", CSES + "2209", 3, ["combinatorics"], r"""
## The task

Count the necklaces of n beads in m colours, where rotations count as the same necklace, modulo 10⁹+7.

## Before you code

- The formula from "Counting Under Symmetry". Loop k from 0 to n − 1 with fast power.
""", "For n = 6, **how many cycles does rotating by k = 4 positions create?**", [
    "2", "1", "3", "4",
], "Rotating by 4 sends bead 0 → 4 → 2 → 0 and bead 1 → 5 → 3 → 1: two cycles of length 3, and gcd(4, 6) = 2. Each cycle must be a single colour, so this rotation fixes m² colourings. Summing m^gcd(k, n) is O(n log n). You can group the k values by their gcd using Euler's totient, but the limits here don't need it."),

X("counting-grids", "Counting Grids", CSES + "2210", 4, ["combinatorics"], r"""
## The task

Count the n × n grids of black and white squares, where rotations of a grid count as the same grid, modulo 10⁹+7.

## Before you code

- Four symmetries: rotations by 0°, 90°, 180° and 270°.
- Count the cycles each one makes on the n² cells. With odd n, the centre cell is a cycle on its own.
""", "**How many cycles does a 90° rotation create on an n × n grid?**", [
    "n²/4 for even n, (n² + 3)/4 for odd n",
    "n²/2",
    "n",
    "n²/4 for every n",
], "Away from the centre, every cell is on a cycle of 4 (its four rotated positions). With even n that's n²/4 cycles. With odd n the centre cell stays put, giving (n² − 1)/4 cycles of 4 plus 1, which is (n² + 3)/4. A 180° rotation pairs cells instead: ⌈n²/2⌉ cycles. The answer is (2^(n²) + 2 · 2^c90 + 2^c180) / 4, with the division done by a modular inverse."),

X("prime-multiples", "Prime Multiples", CSES + "2185", 3, ["combinatorics", "inclusion-exclusion"], r"""
## The task

Given k distinct primes (k ≤ 20) and n up to 10¹⁸, count the numbers from 1 to n that are divisible by at least one of the primes.

## Pattern recap: inclusion–exclusion

|A ∪ B ∪ C| = |A| + |B| + |C| − |A ∩ B| − |A ∩ C| − |B ∩ C| + |A ∩ B ∩ C|.

Over every non-empty subset of the primes, add `n // product` for odd-sized subsets and subtract it for even-sized ones. That's 2^k subsets, about 10⁶ for k = 20.
""", "Products of primes can exceed 10¹⁸. **What's the right fix?**", [
    "Stop building a subset's product once it exceeds n; that subset contributes 0",
    "Compute the product in floating point",
    "Take the product modulo 10⁹+7",
    "Ignore it; overflow only causes small errors",
], "Once the product passes n, `n // product` is 0, so that subset adds nothing. In C++, check before multiplying (`p > n / cur`) so you never overflow, which is undefined behaviour for signed integers. A modular product makes the division meaningless, and doubles lose precision long before 10¹⁸. Python's integers can't overflow, but capping still saves time."),
])

rangeq = dict(slug="range-queries", id="dsa.range", title="Range Queries", lessons=[
X("static-range-minimum-queries", "Static Range Minimum Queries", CSES + "1647", 2, ["sparse-table"], r"""
## The task

An array never changes. Answer many queries "minimum of `a[l..r]`".

## Pattern recap: sparse table

`st[j][i]` is the minimum of the 2^j elements starting at i. Build it in O(n log n), then answer any query with two overlapping blocks in O(1), because min doesn't mind counting an element twice.

```cpp
vector<vector<int>> build(const vector<int>& a) {
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
    int j = 31 - __builtin_clz(r - l + 1);   // floor(log2(length))
    return min(st[j][l], st[j][r - (1 << j) + 1]);
}
```

```python
def build(a):
    st = [a[:]]
    j = 1
    while (1 << j) <= len(a):
        prev, half = st[-1], 1 << (j - 1)
        st.append([min(prev[i], prev[i + half]) for i in range(len(a) - (1 << j) + 1)])
        j += 1
    return st


def query(st, l, r):            # inclusive l..r
    j = (r - l + 1).bit_length() - 1
    return min(st[j][l], st[j][r - (1 << j) + 1])
```

Also used for range max, range gcd, and LCA through an Euler tour.
""", "**Why can't a sparse table answer range *sum* queries in O(1) the same way?**", [
    "The two blocks overlap, and a sum would count the overlap twice; min and max don't care",
    "Sums are too large to store",
    "The table would need O(n²) memory",
    "It can; sums work exactly like min",
], "The O(1) query covers [l, r] with two blocks of length 2^j that usually overlap. That's fine for *idempotent* operations, where f(x, x) = x: min, max, gcd, AND, OR. A sum would count the overlap twice. For static sums use prefix sums, and when the array changes, a Fenwick tree.",
    read=("Sparse Table on cp-algorithms", CPA + "data_structures/sparse-table.html")),

P("fenwick-and-segment-trees", "Fenwick and Segment Trees", 3, ["fenwick", "segment-tree"], r"""
When the array **changes** between queries, prefix sums and sparse tables break. Two structures give O(log n) updates and queries.

**Fenwick tree (BIT)**: short and fast; handles prefix sums (and other invertible operations).

```cpp
struct BIT {
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
};
```

```python
class BIT:
    def __init__(self, n):
        self.t = [0] * (n + 1)          # 1-indexed

    def add(self, i, d):
        while i < len(self.t):
            self.t[i] += d
            i += i & -i

    def prefix(self, i):                # sum of positions 1..i
        s = 0
        while i > 0:
            s += self.t[i]
            i -= i & -i
        return s
```

A range sum is `prefix(r) - prefix(l - 1)`. For **range update, point query**, store differences: `add(l, d)` and `add(r + 1, -d)`, then the value at i is `prefix(i)`.

**Segment tree**: each node stores the answer for its segment. It handles any associative operation (min, max, sum, gcd), finds "the first index with value ≥ x" by descending from the root, and supports **range updates with lazy propagation**: a pending update stored at a node and pushed to its children only when needed.

| Need | Use |
| --- | --- |
| point update, range sum | Fenwick |
| range add, point query | Fenwick on differences |
| range add, range sum | two Fenwicks, or a lazy segment tree |
| point update, range min/max | segment tree |
| first position with value ≥ x | segment tree descent |
| 2D grid, point update, rectangle sum | 2D Fenwick |

Reading: [segment trees](https://cp-algorithms.web.app/data_structures/segment_tree.html) and [Fenwick trees](https://cp-algorithms.web.app/data_structures/fenwick.html) on cp-algorithms.
""", "**Why is a Fenwick tree a poor fit for \"point update, range *minimum*\"?**", [
    "A range minimum isn't a difference of two prefix answers: you can't 'subtract' a prefix minimum",
    "Fenwick trees can only store integers",
    "Its updates are O(n)",
    "It's a perfect fit; range minimum is the classic Fenwick use",
], "A Fenwick tree answers prefixes. For sums, `sum(l..r) = prefix(r) − prefix(l − 1)`, but there's no way to remove `min(1..l−1)` from `min(1..r)`. Prefix minimums alone are possible when values only decrease, but arbitrary ranges need a segment tree (or a sparse table if the array never changes)."),

X("range-updates-and-sums", "Range Updates and Sums", CSES + "1735", 4, ["segment-tree", "lazy"], r"""
## The task

An array supports three operations: add x to every value in a range, set every value in a range to x, and report the sum of a range.

## Before you code

- A lazy segment tree with two kinds of pending tag: "set to v" and "add d".
- Decide how two tags combine before you write any code.
""", "A node already has a pending \"add 5\" and now receives \"set to 2\". **What's its pending tag afterwards?**", [
    "Set to 2 (the earlier add is discarded)",
    "Set to 7",
    "Add 5 (the set is ignored while an add is pending)",
    "Add −3",
], "An assignment overwrites everything before it, including pending adds, so the tag becomes 'set 2'. The other order does combine: 'set 2' then 'add 5' becomes 'set 7', and 'add 5' then 'add 3' becomes 'add 8'. Getting this composition right, and pushing it to both children before descending, is the whole problem. Range sums reach about 2·10⁵ · 10⁶, so use 64-bit.", read=("Segment Tree on cp-algorithms", CPA + "data_structures/segment_tree.html")),

X("hotel-queries", "Hotel Queries", CSES + "1143", 3, ["segment-tree"], r"""
## The task

Hotels stand in a row, each with some free rooms. Groups arrive in order; each needs r rooms in a single hotel and takes the **first** (leftmost) hotel with enough free rooms. Print that hotel, or 0 if none fits.

## Before you code

- A max segment tree over the free rooms.
- Descend from the root: go left whenever the left child's maximum is big enough.
""", "**Why descend the tree instead of binary searching on the index with a range-max query at each step?**", [
    "Descending finds the leftmost index in one O(log n) walk; binary search with range queries costs O(log² n)",
    "Binary search gives the wrong hotel",
    "A segment tree can't answer prefix maximums",
    "There's no difference",
], "Binary search with a range-max query per step is correct, but every step is itself O(log n), so a group costs O(log² n). Walking down from the root, taking the left child whenever its max is ≥ r, finds the same hotel in O(log n). Then subtract r at that leaf and update the path back up. This 'descend to the first position' trick comes up often."),

X("forest-queries-ii", "Forest Queries II", CSES + "1739", 3, ["fenwick", "2d"], r"""
## The task

An n × n grid has trees and empty squares. Either flip one square's state, or count the trees inside a rectangle.

## Before you code

- A 2D Fenwick tree: the 1D loops, nested. O(log² n) per operation.
- A rectangle query uses the same four-term inclusion–exclusion as 2D prefix sums.
""", "**What does `add(x, y, d)` look like in a 2D Fenwick tree?**", [
    "Two nested loops: x moves by x & -x, and for each x, y moves by y & -y",
    "One loop over x, then a separate loop over y",
    "A loop over every cell in row x",
    "A segment tree per row, updated in O(n)",
], "Each row-level Fenwick node owns a whole Fenwick tree over the columns, so the update walks both: `i = x; while i <= n: j = y; while j <= n: t[i][j] += d; j += j & -j` and then `i += i & -i`. That's O(log² n). Two separate loops would update a cross shape instead of the right set of nodes. A flip is `add(+1)` or `add(-1)` depending on the current state, so keep the grid itself too."),
])

CHAPTERS = [classics, greedy, stack_window, bits, math, rangeq]
