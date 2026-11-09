# What Is a Dockerfile?

So far you've only used images other people made. Now you'll make your own.

An image is built from a **Dockerfile**: a plain text file with a list of steps, one per line. It's the recipe for your image.

```dockerfile
FROM alpine:3.20
CMD ["echo", "Hello from my own image!"]
```

That's a complete Dockerfile. It says:

1. **`FROM`**: start from the `alpine:3.20` image. Every Dockerfile starts from an existing image.
2. **`CMD`**: when a container starts, run `echo "Hello from my own image!"`.

Then two commands turn it into something that runs:

```
Dockerfile  --docker build-->  image  --docker run-->  container
 (recipe)                    (packed box)            (running box)
```

**In what order do these things happen?**
