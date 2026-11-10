# Hamiltonian Flights

**Practice:** [Hamiltonian Flights on CSES](https://cses.fi/problemset/task/1690)

## The task

Count the routes from city 1 to city n that visit every city exactly once (n ≤ 20), modulo 10⁹+7.

## Before you code

- `dp[mask][v]` = number of ways to visit exactly the cities in `mask`, starting at city 1 and ending at v.
- That's 2²⁰ · 20 ≈ 2·10⁷ states, so C++ is a must, and memory layout matters.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Which pruning is valid and roughly halves the work?**
