# Override It at Run Time

Your `elixirbytes/brewd:ch5` image has `ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]` and no `CMD`. Because the wrapper passes `"$@"` along, anything after the image name in `docker run` reaches brewd, banner and all:

```sh
docker run --rm elixirbytes/brewd:ch5 --version
# brewd starting on 809264b461d9
# brewd dev
```

Sometimes you want something else entirely: a shell to look around, or a quick `ls`. `--entrypoint` swaps the program, and the words after the image name become its arguments:

```sh
docker run --rm --entrypoint cat elixirbytes/brewd:ch5 /usr/local/bin/entrypoint.sh
```

Docker records what a container actually ran as `.Path` (the program) and `.Args` (its arguments).

## Assignment

Make two containers from **`elixirbytes/brewd:ch5`**, and **don't use `--rm`**: the checks inspect them afterwards.

1. **`eb-oneoff`**: runs the image's normal entrypoint with the single argument **`--version`**, then exits.
2. **`eb-peek`**: replaces the entrypoint with **`ls`**, listing **`/usr/local/bin`**.

Submit from `~/elixirbytes`. The checks read each container's `.Path`, `.Args` and logs.
