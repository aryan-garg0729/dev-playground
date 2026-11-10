# Shortest Routes I

**Practice:** [Shortest Routes I on CSES](https://cses.fi/problemset/task/1671)

## The task

In a directed graph with weighted flights, print the shortest distance from city 1 to every city.

## Before you code

- The template, with fast input reading.
- Distances can reach about 10¹⁴. In C++ or Java use 64-bit integers, and make "infinity" bigger than any real path.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

A C++ solution uses `int dist[]` and `INF = 1e9`. Weights go up to 10⁹, with up to 10⁵ cities. **What goes wrong?**
