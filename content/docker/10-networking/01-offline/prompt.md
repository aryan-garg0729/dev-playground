# Cut the Cable

Every container gets its own private network connection. By default it's plugged into Docker's **default network**, which can reach the internet through your computer. That's how `apk add` worked inside builds.

You can also give a container **no network at all**, with `--network none`. That's useful for jobs that should never talk to anything, like processing a file full of secrets.

This command tries to load a web page and prints `online` or `offline`:

```sh
sh -c 'wget -q -T 3 -O /dev/null http://example.com && echo online || echo offline'
```

# Assignment

1. Run **`cq-online`** from **`alpine:3.20`** with that command, on the default network (no `--network` option).
2. Run **`cq-offline`** from **`alpine:3.20`** with the same command and **`--network none`**.
3. Compare their logs.
4. Run, then submit `dk.networking.offline`.

(If `cq-online` says `offline` too, check your own internet connection, or a VPN that blocks Docker.)
