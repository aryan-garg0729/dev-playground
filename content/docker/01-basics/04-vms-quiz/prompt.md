# Containers vs Virtual Machines

Before Docker, the usual way to give an app its own clean machine was a **virtual machine** (VM): a pretend computer running inside your real one.

A VM works, but it's heavy. Each one boots a **full operating system**: its own Linux or Windows, its own memory, gigabytes of disk. Starting one takes a while, and your laptop can only run a few.

A container skips all that. It **shares the operating system your computer is already running** and only brings the app and its files.

| | Virtual machine | Container |
| --- | --- | --- |
| Has its own operating system | Yes | No, it shares yours |
| Starts in | Seconds to minutes | About a second |
| Size | Gigabytes | Often megabytes |
| How many fit on a laptop | A few | Dozens |

Think of a VM as a **house**: its own foundation, plumbing and wiring. A container is an **apartment**: its own locked rooms, but the building's plumbing and wiring are shared.

(On a Mac, Docker quietly runs one small Linux VM and puts all your containers inside it. You never have to manage it.)

**Which is true about containers?**
