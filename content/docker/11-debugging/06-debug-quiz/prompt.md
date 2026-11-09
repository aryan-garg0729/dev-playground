# Which Tool When?

Your debugging toolbox so far:

| Symptom | Reach for |
| --- | --- |
| Container stopped or keeps exiting | `docker ps -a`, then `docker logs` |
| Wrong settings, ports, mounts or network? | `docker inspect` |
| Computer slow, fan loud | `docker stats`, `docker top` |
| "What did it write?" | `docker diff`, `docker cp` |
| Need to look around inside, test from its point of view | `docker exec -it NAME sh` |

A container is **running**, and its logs look normal. But the app says it can't read its config file and can't reach the database. You want to check, from **inside the container**, whether the file exists and whether the database's name resolves.

**Which tool do you use?**
