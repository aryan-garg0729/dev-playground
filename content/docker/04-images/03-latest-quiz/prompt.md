# The latest Trap

If you don't type a tag, Docker quietly adds `:latest`:

```sh
docker pull nginx        # really: nginx:latest
```

`latest` sounds like "the newest version". **It isn't.** It's just a tag name, like any other:

- The image's publisher decides what `latest` points at, and moves it whenever they like.
- Your computer keeps the `latest` it downloaded months ago until you pull again.
- Your teammate pulls today and gets a different `latest`.

That's the "works on my machine" problem coming back through the side door, which is exactly what Docker was meant to fix.

The rule: **always pick a real version** (`nginx:1.27-alpine`, not `nginx`) for anything that matters. This course does that everywhere.

**Which statement about the `latest` tag is true?**
