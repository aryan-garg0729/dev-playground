# Flight Discount

**Practice:** [Flight Discount on CSES](https://cses.fi/problemset/task/1195)

## The task

Find the cheapest route from city 1 to city n when you may halve the price of one flight on the way (rounded down).

## Before you code

- State graph: `dist[city][used]` with used ∈ {0, 1}. From state 0 you can fly at full price (stay at 0) or half price (move to 1).

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

An alternative without a state graph: run Dijkstra from city 1, and from city n on the reversed graph. **The answer is then the minimum, over every flight u → v with price w, of…**
