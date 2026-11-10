# Bit Tricks

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

## Check yourself

**For a positive x, what does `(x & (x - 1)) == 0` test?**
