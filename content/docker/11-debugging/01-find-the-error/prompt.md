# Find the Error in the Noise

Real apps log a lot. When something breaks, the one useful line is buried among thousands of boring ones. A few `docker logs` options help:

```sh
docker logs --tail 20 cq-x        # only the last 20 lines
docker logs --since 5m cq-x       # only the last 5 minutes
docker logs -t cq-x               # add a timestamp to every line
docker logs -f cq-x               # keep following new lines
```

And the best friend of every log reader is `grep`, which keeps only matching lines:

```sh
docker logs cq-x 2>&1 | grep -i error
```

(`2>&1` again: many apps write errors to the "error" output, so this makes sure `grep` sees those lines too.)

# Assignment

1. Start this noisy container. Copy and paste:

   ```sh
   docker run -d --name cq-noisy alpine:3.20 sh -c 'n=$(( $(date +%s) % 900 + 100 )); i=0; while [ $i -lt 500 ]; do i=$((i+1)); echo "INFO request $i handled"; [ $i -eq 317 ] && echo "ERROR payment failed code=$n" >&2; done; sleep infinity'
   ```

2. Find the one `ERROR` line in its 500+ lines of logs. Its error code is different for everyone.
3. Add it to `answers.txt` as `error_code=NUMBER`.
4. Run, then submit `dk.debug.find-the-error`.
