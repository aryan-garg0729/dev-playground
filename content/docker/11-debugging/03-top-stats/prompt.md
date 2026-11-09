# Who's Eating the CPU?

Your computer's fan is screaming and everything is slow. Is it a container? Two commands tell you:

```sh
docker stats              # live CPU and memory for every running container (Ctrl+C to quit)
docker stats --no-stream  # one snapshot, then stop
# NAME         CPU %     MEM USAGE / LIMIT
# cq-api       0.00%     4.1MiB / 7.6GiB
# cq-lb        0.01%     3.2MiB / 7.6GiB
```

And to see the processes **inside** one container, without `exec`:

```sh
docker top cq-api
# UID    PID     PPID   CMD
# root   4123    4100   webapp
```

# Assignment

1. Start two workers. One of them (randomly) is stuck in a busy loop. Copy and paste:

   ```sh
   b=$(( $(date +%s) % 2 + 1 ))
   for n in 1 2; do
     if [ $n -eq $b ]; then cmd='while :; do :; done'; else cmd='sleep infinity'; fi
     docker run -d --name cq-worker-$n alpine:3.20 sh -c "$cmd"
   done
   ```

2. Use `docker stats --no-stream` to find which one is using the CPU. (`docker top` on each one shows what it's running too.)
3. Add the name to `answers.txt`, as `busiest=cq-worker-1` or `busiest=cq-worker-2`.
4. Stop **both** workers so your fan can rest: `docker stop cq-worker-1 cq-worker-2`. (Stop, don't remove: the check looks at them.)
5. Run, then submit `dk.debug.top-stats`.
