# DSA Interview Prep: course plan

Source: `content/DSA/sheet.md` (the author's practice sheet). Every problem on the sheet is one lesson.

## Audience

Learners who already know the basics (arrays, recursion, big-O, BFS/DFS) and are revising before interviews.
This course is the exception to "build from zero": lessons recap a pattern quickly instead of teaching it from scratch.

## Decisions (2026-11-10, at the author's request)

- **Practice happens on the original sites.** Each problem lesson links to CSES, AtCoder, LeetCode or SPOJ.
  Problem statements are never copied: each lesson gives a one-line paraphrase, the pattern and a nudge.
  CodeQuest grades no code here.
- **Every lesson is an `mcq`.** That keeps the course gradable on the current production setup, which has no grading worker.
  A problem lesson ends with a question about the key idea ("check your approach"); its `explain` gives the full approach, complexity and traps.
- **No bosses, no locks.** Chapters lock only behind an undefeated boss, so a course without bosses is fully open, in any order.
  Every lesson is `free: true`, so the paywall never applies either.
- XP: pattern lessons 5 (quiz), problem lessons 10. `difficulty` reflects the problem's real difficulty (1-5).
- Every template is written twice, C++ then Python, as back-to-back fenced blocks that the app shows as tabs (C++ first).
  C++ snippets assume `#include <bits/stdc++.h>` and `using namespace std;`, and avoid GCC-only builtins like `__lg`.
  Lessons say where Python is too slow for CSES limits. `add_cpp.py` was the one-off that added the C++ twins.

## Lesson shapes

- **Pattern lesson** (`<slug>`, title = pattern name): when to use it, a template, complexity, traps, then one recognition or trace question.
- **Problem lesson** (title = the site's problem name): practice link (+ reading link from the sheet), "The task" (paraphrase),
  optional "Pattern recap" (when the sheet section has only one problem), "Before you code" nudges, then the question.
  A sheet section with 2+ problems gets its own pattern lesson first.

The generator in `.drafts/gen/` (`python3 content/DSA/.drafts/gen/gen.py content/DSA`) writes every `exercise.yaml` and `prompt.md`, and rotates the correct option's position.

## Chapters

| # | Chapter id | Title | Sheet sections |
| --- | --- | --- | --- |
| 1 | dsa.classics | Interview Classics | Interview Classics |
| 2 | dsa.greedy | Searching, Sorting & Greedy | Binary search, sorting + greedy, prefix sums, priority queue / multiset |
| 3 | dsa.stack-window | Stacks & Sliding Windows | Monotonic stack, sliding window |
| 4 | dsa.bits | Bit Manipulation | Bit manipulation |
| 5 | dsa.math | Math & Counting | Algebra, counting / Burnside / inclusion-exclusion |
| 6 | dsa.range | Range Queries | Sparse table, segment / Fenwick trees |
| 7 | dsa.graphs | Graph Traversal | DFS/BFS, topological sort, cycle detection, functional graph, SCC |
| 8 | dsa.paths | Shortest Paths | Dijkstra, Floyd-Warshall, Bellman-Ford |
| 9 | dsa.dsu | Disjoint Set Union & MST | DSU, MST |
| 10 | dsa.trees | Trees | Diameter, binary lifting / LCA, tree DP, Euler tour, small-to-large |
| 11 | dsa.dp | Dynamic Programming | AtCoder Educational DP Contest picks |
| 12 | dsa.strings | Strings | Trie, hashing, KMP, palindromes, suffix array, string DP, minimal rotation |

## Sheet corrections

- "Sliding Window Minimum" linked CSES 1644, which is *Maximum Subarray Sum II*. The lesson uses CSES 3221 (Sliding Window Minimum).
- CSES 2428 is titled *Distinct Values Subarrays II*; 1647 is *Static Range Minimum Queries*.
