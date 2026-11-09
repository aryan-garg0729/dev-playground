# Open the Potion Shop

Mo couldn't wait for us and wrote a Dockerfile for the launch. It "mostly works". Here it is, exactly as Mo committed it:

```dockerfile
FROM golang:1.23-alpine
COPY . /src
RUN cd /src/app
RUN go build -o /usr/local/bin/brewd .
EXPOSE 8080
CMD brewd
```

Then Mo ran it with:

```sh
docker run -d --name eb-shop -p 8080:18081 elixirbytes/brewd:shop
```

Everything that's wrong with this is something you've seen this chapter. Some of it fails loudly, some of it quietly ships junk, and one line is a habit that will bite in a later chapter.

## Assignment

Ship the shop properly, working in `~/elixirbytes`:

1. Save a fixed version as **`Dockerfile.shop`**. Keep Mo's base image, and make sure the image:
   - works in **`/src`**, which holds **only brewd's source**: not `testdata`, not the Dockerfiles, not your `triage.txt` or anything else from `~/elixirbytes`
   - compiles brewd to `/usr/local/bin/brewd`
   - starts brewd by default, in the right form
2. Build it as **`elixirbytes/brewd:shop`** with `docker build -f Dockerfile.shop ...`.
3. Run it as **`eb-shop`** so that **`http://localhost:18081`** on your machine reaches brewd.
4. Self-check: `curl -X POST 'localhost:18081/orders?potion=test'`, then `curl localhost:18081/orders`.

Submit from `~/elixirbytes` to attack. The checks read the image's command and working directory, list `/src` (hidden files included), then order a potion over HTTP and read it back.

Every check is a hit on the boss. Defeat it to unlock **The build cache**.
