# Flight Routes Check

**Practice:** [Flight Routes Check on CSES](https://cses.fi/problemset/task/1682)

## The task

Check whether every city can reach every other city by one-way flights. If not, print a pair (a, b) with no route from a to b.

## Before you code

- The two-BFS check from "Strongly Connected Components". The BFS that fails tells you the pair.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

The BFS from city 1 on the **reversed** graph doesn't reach city 4. **Which pair do you print?**
