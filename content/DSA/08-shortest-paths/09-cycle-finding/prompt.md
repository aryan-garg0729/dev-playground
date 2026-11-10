# Cycle Finding

**Practice:** [Cycle Finding on CSES](https://cses.fi/problemset/task/1197)

## The task

Find any negative cycle in a directed graph and print it, or print NO.

## Before you code

- Bellman–Ford with every distance starting at 0, n rounds, then the parent walk from "Bellman–Ford".

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why start every `dist` at 0 instead of `dist[1] = 0` and the rest at infinity?**
