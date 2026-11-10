# Planets Cycles

**Practice:** [Planets Cycles on CSES](https://cses.fi/problemset/task/1751)

## The task

Every planet has exactly one teleporter to some planet. For every starting planet, count the teleports you make until you reach a planet you've already visited.

## Pattern recap: functional graphs

A graph where every node has exactly one outgoing edge is a **functional graph**. It's a set of cycles with trees hanging off them, all pointing in towards the cycles. For each node, the answer is "steps to reach its cycle + the cycle's length".

Find each cycle by walking until you meet a node from the current walk, then fill in the answers backwards along the walk.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

A planet is 3 teleports away from a cycle of length 5. **What's its answer?**
