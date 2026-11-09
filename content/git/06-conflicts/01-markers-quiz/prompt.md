# Reading Conflict Markers

Git merges silently when two sides changed **different** parts of a file. When both changed the **same lines**, it can't know which you want, so it stops and writes both versions into the file:

```
# On-call

Primary pager: +1-555-0142
Escalation window: 5 minutes
<<<<<<< HEAD
Handover: Tuesdays, 10:00 (moved from Monday)
=======
Handover: Mondays, 09:30, in #ops
>>>>>>> feat/handover
```

This is a **conflict**. The markers are plain text that Git put in your file. It's not broken; it's waiting for you to decide.

You ran `git merge feat/handover` while on `main`.

**Which side is which?**
