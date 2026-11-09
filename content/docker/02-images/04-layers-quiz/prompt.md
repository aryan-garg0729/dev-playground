# Which Steps Add Layers?

In the last lesson, `docker history` showed many more rows than `RootFS.Layers` had entries. Every build step adds a history row, but only some add a **filesystem layer**.

You'll write Dockerfiles in the next chapter. For now, each line below is one build step.

**Which of these steps create a filesystem layer?** Pick every one that does.
