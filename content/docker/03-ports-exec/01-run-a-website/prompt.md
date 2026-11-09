# Run a Website

So far our containers have only printed text. Let's run something real: a website.

**nginx** is a web server, a program that hands web pages to browsers. Start it in the background:

```sh
docker run -d --name cq-demo-web nginx:1.27-alpine
```

It's running, but open a browser and... nothing. Why? **Each container has its own private network.** nginx is waiting for visitors on port 80 *inside* its box, and nothing outside can reach it.

You have to **publish** a port: connect a door on your computer to a door on the container.

```sh
docker run -d --name cq-demo-web2 -p 8090:80 nginx:1.27-alpine
#                                    │    └── port inside the container
#                                    └── port on your computer
```

Now http://localhost:8090 reaches nginx. `docker ps` shows the connection in the PORTS column: `0.0.0.0:8090->80/tcp`.

# Assignment

1. Start a container named **`cq-web`** from **`nginx:1.27-alpine`**, detached, with your port **8081** connected to the container's port **80**.
2. Open http://localhost:8081 in your browser. You should see **Welcome to nginx!**
3. Run, then submit `dk.ports.run-a-website`. Leave `cq-web` running for the next lessons.

If you get `port is already allocated`, something else is using 8081. Find it with `docker ps` and remove it, or stop the other program.
