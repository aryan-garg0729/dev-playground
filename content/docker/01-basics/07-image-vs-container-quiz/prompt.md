# Images vs Containers

These two words get mixed up all the time, so let's pin them down.

- An **image** is the packed box: the app and all its files, frozen. It never changes, and it doesn't *do* anything by itself.
- A **container** is a running copy made *from* an image. You can make as many as you like from one image.

The classic way to remember it: **an image is a recipe, a container is a cake.** One recipe, many cakes. Eating one cake doesn't change the recipe.

Try it:

```sh
docker run hello-world
docker run hello-world
docker ps -a
# CONTAINER ID   IMAGE         COMMAND    STATUS                      NAMES
# 4f1c2a9b8e3d   hello-world   "/hello"   Exited (0) 3 seconds ago    eager_turing
# 9b2d7e6a1c0f   hello-world   "/hello"   Exited (0) 10 seconds ago   brave_curie
```

`docker ps -a` lists every container, even finished ones. Docker even gives each one a random silly name.

**You run `docker run hello-world` three times on a fresh computer. How many images and containers do you end up with?**
