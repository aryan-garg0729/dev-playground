# How ENTRYPOINT and CMD Combine

When a container starts, Docker builds its command line as **ENTRYPOINT + CMD**:

| ENTRYPOINT | CMD | Container runs |
| --- | --- | --- |
| *(none)* | `["brewd"]` | `brewd` |
| `["brewd"]` | *(none)* | `brewd` |
| `["/usr/local/bin/entrypoint.sh"]` | `["--version"]` | `entrypoint.sh --version` |

And `docker run` can override each part:

- `docker run IMAGE arg1 arg2` **replaces CMD** with `arg1 arg2`
- `docker run --entrypoint prog IMAGE` **replaces ENTRYPOINT** (and clears CMD)

Priya's image has:

```dockerfile
ENTRYPOINT ["brewd"]
CMD ["--version"]
```

**What does `docker run --rm IMAGE --help` run?**
