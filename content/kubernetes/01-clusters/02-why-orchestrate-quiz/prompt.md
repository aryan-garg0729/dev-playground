# What Docker Alone Can't Do

Mo asks the obvious question: *"We already have Docker. Why do we need another thing?"*

Docker is a **container platform**: it builds images and runs containers on one machine. Kubernetes is a **container orchestrator**: it decides which containers should run on which of many machines, and keeps it that way.

The gap shows up as soon as you have more than one server:

- **Self-healing.** Containers are short-lived. They crash, run out of memory, or vanish with their host. Someone has to notice and start a replacement.
- **Scaling.** Black Moon needs 10 copies of brewd at 9am and 2 by noon.
- **One address for many copies.** Clients need one stable name, not a list of IPs that keeps changing.
- **Rollouts.** Ship a new version without dropping orders, and take it back when it's bad.

You tell Kubernetes *what* you want ("3 copies of brewd, reachable as `brewd`"), and it works out *how*, continuously.

**Which of these needs an orchestrator, not just Docker?**
