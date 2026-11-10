# Next Prime

**Practice:** [Next Prime on CSES](https://cses.fi/problemset/task/3396)

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

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why does it matter that the gaps between primes are small?**
