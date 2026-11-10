# Required Substring

**Practice:** [Required Substring on CSES](https://cses.fi/problemset/task/1112)

## The task

Count the strings of length n over the letters A–Z that contain a given pattern of length m, modulo 10⁹+7.

## Before you code

- Count the strings that *don't* contain it: a DP over (length so far, KMP state = how many pattern characters are matched), never letting the state reach m.
- The answer is 26ⁿ minus that count.

## Check your approach

Solve it on CSES first, or give it a real try (about 30 minutes) before reading on. The explanation after you answer walks through the full approach.

From KMP state k (k characters matched), you add the character c. **Which state do you move to?**
