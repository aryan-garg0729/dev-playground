# Night Shift Triage

It's 03:12. Dex has been paged three times in an hour and hands you the laptop:

> "Three containers, three problems. Tell me what happened to the first, read me the batch number off the second, and put the third out **gracefully**. Its cleanup handler has to run. Write it all down. I'm getting coffee."

Run this from `~/elixirbytes` to set up the incident. It's random on purpose: your numbers won't match anyone else's.

```sh
docker rm -f eb-case-mana eb-case-fire eb-case-frost 2>/dev/null
docker run -d --name eb-case-mana alpine:3.20 sh -c 'echo "brewing mana"; sleep 2; echo "ERROR: cauldron cracked"; exit $(( $(od -An -N1 -tu1 /dev/urandom) % 90 + 10 ))'
docker run -d --init --name eb-case-fire alpine:3.20 sleep infinity
docker run -d --name eb-case-frost alpine:3.20 sh -c 'echo "frost batch $(od -An -N2 -tu2 /dev/urandom | tr -d " ") ready"; exec sleep infinity'
```

## Assignment

Everything you need is from this chapter. In `~/elixirbytes`, create **`triage.txt`** with two `key=value` lines:

```
mana_exit=<eb-case-mana's exit code>
frost_batch=<the number in eb-case-frost's "frost batch N ready" log line>
```

Then deal with the containers:

1. **`eb-case-mana`** has crashed. Find its exit code and **don't remove it**: it's evidence.
2. **`eb-case-frost`** is healthy. Read its logs, but **don't restart, stop or recreate it**.
3. **`eb-case-fire`** must be shut down **gracefully**: the polite signal, not the kill switch. Its exit code proves which one you sent.

Submit from `~/elixirbytes` to attack. The checks read the real exit code and batch number from Docker and compare them with your `triage.txt`, so the file has to match what actually happened. If you rerun the setup, the numbers change and you'll need to update `triage.txt`.

Every check is a hit on the boss. Defeat it to unlock **Images and tags**.
