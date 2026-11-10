# High Score

**Practice:** [High Score on CSES](https://cses.fi/problemset/task/1673)

## The task

Find the maximum score of a path from room 1 to room n, where each tunnel adds its score (possibly negative). If the score can be made arbitrarily large, print −1.

## Before you code

- Negate the scores and run Bellman–Ford from room 1.
- Not every negative cycle matters: only one that's reachable from room 1 **and** can reach room n.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

In round n, room 7 improves again, but room n can't be reached from room 7. **What does that mean for the answer?**
