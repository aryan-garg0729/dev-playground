# Why Does It Rebuild Everything?

Build the same Dockerfile twice and every step says `CACHED`: the second build takes about a second. Docker remembers each step's result under a **cache key** made from:

1. the cache key of the step before it
2. the instruction text itself
3. for `COPY` and `ADD`, a checksum of the files being copied

If all three match, the step is `CACHED`. If **any** step misses, every step after it runs again, because each builds on the one before.

Once brewd has real dependencies, downloading them will take four minutes:

```dockerfile
1  FROM golang:1.23-alpine
2  WORKDIR /src
3  COPY app/ .
4  RUN go mod download            # four minutes, once we have dependencies
5  RUN go build -o /usr/local/bin/brewd .
6  CMD ["brewd"]
```

Mo fixes a typo in a comment in `app/main.go` and rebuilds.

**Which steps run again?**
