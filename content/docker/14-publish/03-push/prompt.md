# Push an Image

To push an image, its **name must start with the registry's address**. So first give it a new name with `docker tag`, then push:

```sh
docker tag cq/demo:1.0 localhost:15000/cq/demo:1.0
docker push localhost:15000/cq/demo:1.0
# The push refers to repository [localhost:15000/cq/demo]
# 5f70bf18a086: Pushed
# 1.0: digest: sha256:... size: 528
```

Only the layers the registry doesn't already have are uploaded. Push a second version that shares the same base, and the base layers say `Layer already exists`.

A registry has a simple web API, so you can ask it what it holds:

```sh
curl http://localhost:15000/v2/cq/demo/tags/list
# {"name":"cq/demo","tags":["1.0"]}
```

# Assignment

1. Tag **`cq/webapp:prod`** as **`localhost:15000/cq/webapp:1.0`**.
2. Push it.
3. Ask the registry which tags it has for `cq/webapp`.
4. Run, then submit `dk.publish.push`.
