# Concert Tickets

**Practice:** [Concert Tickets on CSES](https://cses.fi/problemset/task/1091)

## The task

n tickets have prices. m customers arrive in order, each with a maximum price. Each customer gets the most expensive ticket they can afford, and that ticket is gone. Print what each customer pays, or −1.

## Before you code

- An ordered multiset solves it directly.
- In Python, try sorting the prices and using a DSU that jumps from a sold ticket to the nearest unsold one on its left.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Prices are `[3, 5, 5, 7, 8]`, and three customers arrive, each willing to pay up to 6. **What gets printed?**
