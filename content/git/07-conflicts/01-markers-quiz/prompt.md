# Reading Conflict Markers

Git merges by itself when two sides changed **different** parts of a file. When both changed the **same lines**, it can't know which you want, so it stops and writes both versions into the file:

```
# Contact

Email: hello@example.com
<<<<<<< HEAD
Office hours: Tuesdays, 10:00
=======
Office hours: Mondays, 09:30, online
>>>>>>> update-hours
```

This is a **conflict**. It sounds scary, but it's normal and it happens to everyone. The markers are plain text that Git put in your file. Nothing is broken: Git is waiting for you to decide.

You ran `git merge update-hours` while on `main`.

**Which side is which?**
