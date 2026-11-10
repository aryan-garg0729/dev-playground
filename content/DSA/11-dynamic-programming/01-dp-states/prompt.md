# Designing a DP State

Every DP answers three questions:

1. **State**: what do you need to know about a partial solution to finish it? (`dp[i][w]`: best value from the first i items with total weight w.)
2. **Transition**: how does a state follow from smaller ones? (take item i, or skip it)
3. **Base and order**: which states are trivial, and in what order do you fill in the rest?

Then check the cost: states × transitions per state. If that's too big, change the state, not the code.

**Swap state and value.** If `dp[i][w] = best value` has too many weights but the values are small, use `dp[i][v] = least weight that reaches value v`. The next problem is exactly this.

**Common shapes in interviews**

| Shape | Examples |
| --- | --- |
| prefix `dp[i]` | climbing stairs, house robber, word break, decode ways |
| two sequences `dp[i][j]` | LCS, edit distance, regex matching |
| interval `dp[l][r]` | burst balloons, matrix chain, merging stones |
| subset `dp[mask]` | travelling salesman, assignments with n ≤ 20 |
| tree `dp[v][state]` | house robber III, tree matching |
| digits `dp[pos][tight][...]` | count numbers ≤ N with a property |

Practise recovering the solution, not just its value: store the choice made at each state, or walk back through the table comparing values.

The rest of this chapter follows the AtCoder Educational DP Contest problems on your sheet ([contest page](https://atcoder.jp/contests/dp/tasks)).

## Check yourself

A knapsack has n = 100 items, capacity W = 10⁹ and every value ≤ 1000. **Which state is feasible?**
