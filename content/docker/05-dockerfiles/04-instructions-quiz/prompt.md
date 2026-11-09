# The Main Instructions

Most Dockerfiles use only a handful of instructions:

| Instruction | What it does | Happens when |
| --- | --- | --- |
| `FROM image` | Start from an existing image. Always first. | Build |
| `WORKDIR /path` | Go to this folder inside the image (and create it). Later steps work from here. | Build |
| `COPY src dest` | Copy files from your folder into the image | Build |
| `RUN command` | Run a command and save its result into the image (install things, compile code) | Build |
| `ENV NAME=value` | Set an environment variable for containers | Build (used at run) |
| `EXPOSE 8080` | Note which port the app listens on (a label only, it doesn't publish anything) | Build |
| `CMD ["program", "arg"]` | The default command when a container starts | **Run** |

The big one to get right: **`RUN` happens once, at build time. `CMD` happens every time a container starts.**

**You want `curl` installed in your image, so that every container already has it. Which instruction installs it?**
