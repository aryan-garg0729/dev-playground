# Round Trip II

**Practice:** [Round Trip II on CSES](https://cses.fi/problemset/task/1678)

## The task

Same as Round Trip, but the flights are one-way: find any directed cycle.

## Before you code

- Three colours, as in "Cycle Detection". When you reach a grey node, the cycle is on your DFS path.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

The DFS path is `1 → 4 → 2 → 5`, and city 5 has a flight to city 4 (grey). **Which cycle do you print?**
