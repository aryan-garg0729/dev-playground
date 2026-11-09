# -v vs --mount

There are two ways to write a mount:

```sh
# -v: short, positional; the TYPE is guessed from the first part
docker run -v eb-orders:/data ...            # volume
docker run -v "$PWD/orders":/data ...        # bind mount (it's a path)

# --mount: long, explicit key=value; errors instead of guessing
docker run --mount type=volume,src=eb-orders,dst=/data ...
docker run --mount type=bind,src="$PWD/orders",dst=/data,readonly ...
```

Mo is in `~/elixirbytes`, where an `orders/` folder exists, and types:

```sh
docker run -d --name eb-mo -v orders:/data elixirbytes/brewd:1.6.0
```

**What happens?**
