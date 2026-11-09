# Read the Exit Code

A quick recap of the last two lessons:

1. `0` means success.
2. `1`-`125` mean the program exited on purpose and chose the number.
3. Above 128 means a signal ended it, and the code is **128 + the signal number**.

| Signal | Number | Usually sent by |
| --- | --- | --- |
| SIGINT | 2 | Ctrl-C in a terminal |
| SIGKILL | 9 | `docker kill`, the out-of-memory killer, or `docker stop` once its timeout runs out |
| SIGTERM | 15 | `docker stop` (first attempt) |

Dex gets paged at 2 a.m.:

```
eb-brewer   Exited (130) 4 minutes ago
```

**What ended eb-brewer?**
