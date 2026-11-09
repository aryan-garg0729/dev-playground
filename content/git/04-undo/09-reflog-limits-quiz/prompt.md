# What the Reflog Can't Save

The reflog makes Git feel indestructible, and mostly it is. It has limits though:

- It only records moves of `HEAD` and branches, which means **commits**.
- It lives in `.git/logs/` of **one clone**. It's never pushed or shared.
- Entries expire after a few weeks to months, after which `git gc` can delete the commits for real.

**Which of these can the reflog NOT bring back? Pick all that apply.**
