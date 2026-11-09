# Run Commands Inside

`docker run` starts a **new** container. But what if you want to look inside one that's **already running**?

That's what `docker exec` does: it runs an extra command inside a running container.

```sh
docker exec cq-web ls /usr/share/nginx/html
# 50x.html
# index.html
```

That's nginx's folder of web pages, and `index.html` is the "Welcome to nginx!" page you saw. Read it:

```sh
docker exec cq-web cat /usr/share/nginx/html/index.html
```

You can change things in there too. The `sh -c '...'` part lets you run a small shell command inside the box (you need it for things like `>`, which writes to a file):

```sh
docker exec cq-web sh -c 'echo "<h1>Hi</h1>" > /usr/share/nginx/html/index.html'
```

# Assignment

1. List the files in `/usr/share/nginx/html` inside **`cq-web`**.
2. Replace `index.html` inside `cq-web` with a page that says **Hello from inside the container**. Copy and paste:

   ```sh
   docker exec cq-web sh -c 'echo "<h1>Hello from inside the container</h1>" > /usr/share/nginx/html/index.html'
   ```

3. Refresh http://localhost:8081. Your page appears straight away, with no restart.
4. Run, then submit `dk.ports.exec`.
