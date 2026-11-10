# Shortest Routes II

**Practice:** [Shortest Routes II on CSES](https://cses.fi/problemset/task/1672)

## The task

Up to 500 cities with two-way roads of given lengths, and many queries: the shortest distance from a to b, or −1 if there's no route.

## Before you code

- Floyd–Warshall once; then each query is a lookup.
- At n = 500 that's 1.25·10⁸ inner steps, too slow for plain Python. Use C++ for this one.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

The input has two roads between cities 2 and 3, of lengths 7 and 4. **What should `dist[2][3]` start at?**
