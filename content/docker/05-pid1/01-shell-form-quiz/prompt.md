# Shell Form, Exec Form

`CMD` and `ENTRYPOINT` each come in two forms:

```dockerfile
CMD ["brewd", "--version"]     # exec form (JSON array): run this program with these arguments
CMD brewd --version            # shell form: Docker stores ["/bin/sh", "-c", "brewd --version"]
```

Exec form runs your program **directly**. Shell form starts `/bin/sh` and hands it the whole line as a string. That buys you `&&`, `$VARS` and pipes, and costs you control over **which process is PID 1**.

PID 1 matters because it's the only process `docker stop` sends SIGTERM to. Whatever is PID 1 has to shut the rest down.

Mo wants a startup message and writes:

```dockerfile
CMD echo "starting brewd" && brewd
```

**In the running container, which process is PID 1?**
