# Make the Image the Program

Let's turn the webapp image into "the webapp program, in a box". Swap the last line of the Dockerfile:

```dockerfile
ENTRYPOINT ["webapp"]
```

Now options go straight to the app:

```sh
docker run --rm cq/webapp:ep --version
# webapp dev
```

And when you really do want a different program (to look around, say), replace the entrypoint with the `--entrypoint` flag:

```sh
docker run --rm --entrypoint ls cq/webapp:ep /src
#                            └── program  └── its arguments
```

# Assignment

In `~/docker-course`:

1. In your `Dockerfile`, replace the last line, `CMD ["webapp"]`, with `ENTRYPOINT ["webapp"]`. Leave everything else as it is.
2. Build it as **`cq/webapp:ep`** (no build argument, so the version is `dev`).
3. Run a container named **`cq-oneoff`** from `cq/webapp:ep`, passing just `--version`.
4. Run a container named **`cq-shell`** from `cq/webapp:ep` that uses `--entrypoint` to run `ls /src` instead.
5. Run, then submit `dk.entrypoint.entrypoint`.
