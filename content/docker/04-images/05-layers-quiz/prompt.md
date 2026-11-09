# Why Layers Help

Layers aren't just a neat way to organise files. They save a lot of time and disk space.

Every layer is identified by its contents. If two images contain the **exact same layer**, Docker stores and downloads it only once.

Here's a pull of a second image that's built on the same base as one you already have:

```
docker pull myapp:2.0
a1b2c3d4e5f6: Already exists
b2c3d4e5f6a7: Already exists
c3d4e5f6a7b8: Pull complete
Status: Downloaded newer image for myapp:2.0
```

**Why did that pull finish so quickly?**
