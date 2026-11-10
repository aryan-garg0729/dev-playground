# Monsters

**Practice:** [Monsters on CSES](https://cses.fi/problemset/task/1194)

## The task

A grid has walls, monsters and you. Every step, you and every monster move one square. Reach the border of the grid without a monster ever being able to catch you. Print a path, or say it's impossible.

## Before you code

- First, a multi-source BFS from **all** monsters gives each cell the earliest time any monster can be there.
- Then a BFS from your position, only stepping onto cells you reach strictly before every monster.
- Keep parent pointers so you can print the path.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**Why is it enough to plan against the monsters' *earliest* arrival times, instead of simulating how they move?**
