# Where Did My Data Go?

Our webapp has a tiny notes feature. Send it a note and it saves it in a file, `/data/notes.txt`:

```sh
docker run -d --name cq-temp -p 8090:8080 cq/webapp:1.6
curl -X POST "http://localhost:8090/notes?text=buy+milk"
# saved: buy milk
curl http://localhost:8090/notes
# buy milk
```

(`curl` sends web requests from the terminal. `-X POST` means "send something" rather than "fetch something".)

Now upgrade it, the way you would when a new version comes out: remove the container and start a new one.

```sh
docker rm -f cq-temp
docker run -d --name cq-temp -p 8090:8080 cq/webapp:1.6
curl http://localhost:8090/notes
#               <- empty!
```

Imagine that was a database full of customer orders.

**Why is the note gone?**
