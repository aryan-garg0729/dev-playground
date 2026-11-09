# CMD vs ENTRYPOINT

You've been typing commands after the image name since chapter 1:

```sh
docker run --rm alpine:3.20 echo hi
```

That works because Alpine's default command is set with **`CMD`**, and **anything you type after the image name replaces `CMD`**.

There's a second instruction, **`ENTRYPOINT`**, that sets a program that **always** runs. With an entrypoint, the words after the image name don't replace the program. They're **added to it** as arguments.

What a container runs is simply:

```
ENTRYPOINT  +  CMD
(the program)  (default arguments, replaced by whatever you type)
```

| Dockerfile | `docker run IMG` runs | `docker run IMG --version` runs |
| --- | --- | --- |
| `CMD ["webapp"]` | `webapp` | `--version` (fails: not a program) |
| `ENTRYPOINT ["webapp"]` | `webapp` | `webapp --version` |

Use `CMD` alone for general-purpose images (like alpine, where you run all sorts of things). Use `ENTRYPOINT` when the image *is* one program, so `docker run IMG --help` feels like running the program itself.

An image has:

```dockerfile
ENTRYPOINT ["ls"]
CMD ["-l", "/tmp"]
```

**What does `docker run IMG --version` run?**
