# From Scratch

Let's go all the way down. `scratch` is a special name meaning "an empty image". The final stage will contain **one file**: our program.

```dockerfile
FROM scratch
COPY --from=build /out/webapp /webapp
ENTRYPOINT ["/webapp"]
```

There's no `/usr/local/bin` folder in an empty image and no shell to search the PATH, so we copy the program to `/webapp` and call it by its full path.

# Assignment

In `~/docker-course`:

1. In your `Dockerfile`, change the final stage to start `FROM scratch`. Copy the program to **`/webapp`** and make the entrypoint `["/webapp"]`. Keep the `ARG`, `ENV`, `LABEL` and `EXPOSE` lines: they're settings, not files, so they work on scratch too.
2. Build it as **`cq/webapp:scratch`**.
3. Compare sizes with `docker images cq/webapp`.
4. Try to get a shell: `docker run --rm --entrypoint /bin/sh cq/webapp:scratch`. It fails, because there's no shell to run!
5. Run, then submit `dk.multistage.scratch`.
