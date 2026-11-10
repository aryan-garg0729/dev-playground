# Exponentiation II

**Practice:** [Exponentiation II on CSES](https://cses.fi/problemset/task/1712)

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

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**How do you reduce the exponent `b^c` before using it?**
