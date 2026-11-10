# Towers

**Practice:** [Towers on CSES](https://cses.fi/problemset/task/1073)

## The task

Cubes arrive in a fixed order. Each cube goes on top of an existing tower whose top cube is **strictly bigger**, or starts a new tower. Minimize the number of towers.

## Before you code

- Keep the tower tops in a sorted list. Which tower should a cube go on, so the others stay as useful as possible?
- This is the patience-sorting idea behind the O(n log n) longest increasing subsequence.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

The sorted tower tops are `[3, 5, 5, 9]` and the next cube has size 5. **Where does it go?**
