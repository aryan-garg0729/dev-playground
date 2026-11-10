# Traffic Lights

**Practice:** [Traffic Lights on CSES](https://cses.fi/problemset/task/1163)

## The task

A street runs from 0 to x. Traffic lights are added one at a time at given positions. After each addition, print the length of the longest stretch of street with no light.

## Before you code

- Online approach: a set of light positions plus a multiset of gap lengths. Each new light splits one gap into two.
- There's also an offline trick that needs neither structure. The question is about it.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Offline trick: add all the lights, then remove them in reverse order and record the answers backwards. **Why is that easier?**
