# Movie Festival

**Practice:** [Movie Festival on CSES](https://cses.fi/problemset/task/1629)

## The task

n movies have start and end times. You watch whole movies, one at a time. Maximize the number of movies you watch.

## Before you code

- Which key from "Greedy by Sorting" applies?
- Decide what happens when one movie ends exactly when the next starts. The CSES sample answers it.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

Movies are sorted by end time, and `last_end` is when the last chosen movie ends. **When do you take the next movie?**
