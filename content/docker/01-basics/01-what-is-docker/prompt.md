# What Is Docker?

You've probably heard this one:

> "But it works on my machine!"

An app almost never runs alone. It needs a certain version of Python or Node, some libraries, a few settings and maybe a config file. Your laptop has the right mix. Your teammate's laptop has a slightly different mix. The server has a third one. Something breaks, and everyone blames everyone.

**Docker fixes this by putting the app in a box.** The box holds the app *and everything it needs to run*. You hand the box to someone else, and it runs exactly the same on their machine, on a server or in the cloud.

That box is called a **container**.

Think of real shipping containers. Before them, every ship was loaded by hand with sacks, barrels and crates of every shape. Then the world agreed on one standard steel box. Cranes, ships and trucks only need to handle *the box*, not whatever is inside. Docker does the same thing for software.

**Why do teams use it?**

- **Same everywhere**: no more "works on my machine".
- **Quick setup**: a new developer runs one command instead of following a 20-step install guide.
- **Clean**: an app's mess stays inside its box. Delete the box and your laptop is clean again.

In this course you'll run real Docker on your own computer, from your first container to shipping your own images.

**What does Docker do?**
