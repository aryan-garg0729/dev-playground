# What Are We Shipping?

Our app works. Now look at what we're actually shipping:

```sh
docker images cq/webapp:1.4
# REPOSITORY   TAG   SIZE
# cq/webapp    1.4   377MB
```

377 MB for a small web server. Let's look inside:

```sh
docker run --rm --entrypoint go cq/webapp:1.4 version
# go version go1.23.12 linux/arm64      <- the entire Go compiler is in there

docker run --rm --entrypoint ls cq/webapp:1.4 -l /usr/local/bin/webapp
# -rwxr-xr-x  1 root root  6.9M ...  /usr/local/bin/webapp
```

The program itself is about 7 MB. Everything else is luggage.

Big images download slowly, start slowly on new servers, and every extra program inside is one more thing an attacker could use.

**Where are the other ~370 MB coming from?**
