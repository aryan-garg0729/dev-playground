# Binomial Coefficients

**Practice:** [Binomial Coefficients on CSES](https://cses.fi/problemset/task/1079)

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

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

The code raises a value to the power `MOD - 2` only once. **Why is `inv[i - 1] = inv[i] * i` correct?**
