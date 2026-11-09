# Read the Logs

A detached container doesn't print to your terminal. Docker collects everything its main process writes to **stdout and stderr** and keeps it as the container's **logs**:

```sh
docker run -d --name eb-demo alpine:3.20 sh -c 'echo first; echo second; echo third'

docker logs eb-demo            # everything
# first
# second
# third

docker logs --tail 1 eb-demo   # only the last N lines
# third

docker logs -f eb-demo         # follow: keep printing new lines (Ctrl-C to stop following)
```

Logs survive after the container exits, which is what makes `docker logs` the first command anyone runs on a crashed container. They're deleted with the container.

This is also why containerised apps should log to stdout instead of to a file: Docker can only show what goes to stdout and stderr.

## Assignment

Our brewer worker will eventually log one line per potion. For now, fake it.

1. Start **`eb-ticker`** from `alpine:3.20`, detached, running a shell loop that prints `potion 1 brewed` through `potion 5 brewed`, one per line and about a second apart, then prints **`all potions brewed`** and exits:

   ```sh
   docker run -d --name eb-ticker alpine:3.20 sh -c 'for i in 1 2 3 4 5; do echo "potion $i brewed"; sleep 3; done; echo "all potions brewed"'
   ```

2. Run `docker logs -f eb-ticker` and watch the lines arrive. Following ends when the container exits.
3. Submit once it has exited. The checks read the exit code, the last two log lines and the full log.
