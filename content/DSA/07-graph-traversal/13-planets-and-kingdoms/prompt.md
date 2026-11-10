# Planets and Kingdoms

**Practice:** [Planets and Kingdoms on CSES](https://cses.fi/problemset/task/1683)

## The task

Group planets into kingdoms: two planets share a kingdom when each can reach the other through teleporters. Print the number of kingdoms and a kingdom label for every planet.

## Before you code

- The kingdoms are exactly the SCCs. Use Kosaraju, written iteratively.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

In Kosaraju's second pass, **why take the nodes in decreasing finish time?**
