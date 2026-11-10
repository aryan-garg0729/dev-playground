# J - Sushi

**Practice:** [J - Sushi on AtCoder](https://atcoder.jp/contests/dp/tasks/dp_j)

## The task

n dishes each start with 1 to 3 pieces of sushi. Each step, roll a die that picks a dish uniformly at random, and eat one piece from it if there are any left. Find the expected number of rolls until everything is eaten: expected-value DP.

## Before you code

- Dishes with the same number of pieces are interchangeable, so the state is just how many dishes have 1, 2 and 3 pieces: `E[a][b][c]`.
- A roll that picks an empty dish leads back to the same state. Solve that equation for E.

## Check your approach

Solve it on AtCoder first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

From state (a, b, c) with z = n − a − b − c empty dishes, `E = 1 + (z/n)·E + (a/n)·E[a−1][b][c] + …`. **How do you handle the E on the right-hand side?**
