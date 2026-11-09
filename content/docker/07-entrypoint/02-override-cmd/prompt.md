# Replace the Command

Our webapp image only has a `CMD`:

```sh
docker image inspect -f 'entrypoint={{json .Config.Entrypoint}} cmd={{json .Config.Cmd}}' cq/webapp:1.3
# entrypoint=null cmd=["webapp"]
```

So whatever you type after the image name replaces `webapp` completely. That's handy for poking around an image:

```sh
docker run --rm cq/webapp:1.3 ls /src      # see what was copied in
docker run --rm cq/webapp:1.3 env          # see its settings
```

To run the app with an option, you have to type the program's name too:

```sh
docker run --rm cq/webapp:1.3 webapp --version
# webapp 1.3
```

# Assignment

1. Run a container named **`cq-peek`** from **`cq/webapp:1.3`** whose command is `ls /usr/local/bin` (no `--rm`: the checks need it).
2. Run a container named **`cq-version`** from **`cq/webapp:1.3`** whose command is `webapp --version`.
3. Run, then submit `dk.entrypoint.override-cmd`.
