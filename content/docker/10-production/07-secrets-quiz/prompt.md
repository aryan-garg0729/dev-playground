# The Leaky Build Arg

Mo needs an API key during the build to download the potion catalogue:

```dockerfile
FROM alpine:3.20
ARG POTION_KEY
RUN wget --header "Authorization: Bearer $POTION_KEY" -qO /catalogue.json https://potions.example/catalogue
```

```sh
docker build --build-arg POTION_KEY=sk-live-8f3a... -t elixirbytes/catalogue .
docker push registry.example.com/elixirbytes/catalogue
```

You checked in chapter 4 that an `ARG` isn't in the container's environment, so Mo is confident the key is safe.

**With Mo's version, who can read the key?**
