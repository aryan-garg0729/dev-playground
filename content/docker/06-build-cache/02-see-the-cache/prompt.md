# See the Cache

By default, `docker build` shows a tidy progress display that hides most of the detail. To see everything, use plain mode:

```sh
docker build --progress=plain -t cq/demo .
```

Every step gets a number and its own lines. A reused step looks like this:

```
#8 [4/4] RUN go build -o /usr/local/bin/webapp .
#8 CACHED
```

To keep a build's output so you can read it later, send it into a file with `tee`. It shows the output on screen *and* saves it:

```sh
docker build --progress=plain -t cq/demo . 2>&1 | tee build.log
```

(`2>&1` matters: Docker writes its progress to the "error" output, and this sends it along with the normal output.)

# Assignment

In `~/docker-course`:

1. Build the app once, so the cache is fresh: `docker build -t cq/webapp:cache .`
2. Build it **again**, in plain mode, saving the output to **`build.log`**:

   ```sh
   docker build --progress=plain -t cq/webapp:cache . 2>&1 | tee build.log
   ```

3. Open `build.log` and find the `RUN go build` step. It should say `CACHED`.
4. Run, then submit `dk.cache.see-the-cache`.
