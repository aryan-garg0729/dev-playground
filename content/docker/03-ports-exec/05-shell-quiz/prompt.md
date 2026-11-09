# Leaving the Shell

Quick check on what you just did.

You ran `docker exec -it cq-web sh`, looked around, and typed `exit`.

Remember the rule from the last chapter: **a container lives as long as its main command.** For `cq-web`, the main command is nginx. `ps` showed it as PID 1, the first process in the box.

**What happens to `cq-web` when you type `exit` in the exec shell?**
