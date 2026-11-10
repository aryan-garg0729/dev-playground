# Flight Routes

**Practice:** [Flight Routes on CSES](https://cses.fi/problemset/task/1196)

## The task

Print the prices of the k cheapest routes from city 1 to city n (k ≤ 10; routes may revisit cities).

## Before you code

- Dijkstra where each city may be popped up to k times. Each pop of city n is the next route price.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why is it correct to stop expanding a city once it has been popped k times?**
