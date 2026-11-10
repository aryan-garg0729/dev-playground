# LRU Cache

**Practice:** [LRU Cache on LeetCode](https://leetcode.com/problems/lru-cache/)

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

## Check your approach

Solve it on LeetCode first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why does the list need to be *doubly* linked?**
