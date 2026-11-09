# Nicknames on the Network

Other containers find `cq-api` by its container name. But what happens when you replace it with a new version, under a new name? Everything that called `cq-api` would need changing.

Give it a **network alias** instead: an extra name on the network that describes its **job**, not the container.

```sh
docker run -d --name cq-api-v2 --network cq-net --network-alias backend cq/webapp:1.6
```

Now neighbours can use `http://backend:8080`. Next time, start the replacement with the same alias, remove the old one, and nobody else has to change a thing.

# Assignment

1. Start **`cq-api-v2`** from **`cq/webapp:1.6`**, detached, on **`cq-net`**, with the network alias **`backend`** and the setting **`GREETING=Hi`** (so you can tell it apart).
2. Check from a neighbour:

   ```sh
   docker run --rm --network cq-net alpine:3.20 wget -qO- -T 3 http://backend:8080/
   ```

3. Run, then submit `dk.networking.aliases`.
