# Brackets or No Brackets?

`CMD` and `ENTRYPOINT` can be written two ways:

```dockerfile
# exec form: a list in brackets
CMD ["webapp", "--port", "8080"]

# shell form: plain text
CMD webapp --port 8080
```

They look the same, but they aren't:

- **Exec form** (brackets): Docker runs `webapp` directly. It is the container's main process.
- **Shell form** (no brackets): Docker runs a **shell** and hands it your text, as `/bin/sh -c "webapp --port 8080"`. The shell is now the main process, and `webapp` runs underneath it.

Shell form lets you use shell tricks like `$HOME` or `&&`. But a shell in the middle causes trouble when it's time to stop, which is the next lesson.

Rule of thumb: **use the brackets** for `CMD` and `ENTRYPOINT`.

An image is built with `CMD webapp` (no brackets).

**What does `docker image inspect -f '{{json .Config.Cmd}}'` show?**
