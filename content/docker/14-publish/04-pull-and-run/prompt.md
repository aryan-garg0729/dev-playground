# Pull It Back

Let's play "production server". A fresh server has nothing on it, so delete your local copy of the name and get it back from the registry:

```sh
docker rmi localhost:15000/cq/webapp:1.0     # remove the local name
docker pull localhost:15000/cq/webapp:1.0    # download it from the registry
```

`docker run` would also pull it automatically if it's missing, just like `hello-world` in chapter 1.

# Assignment

1. Remove the local `localhost:15000/cq/webapp:1.0` name, then pull it from your registry.
2. Run **`cq-shipped`** from **`localhost:15000/cq/webapp:1.0`**, detached. (No port needed.)
3. Run, then submit `dk.publish.pull-and-run`.

On a real server you'd use the registry's real address instead of `localhost:15000`, and that's the only difference.
