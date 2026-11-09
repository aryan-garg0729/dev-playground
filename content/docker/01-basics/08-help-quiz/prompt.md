# Ask Docker for Help

Docker has hundreds of options. Nobody remembers them all, and you don't need to. Docker can explain itself:

```sh
docker --help          # every command
docker run --help      # every option for one command
```

The output is long, so search it. `grep` keeps only the lines that contain a word:

```sh
docker run --help | grep name
#       --name string     Assign a name to the container
```

Reading `--help` is a real skill, and pros use it all the time.

Remember those finished hello-world containers piling up in `docker ps -a`? There's a `docker run` option that **removes the container automatically when it finishes**.

Run `docker run --help` (try `| grep -i remove`) and find it.

**Which option is it?**
