# Your First Dockerfile

So far we've only run images other people built. A **Dockerfile** is the recipe for building your own: a text file of instructions, run top to bottom, where each one builds on the result of the one before.

```dockerfile
FROM python:3.12-slim          # start from an existing image
WORKDIR /app                   # cd into /app (created if missing) for every later step
COPY greet.py .                # copy files from your machine into the image
RUN pip install --no-cache-dir requests   # run a command at BUILD time, keep its file changes
CMD ["python", "greet.py"]     # the default command when a container STARTS
```

```sh
docker build -t demo/greet:v1 .
#            │                 └── the build context: the folder whose files COPY can see
#            └── tag the result
docker run --rm demo/greet:v1
```

Two instructions trip people up:

- **`RUN` runs once, at build time.** Its result (installed packages, a compiled binary) is saved into the image.
- **`CMD` doesn't run at build time.** It's only the default command for containers started later, and `docker run image other-command` replaces it.

Write `CMD` in its **JSON form**, `["program", "arg"]`, not as `CMD program arg`. The other form quietly wraps your command in `/bin/sh -c`; chapter 5 shows when that bites.

## Assignment

Priya's ticket: *"brewd gets an image. One stage is fine for now; we'll make it small later."*

1. Create **`~/elixirbytes/Dockerfile`** (in `~/elixirbytes`, not in `app/`) that:
   1. starts **from `golang:1.23-alpine`**
   2. works in **`/src`**
   3. copies the **contents** of `app/` into `/src`
   4. compiles brewd to **`/usr/local/bin/brewd`**. `go build -o <path> .` compiles the package in the current folder to `<path>`.
   5. starts brewd by default, using the JSON form
2. Build it from `~/elixirbytes` as **`elixirbytes/brewd:v1`**. The first build downloads the Go toolchain image (about 250 MB), so give it a minute.
3. Try it: `docker run --rm elixirbytes/brewd:v1 brewd --version` should print `brewd dev`.
4. Submit from `~/elixirbytes`.

The checks read the image's default command and working directory, look for the binary with `which`, and run `brewd --version` inside the image.
