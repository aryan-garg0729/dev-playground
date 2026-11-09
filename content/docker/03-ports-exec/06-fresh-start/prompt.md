# Where Did My Change Go?

You changed `index.html` and left a note in `cq-web`. Let's find out what happens to those changes when the container goes away.

Every container gets its **own private scratch space** on top of the image. Anything you change, add or delete inside the container lands there. The image itself is never touched.

- `docker stop` + `docker start`: same container, same scratch space. Your changes are still there.
- `docker rm`: the scratch space is deleted with the container. Your changes are **gone**.
- `docker run` again: a fresh container, a fresh, empty scratch space, straight from the untouched image.

Recipe and cake again: you can put sprinkles on your cake, but the recipe doesn't change. Bin the cake, bake a new one, and there are no sprinkles.

That's a feature, not a bug: every new container starts clean and predictable. (How to *keep* data on purpose is the storage chapter.)

# Assignment

1. Remove **`cq-web`** completely.
2. Start it again with exactly the same command as before: named `cq-web`, from `nginx:1.27-alpine`, detached, `-p 8081:80`.
3. Refresh http://localhost:8081. It's the default **Welcome to nginx!** page again, and `/tmp/note.txt` is gone.
4. Run, then submit `dk.ports.fresh-start`.
