# Too Many Commands

Look at what it took to run the load balancer in the networking chapter:

```sh
docker network create cq-net
docker run -d --name cq-app-a --network cq-net -e GREETING="Server A" cq/webapp:1.6
docker run -d --name cq-app-b --network cq-net -e GREETING="Server B" cq/webapp:1.6
docker run -d --name cq-lb --network cq-net -p 8095:80 -v "$(pwd)/lb":/etc/nginx/conf.d:ro nginx:1.27-alpine
```

Four long commands, in the right order, with every option spelled exactly right. Now imagine a teammate trying to set it up from your notes. Or ten services instead of three.

**Docker Compose** fixes this. You write everything down **once**, in a file called `compose.yaml`:

```yaml
services:
  app-a:
    image: cq/webapp:1.6
    environment:
      GREETING: Server A
  lb:
    image: nginx:1.27-alpine
    ports:
      - "8095:80"
```

Then one command starts the lot, and another removes it:

```sh
docker compose up -d
docker compose down
```

The file lives next to your code, so the whole team gets the same setup.

**What is Docker Compose?**
