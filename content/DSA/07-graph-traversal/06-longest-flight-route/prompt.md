# Longest Flight Route

**Practice:** [Longest Flight Route on CSES](https://cses.fi/problemset/task/1680)

## The task

In a directed flight network with no cycles, find the route from city 1 to city n that visits the most cities, and print it.

## Before you code

- DP in topological order: `best[v] = max(best[u] + 1)` over flights u → v, with routes starting only at city 1.
- Store each city's parent to rebuild the route.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Some cities can't be reached from city 1. **How should `best` start?**
