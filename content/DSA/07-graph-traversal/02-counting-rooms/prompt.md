# Counting Rooms

**Practice:** [Counting Rooms on CSES](https://cses.fi/problemset/task/1192)

## The task

A map is a grid of floor (`.`) and wall (`#`) squares. Count the rooms: connected areas of floor, moving up, down, left and right.

## Before you code

- Loop over every cell. Each unvisited floor cell starts a new flood fill and a new room.
- Write the flood fill iteratively, for the reason in "DFS and BFS".

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

The map can be 1000 × 1000. **In Python, what's the main risk of a recursive flood fill?**
