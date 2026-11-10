# Round Trip

**Practice:** [Round Trip on CSES](https://cses.fi/problemset/task/1669)

## The task

In an undirected road network, find any round trip (a cycle through at least 3 cities) and print it, or say none exists.

## Before you code

- DFS with parent tracking. When you see a visited neighbour that isn't your parent, walk the parent pointers back to it.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why does "visited and not my parent" guarantee a cycle of at least 3 cities here?**
