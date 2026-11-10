# Counting Grids

**Practice:** [Counting Grids on CSES](https://cses.fi/problemset/task/2210)

## The task

Count the n × n grids of black and white squares, where rotations of a grid count as the same grid, modulo 10⁹+7.

## Before you code

- Four symmetries: rotations by 0°, 90°, 180° and 270°.
- Count the cycles each one makes on the n² cells. With odd n, the centre cell is a cycle on its own.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

**How many cycles does a 90° rotation create on an n × n grid?**
